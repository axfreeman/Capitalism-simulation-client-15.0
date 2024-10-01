// PATH: go-auth/controllers/auth.go

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorilla-client/api"
	"gorilla-client/config"
	"gorilla-client/db"
	"gorilla-client/models"
	"gorilla-client/utils"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Store = sessions.NewCookieStore([]byte("super-secret-password")) //TODO security
var Tpl *template.Template
var hash []byte

// registerHandler serves form for registering new users
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Enter RegisterHandler")
	Tpl.ExecuteTemplate(w, "register.html", nil)
}

// Services a post request to create a new registered user from user data in a form,
// validates the form, and checks for duplicates in the local RegisteredUser database.
//
// Synchronises with the server.
//
// Creates a local RegisteredUser record.
func RegisterAuthHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var req *http.Request
	var res *http.Response
	var ServerData ServerUserDetails

	utils.TraceInfo(utils.BrightGreen, "Enter RegisterAuthHandler")

	// validate the form
	if r.ParseForm() != nil {
		Tpl.ExecuteTemplate(w, "register.html", "Form incorrectly filled out. Try again")
	}

	// validate user name
	username := r.FormValue("username")
	if len(username) < 2 {
		Tpl.ExecuteTemplate(w, "register.html", "Username is too short")
		return
	}

	utils.TraceInfo(utils.BrightGreen, "User Name is valid")

	// check if username already exists in the local database
	if _, err = db.DataBase.FindRegisteredUser(username); err == nil {
		utils.TraceInfo(utils.BrightGreen, "User already exists")
		Tpl.ExecuteTemplate(w, "register.html", MessageData{Message: "User already exists", Username: "admin"})
		return
	}
	utils.TraceInfo(utils.BrightGreen, "User Name is new")

	// create hash from password
	password := r.FormValue("password")

	if hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		utils.TraceError(fmt.Sprint("bcrypt err:", err))
		Tpl.ExecuteTemplate(w, "register.html", MessageData{Message: fmt.Sprintf("Encryption problem. Please report this to the developer\n%v", err), Username: "admin"})
		return
	}
	utils.TraceInfo(utils.BrightGreen, "Pasword is valid")

	// Create a local Registered user
	registeredUser := models.NewRegisteredUser(username, string(hash), "")
	registeredUserServerRequest := models.RegisteredUserServerRequest{UserName: username}
	utils.TraceInfo(utils.BrightGreen, "Prototype new registered user created")

	// send the skeleton details to the server to construct a fullblown user
	// unless one already exists, in which case retrieve the details
	body, _ := json.Marshal(registeredUserServerRequest)

	req, err = http.NewRequest("POST", config.Config.ApiSource+"/admin/register", bytes.NewBuffer(body))

	if err != nil {
		utils.TraceErrorf("Error constructing server request: %v", err)
		Tpl.ExecuteTemplate(w, "register.html", MessageData{Message: fmt.Sprintf("Error constructing server request:%v", err), Username: "admin"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("x-api-key", config.Config.AdminKey)
	client := &http.Client{}

	res, err = client.Do(req)
	if err != nil {
		utils.TraceErrorf("Server returned error:%v", err)
		Tpl.ExecuteTemplate(w, "register.html", MessageData{Message: fmt.Sprintf("Server returned error:%v", err), Username: "admin"})
		return
	}
	respBody, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	utils.TraceInfof(utils.BrightGreen, "Server returned status %d and said:%s", res.StatusCode, string(respBody))

	// registered already on the server? No worries.
	if res.StatusCode == http.StatusConflict {
		utils.TraceInfo(utils.BrightGreen, "This user is already registered on the server. No new action taken")
		return
	}

	if res.StatusCode != http.StatusCreated {
		errorReport := fmt.Sprintf("The server could not create the new record and returned status code %d", res.StatusCode)
		utils.TraceError(errorReport)
		Tpl.ExecuteTemplate(w, "register.html", MessageData{Message: fmt.Sprintf("Could not create the new record: status code:%d", res.StatusCode), Username: "admin"})
		return
	}

	// WAS status, err := api.AdminPostRequest(config.Config.ApiSource+"/admin/register", body) DEPRECATED

	if res.StatusCode == http.StatusConflict {
		utils.TraceInfof(utils.BrightGreen, "User %s is already registered on the server. No worries", username)
	} else {
		utils.TraceInfof(utils.BrightGreen, "Server is registering the user")
	}

	//retrieve the apikey that the server generated
	//TODO generate secure apikeys (in the API)
	json.Unmarshal(body, &ServerData)
	registeredUser.ApiKey = ServerData.ApiKey

	// Save the user to the database
	db.DataBase.CreateRegisteredUser((registeredUser))
	Tpl.ExecuteTemplate(w, "login.html", nil)
}

// loginHandler serves a form for users to login with
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Enter LoginHandler")
	Tpl.ExecuteTemplate(w, "login.html", nil)
	utils.TraceInfo(utils.BrightGreen, "Exit LoginHandler")
}

// loginAuthHandler authenticates user login
func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	var registeredUser *models.RegisteredUser
	var err error

	utils.TraceInfo(utils.BrightGreen, "Enter LoginAuthHandler")

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	utils.TraceInfo(utils.BrightGreen, fmt.Sprintf("Request to log in from User %s with password %s", username, password))
	if registeredUser, err = db.DataBase.FindRegisteredUser(username); err != nil {
		utils.TraceError(fmt.Sprintf("User %s is not registered", username))
		Tpl.ExecuteTemplate(w, "login.html", "Check the username and the password")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(registeredUser.Password), []byte(password))
	if err != nil {
		utils.TraceError("Incorrect password")
		Tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	// Send the Registered User's name to the server and retrieve a fullblown user.
	user := models.NewUser(username)
	status, err := api.AdminGetRequest(config.Config.ApiSource+"/admin/user/"+username, &user)
	utils.TraceInfo(utils.BrightGreen, fmt.Sprintf("The server responded with status %d and error %v", status, err))
	if status != http.StatusOK {
		utils.TraceError("The server doesn't know this user, sorry")
		Tpl.ExecuteTemplate(w, "login.html", "Check username and password")
		return
	}
	// Override local registeredUser store with the apikey supplied by the server.
	// They should be the same anyhow but this is an added precaution.
	registeredUser.ApiKey = user.ApiKey

	// save the name in the authentication store
	session, _ := Store.Get(r, "session") // session struct has field make(map[interface{}]interface{})
	session.Values["userID"] = username
	session.Save(r, w) // save before writing to response/return from handler
	utils.TraceInfof(utils.BrightGreen, "User %s has successfully logged in with apikey %s", registeredUser.UserName, registeredUser.ApiKey)

	// Add the fullblown user to the client list of logged-in users
	models.LoggedInUsers[username] = user

	//Grab all the templates from the server
	//See note in DOCS folder
	api.FetchRemoteTemplates()

	//Grab this user's data from the server TODO degrade gracefully if this doesn't work
	utils.TraceInfof(utils.BrightGreen, "the user's current simulation is %d", user.CurrentSimulationID)

	// display the welcome screen
	user.CurrentPage = models.CurrentPageType{Url: "welcome.html", Id: 0}
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, MessageData{Message: "", Username: user.UserName})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Entered LogoutHandler")
	session, _ := Store.Get(r, "session")
	delete(session.Values, "userID")
	session.Save(r, w)
	Tpl.ExecuteTemplate(w, "login.html", "Logged Out")
}

// Auth adds authentication code to handler before returning handler
// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
func Auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session")
		content, ok := session.Values["userID"]
		if !ok {
			http.Redirect(w, r, "auth/login", http.StatusFound)
			return
		}
		utils.TraceInfof(utils.BrightGreen, "Auth was called and retrieved %s", content)

		// Check that the cookie refers to a logged in user
		_, ok = models.LoggedInUsers[content.(string)]
		if !ok {
			http.Redirect(w, r, "auth/login", http.StatusFound)
			return
		}

		// ServeHTTP calls f(w, r)
		// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
		HandlerFunc.ServeHTTP(w, r)
	}
}

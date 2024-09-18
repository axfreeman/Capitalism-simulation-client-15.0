// controllers.display.go

package controllers

import (
	"encoding/json"
	"errors"
	"gorilla-client/models"
	"gorilla-client/utils"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

// Simplified message type to pass into templates
// without calculating Views
type MessageData struct {
	Message  string
	Username string
}

type ServerUserDetails struct {
	Username string `json:"username"`
	ApiKey   string `json:"apikey"`
}

// Fetch the current user from the cookie Store
func CurrentUser(r *http.Request) *models.User {
	session, _ := Store.Get(r, "session")
	content := session.Values["userID"]
	return models.LoggedInUsers[content.(string)]
}

// Display the data that is available for the user who made this call
// Fetch the data from the client local store, not from the server
func AllData(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	utils.TraceInfof(utils.Green, "Get Data for user %s", user.UserName)
	data, err := json.MarshalIndent(user, " ", " ")

	if err != nil {
		utils.TraceErrorf("Error %v retrieving base data", err)
	}

	utils.TraceInfof(utils.Blue, "User %s asked to view base data", user.UserName)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Display the tables that are available for the user who made this call
// Fetch the data from the client local store, not from the server
func TableData(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	utils.TraceInfof(utils.Green, "Get Table Data for user %s", user.UserName)
	output := user.CreateTemplateData("This user's display data")
	templateData, _ := json.MarshalIndent(output, " ", " ")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(templateData))
}

// Report an error by redisplaying the current template with an error message
//
//	user.CurrentPageDetail.Url must be set with the template name
//
//	user: the current user
//	w: the ResponseWriter to which the message should be sent
//	message: the error message
func ReportError(user *models.User, w http.ResponseWriter, message string) {
	t := user.CreateTemplateData(message)
	utils.TraceError(t.Message)

	// use standard error page if no Current Page is set
	if len(user.CurrentPage.Url) < 1 {
		user.CurrentPage = models.CurrentPageType{Url: "errors.html", Id: 0}
	}
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, t)
}

// The state which follows each action.
var nextStates = map[string]string{
	`demand`:  `SUPPLY`,
	`supply`:  `TRADE`,
	`trade`:   `PRODUCE`,
	`produce`: `CONSUME`,
	`consume`: `INVEST`,
	`invest`:  `DEMAND`,
}

// pages for which redirection is OK.
func useLastVisited(last string) bool {
	if last == "" {
		return false
	}
	switch last {
	case
		`commodities.html`,
		`industries.html`,
		`classes.html`,
		`industry_stocks.html`,
		`class_stocks.html`,
		`index.html`,
		`/`:
		return true
	}
	return false
}

func FetchIDfromURL(r *http.Request) (int, error) {
	var idAsString string
	var ok bool
	var err error
	var id int
	if idAsString, ok = mux.Vars(r)["id"]; !ok {
		return 0, errors.New("unrecognised object id")
	}
	if id, err = strconv.Atoi(idAsString); err != nil {
		return 0, errors.New("malformed id")
	}
	return id, nil
}

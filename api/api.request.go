// api.server.go
// container for interacting with remote server

package api

import (
	"bytes"
	"encoding/json"
	"gorilla-client/config"
	"gorilla-client/db"
	"gorilla-client/models"
	"gorilla-client/utils"
	"log"

	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Prepare and send a request by a normal user for a protected service
// to the server using an api key.
//
//	apiKey:the api key
//	url:appended to apiSource to tell the server what to do.
//	Returns:byte array with server response, error if anything went wrong
func UserGetRequest(apiKey string, url string) ([]byte, error) {
	// uncomment for more verbose diagnostics
	// utils.TraceInfof(utils.Cyan, "UserGetRequest was called with apiKey %s and path %s", apiKey, url)
	resp, err := http.NewRequest("GET", config.Config.ApiSource+url, bytes.NewBuffer([]byte(`{"origin":"Simulation-client"}`)))
	if err != nil {
		utils.TraceInfof(utils.Red, "Malformed client request:%v", err)
		return nil, err
	}

	resp.Header.Add("Content-Type", "application/json")
	resp.Header.Add("x-api-key", apiKey)

	client := &http.Client{Timeout: time.Second * 5} // Timeout after 5 seconds
	res, _ := client.Do(resp)
	if res == nil {
		utils.TraceInfo(utils.Red, "Server is down or misbehaving")
		return nil, errors.New("server did not respond to ServerRequest")
	}

	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		utils.TraceInfo(utils.Red, fmt.Sprintf("Server rejected the request with status %s", res.Status))
		utils.TraceInfo(utils.Red, fmt.Sprintf("It said %s", string(b)))
		return nil, errors.New(string(b))
	}
	// uncomment for more verbose diagnostics
	// utils.TraceInfo(utils.Cyan, fmt.Sprintf("UserGetRequest succeeded"))
	return b, nil
}

// Sends a GET request to the api server using admin credentials
//
//	url: the endpoint of the API server to which the reqeust is being sent
//	target: the result will be placed here
//
//	 returns:
//	  status
//	  error if there is a failure, nil otherwise
func AdminGetRequest(url string, target any) (int, error) {
	var err error
	var resp *http.Request
	utils.TraceInfo(utils.Cyan, fmt.Sprintf("Admin request with url %s", url))
	resp, err = http.NewRequest("GET", url, nil)
	if err != nil {
		utils.TraceInfo(utils.Cyan, fmt.Sprintf("Error constructing server request:%v", err))
		return http.StatusInternalServerError, err //500
	}

	resp.Header.Add("x-api-key", config.Config.AdminKey)
	client := &http.Client{Timeout: time.Second * 2} // Timeout after 2 seconds
	res, err := client.Do(resp)
	if err != nil {
		utils.TraceInfo(utils.Cyan, fmt.Sprintf("Server returned error:%v", err))
		return http.StatusInternalServerError, err //500
	}
	if res == nil {
		utils.TraceInfo(utils.Cyan, "Server response was empty")
		return http.StatusNonAuthoritativeInfo, errors.New("server send an empty response") //203
	}

	if res.StatusCode != 200 {
		utils.TraceInfo(utils.Cyan, fmt.Sprintf("server responded with status code %d", res.StatusCode))
		return res.StatusCode, nil
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	jsonErr := json.Unmarshal(body, target)
	if jsonErr != nil {
		errorReport := fmt.Sprintf("could not unmarshal server response because of this error: %v", jsonErr)
		utils.TraceInfo(utils.Cyan, errorReport)
		utils.TraceInfo(utils.Cyan, fmt.Sprintf("The text from the server was %s", string(body)))
		return 422, errors.New(errorReport)
	}
	utils.TraceInfo(utils.Cyan, fmt.Sprintf("Request for data from endpoint %s accepted", url))
	return http.StatusOK, nil
}

// Sends a POST request to the api server using admin credentials
//
//	url: the endpoint of the API server to which the reqeust is being sent
//	target: the result will be placed here
//
//	returns:
//	 status
//	 error if there is a failure, nil otherwise
func AdminPostRequest(url string, body []byte) (int, error) {
	var err error
	var req *http.Request
	var res *http.Response

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("x-api-key", config.Config.AdminKey)

	if err != nil {
		utils.TraceInfof(utils.Cyan, "Error constructing server request: %v", err)
		return http.StatusBadRequest, err
	}
	client := &http.Client{}

	res, err = client.Do(req)
	if err != nil {
		utils.TraceInfo(utils.Cyan, fmt.Sprintf("Server returned error:%v", err))
		return res.StatusCode, err
	}
	respBody, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	utils.TraceInfof(utils.Cyan, "Server returned status %d and said:%s", res.StatusCode, string(respBody))

	// registered already on the server? No worries.
	if res.StatusCode == http.StatusConflict {
		utils.TraceInfo(utils.Cyan, "This user is already registered on the server. No new action taken")
		return res.StatusCode, err
	}

	if res.StatusCode != http.StatusCreated {
		errorReport := fmt.Sprintf("The server could not create the new record and returned status code %d", res.StatusCode)
		utils.TraceError(errorReport)
		return res.StatusCode, errors.New(errorReport)
	}

	return http.StatusCreated, nil
}

// Loads Templates
// See note in DOCS folder
func FetchRemoteTemplates() error {
	status, err := AdminGetRequest(config.Config.ApiSource+`/templates/templates`, &models.TemplateList)
	if err != nil {
		errorReport := fmt.Sprintf("Could not retrieve template information from server. Status %d, error message folows:\n%v", status, err)
		utils.TraceInfo(utils.BrightRed, errorReport)
		return errors.New(errorReport)
	}
	utils.TraceInfo(utils.Cyan, "Templates retrieved from server")
	return nil
}

// Populate the RegisteredUser database with data fetched from the remote server
func LoadRegisteredUsers() error {
	var RegisteredUserList []models.RegisteredUser // Temporary storage for initializing
	utils.TraceInfo(utils.BrightCyan, "Loading remote users")
	_, err := AdminGetRequest(config.Config.ApiSource+`/admin/users`, &RegisteredUserList)
	if err != nil {
		log.Fatal("server failed to return user data. Cannot continue")
	}

	for _, item := range RegisteredUserList {
		item.Password = `insecure` // TODO store hashed passwords on the server
		hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.TraceError(fmt.Sprint("bcrypt err:", err))
			return err
		}
		item.Password = string(hash)
		db.DataBase.CreateRegisteredUser(&item)
	}

	utils.TraceInfo(utils.BrightCyan, "Registered Users Loaded")
	return nil
}

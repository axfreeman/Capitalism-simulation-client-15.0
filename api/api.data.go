// api.data.go
// DataObject is the intermediary between the client and the server.

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorilla-client/models"
	"gorilla-client/utils"
)

// Retrieves the data for a single table from the server.
// Unmarshals the server response into the DataList of the receiver
//
//	apiKey: sent to the server to identify and authorize the user
//	d: target of the data
//
//	Return: nil if it worked
//	Return: error string if there was an error
func Fetch(apiKey string, d *models.Tabler) error {
	utils.TraceInfo(utils.BrightCyan, fmt.Sprintf("Fetching a table from server with api key %s and path %s", apiKey, d.ApiUrl))

	response, err := UserGetRequest(apiKey, d.ApiUrl)
	if err != nil {
		errorReport := fmt.Sprintf("ServerRequest produced the error %v", err)
		utils.TraceInfo(utils.Red, errorReport)
		return errors.New(errorReport)
	}

	if len(string(response)) == 0 {
		utils.TraceInfo(utils.BrightCyan, "INFORMATION: a server response to a fetch request was empty")
		return nil // no response is not an error, but don't process the result
	}

	// Populate the table
	jsonErr := json.Unmarshal(response, &d.Table)
	if jsonErr != nil {
		utils.TraceInfof(utils.Red, "Server response could not be unmarshalled because: %v", jsonErr)
		utils.TraceInfof(utils.Red, "The server response was %s\n", response)
		return errors.New("server response could not be unmarshalled")
	}
	return nil
}

// Fetches a simulation and associated tables from the api server.
// NOTE the server works out who the user is from the apiKey
// NOTE the server must first be told this user's current simulation ID
//
//	user: supplies apiKey and simulationID that uniquely identify the simulation
//
//	returns:
//	  err if anything goes wrong
func FetchTables(user *models.User) error {
	// Fetch all the simulations for this user (regardless of ID)
	err := Fetch(user.ApiKey, &user.Simulations)
	if err != nil {
		return err
	}

	// Fetch all the tables in this simulation
	// NOTE the server knows the simulationID because it knows about the user
	newTableSet := models.NewTableSet()
	for key, value := range newTableSet {
		err = Fetch(user.ApiKey, &value)
		if err != nil {
			utils.TraceErrorf("Could not retrieve server data with key %s because of error %s", key, err.Error())
		}
	}

	// set the stocklist of every industry
	// TODO filter so the stocklist only contains relevant stocks
	industries := *(newTableSet[`industries`].Table.(*[]models.Industry))
	stocks := newTableSet[`industry stocks`].Table.(*[]models.IndustryStock)
	for ind := range industries {
		industries[ind].Stocks = stocks
	}

	user.TableSets = append(user.TableSets, &newTableSet)
	return nil
}

package controllers

import (
	"encoding/json"
	"fmt"
	"gorilla-client/api"
	"gorilla-client/models"
	"gorilla-client/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Structure to hold result of the server's response to a clone request
type CloneResult struct {
	Message       string `json:"message"`
	StatusCode    int    `json:"statusCode"`
	Simulation_id int    `json:"simulation_id"`
}

// Create a new simulation for the user, from the template specified by the 'id' parameter.
// Set the user's currentSimulationID to point to this simulation
// Create a Manager for this simulation
// Initialise the Manager's viewed and comparator timeStamps to 0
// Create the first stage of the simulation.

func CreateSimulation(w http.ResponseWriter, r *http.Request) {
	var s string
	var ok bool
	var err error
	var body []byte
	var manager *models.Manager

	user := CurrentUser(r)
	utils.TraceInfof(utils.Green, "Clone Simulation was called by user %s", user.UserName)
	user.CurrentPage = models.CurrentPageType{Url: "user-dashboard.html", Id: 0}

	if s, ok = mux.Vars(r)["id"]; !ok {
		ReportError(user, w, "Unrecognisable URL. Please report this to the developer")
		return
	}

	requestedSimulation, _ := strconv.Atoi(s)
	utils.TraceInfof(utils.Green, "Request to clone simulation %d", requestedSimulation)

	// Ask server to create clone and supply simulation id. Do not load tables yet
	if body, err = api.UserGetRequest(user.ApiKey, `/clone/`+s); err != nil {
		ReportError(user, w, fmt.Sprintf("There was a problem. Please report this to the developer%v", err))
		return
	}

	// read the simulation id
	var result CloneResult
	if err = json.Unmarshal(body, &result); err != nil {
		ReportError(user, w, fmt.Sprintf("There was a problem with the server's response. Please report this to the developer%v", err))
		return
	}

	utils.TraceInfo(utils.Green, ("Server responded to clone request:"))
	utils.TraceInfo(utils.Green, ` `+string(body))

	// Set the current simulation
	utils.TraceInfof(utils.Green, "Setting current simulation to %d", result.Simulation_id)
	user.CurrentSimulationID = result.Simulation_id

	// TODO New code starts here

	// Fetch the manager of this object
	if manager, err = api.ReplacementFetchManager(user, result.Simulation_id); err != nil {
		utils.TraceErrorf("Could not retrieve the manager object with apikey %s", user.ApiKey)
		ReportError(user, w, "oops")
		return
	}
	utils.TraceInfo(utils.BrightRed, " Retrieved the Manager object, phew")

	// Create a new Simulation object
	newSimulation := models.NewSimulation()

	// Make a fresh copy of the manager
	newSimulation.Manager = *manager
	user.Simulations[user.CurrentSimulationID] = newSimulation

	// Set the manager's timeStamps and initial state, and create the States map
	user.Simulations[user.CurrentSimulationID].Manager.TimeStamp = 0
	user.Simulations[user.CurrentSimulationID].Manager.ViewedTimeStamp = 0
	user.Simulations[user.CurrentSimulationID].Manager.ComparatorTimeStamp = 0
	user.Simulations[user.CurrentSimulationID].Manager.States = make(map[int]string)
	user.ReplacementSetCurrentState("DEMAND")

	utils.TraceInfo(utils.BrightRed, " Set up the managers initial state and timestamps, phew")

	// Fetch the data from the first Stage
	if err = api.ReplacementFetchStage(user); err != nil {
		utils.TraceErrorf("Could not retrieve the data for simulation with id %d using apikey %s", user.CurrentSimulationID, user.ApiKey)
		ReportError(user, w, "oops")
		return
	}
	utils.TraceInfo(utils.BrightRed, " Retrieved the Data, phew")

	// Convert the data to add pointers in place of Id field
	api.ConvertStage(newSimulation.Stages[user.Simulations[user.CurrentSimulationID].Manager.ViewedTimeStamp])

	utils.TraceInfo(utils.BrightRed, " Converted the Data, phew")

	// TODO Deprecated code from here
	// Fetch everything for the new simulation from the server.
	// (until now we only told the server to create it - now we want it).
	// Add this to the user's Tables
	err = api.CreateStage(user)
	if err != nil {
		utils.TraceErrorf("Could not retrieve the requested data with apikey %s and simulation id %d", user.ApiKey, result.Simulation_id)
		ReportError(user, w, "oops")
		return
	}

	utils.TraceInfo(utils.Green, ("Setting current state to DEMAND"))
	user.SetCurrentState("DEMAND")

	simstring, _ := json.MarshalIndent(user.Managers, " ", " ")
	utils.TraceLogf(utils.BrightYellow, "FetchTables retrieved the simulation %s", string(simstring))
	tablestring, _ := json.MarshalIndent(user.Stages, " ", " ")
	utils.TraceLogf(utils.BrightYellow, "FetchTables retrieved the tables %s", string(tablestring))

	// Initialise the timeStamp so that we are viewing the first Stage.
	// As the user moves through the circuit, this timestamp will move forwards.
	// Each time we move forward, a new Stage will be created.
	// This allows the user to view and compare with previous stages of the simulation.
	*user.GetViewedTimeStamp() = 0
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

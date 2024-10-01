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

	// Fetch the manager of this object
	if manager, err = api.FetchManager(user, result.Simulation_id); err != nil {
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
	user.GetCurrentSimulation().Manager.TimeStamp = 0
	user.GetCurrentSimulation().Manager.ViewedTimeStamp = 0
	user.GetCurrentSimulation().Manager.ComparatorTimeStamp = 0
	user.GetCurrentSimulation().Manager.States = make(map[int]string)
	user.SetCurrentState("DEMAND")

	utils.TraceInfo(utils.BrightRed, " Set up the manager's initial state and timestamps, phew")

	// Fetch the data from the first Stage
	if err = api.FetchStage(user); err != nil {
		utils.TraceErrorf("Could not retrieve the data for simulation with id %d using apikey %s", user.CurrentSimulationID, user.ApiKey)
		ReportError(user, w, "oops")
		return
	}
	utils.TraceInfo(utils.BrightRed, " Retrieved the Data, phew")

	// Convert the data to add pointers in place of Id field
	api.ConvertStage(user.GetCurrentStage())
	// WAS 	api.ConvertStage(newSimulation.Stages[user.Simulations[user.CurrentSimulationID].Manager.ViewedTimeStamp])

	utils.TraceInfo(utils.BrightRed, " Converted the Data, phew")

	simstring, _ := json.MarshalIndent(user.Simulations[user.CurrentSimulationID], " ", " ")
	utils.TraceLogf(utils.BrightYellow, "FetchTables retrieved the simulation %s", string(simstring))

	// Initialise all timeStamps so we are viewing the first Stage.
	// As the user moves through the circuit, timestamp will move forwards.
	// Each time we move forward, a new Stage will be created.
	// The user can move the ViewedTimeStamp backwards and forward to view
	// the history. At present the ComparatorTimeStamp stays one step
	// behind the ViewedTimeStamp but we may change this in future, for
	// example to compare one period with the previous one.
	user.GetCurrentSimulation().Manager.TimeStamp = 0
	user.GetCurrentSimulation().Manager.ViewedTimeStamp = 0
	user.GetCurrentSimulation().Manager.ComparatorTimeStamp = 0
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

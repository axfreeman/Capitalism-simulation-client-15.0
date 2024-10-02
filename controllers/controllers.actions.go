// display.actions.go
// This module processes the actions that take the simulation through
// a circuit - Demand, Supply, Trade, Produce, Consume, Invest

package controllers

import (
	"gorilla-client/api"
	"gorilla-client/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// Handles requests for the server to take an action comprising a stage
// of the circuit (demand,supply, trade, produce, invest), corresponding
// to a button press. This is specified by the URL parameter 'act'.
//
// Having requested the action from ths server, sets 'state' to the next
// stage of the circuit and redisplays whatever the user was looking at.
//
//	user.CurrentPageDetail.Url will be used to display errors if set
//	otherwise, a standard error page will be displayed
func ActionHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var action string
	var ok bool

	user := CurrentUser(r)
	utils.TraceInfof(utils.Green, "Processing action for user %s", user.UserName)

	// Find the requested action
	if action, ok = mux.Vars(r)["action"]; !ok {
		ReportError(user, w, "Poorly specified action in the URL")
		return
	}
	utils.TraceInfof(utils.Green, "User requested action %s", action)

	// Tell the API server to perform the action
	if _, err = api.UserGetRequest(user.ApiKey, `/action/`+action); err != nil {
		ReportError(user, w, "The server could not complete the action")
		return
	}

	// Create a new Stage and Append it to Datasets. Set the TimeStamps,
	// moving the comparator to immediately preceding stage
	simulation := user.GetCurrentSimulation()
	manager := &simulation.Manager
	manager.ComparatorTimeStamp = manager.ViewedTimeStamp
	manager.ViewedTimeStamp += 1
	manager.TimeStamp += 1

	// Fetch the data from the server and append it to Stages.
	if err = api.FetchStage(user); err != nil {
		ReportError(user, w, "The server completed the action but did not send back any data.")
		return
	}
	utils.TraceInfof(utils.Green, "Fetched a new set of tables")

	// Convert the data to add pointers in place of Id field
	api.ConvertStage(user.GetCurrentStage())
	utils.TraceInfof(utils.Green, "Converted the tables")

	// Set the state so that the simulation can proceed to the next action.
	user.SetCurrentState(nextStates[action])

	// Choose which page to display, depending on what the user was looking at
	utils.TraceInfof(utils.Green, "The last page this user visited was %v ", user.CurrentPage.Url)

	if useLastVisited(user.CurrentPage.Url) {
		Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
	} else {
		Tpl.ExecuteTemplate(w, "user-dashboard.html", user.CreateTemplateData(""))
	}
}

// View the previous stage of the simulation
// Comparator stays one step behind Viewed
// Later we can develop more sophisticated logic
// If Comparator is at the start, do nothing
func Back(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.Green, "Back was requested")
	u := CurrentUser(r)
	m := &CurrentUser(r).GetCurrentSimulation().Manager

	// View one earlier stage and compare it with the preceding
	if m.ComparatorTimeStamp > 0 {
		m.ComparatorTimeStamp--
		m.ViewedTimeStamp--
	}
	utils.TraceInfof(utils.Green, "Viewing timeStamp %d with comparator %d", m.ViewedTimeStamp, m.ComparatorTimeStamp)

	// Display appropriate page depending what the user was looking at
	if useLastVisited(u.CurrentPage.Url) {
		Tpl.ExecuteTemplate(w, u.CurrentPage.Url, u.CreateTemplateData(""))
	} else {
		Tpl.ExecuteTemplate(w, "index.html", u.CreateTemplateData(""))
	}
}

// View the previous stage of the simulation
// Comparator stays at least one step behind Viewed
// Later we can develop more sophisticated logic
// If Viewed is at the current stage, do nothing
func Forward(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.Green, "Forward was requested")
	u := CurrentUser(r)
	m := &CurrentUser(r).GetCurrentSimulation().Manager

	if m.ViewedTimeStamp < m.TimeStamp {
		m.ViewedTimeStamp++
		m.ComparatorTimeStamp++
	}

	utils.TraceInfof(utils.Green, "Viewing %d with comparator %d", m.ViewedTimeStamp, m.ComparatorTimeStamp)
	if useLastVisited(u.CurrentPage.Url) {
		Tpl.ExecuteTemplate(w, u.CurrentPage.Url, u.CreateTemplateData(""))
	} else {
		Tpl.ExecuteTemplate(w, "index.html", u.CreateTemplateData(""))
	}
}

// TODO not working yet
func SwitchSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData("Sorry, Switching Simulations is not ready yet"))
}

// TODO not working yet
func DeleteSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData("Sorry, Deleting a Simulation is not ready yet"))

}

// TODO not working yet
func RestartSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData("Sorry, Restarting a Simulation is not ready yet"))
}

// Quick and Dirty download method
// TODO rewrite
func Download(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	newStage := api.FetchStage(user)
	utils.UNUSED(newStage)
	// type listItem struct {
	// 	filename string
	// 	object   any
	// }
	// var f *os.File
	// var err error
	// outputList := make([]listItem, 5)
	// outputList[0] = listItem{`commodities.json`, (*newStage)[`commodities`]}
	// outputList[1] = listItem{`industries.json`, (*newStage)[`industries`]}
	// outputList[2] = listItem{`classes.json`, (*newStage)[`classes`]}
	// outputList[3] = listItem{`industry-stocks.json`, (*newStage)[`industry stocks`]}
	// outputList[4] = listItem{`class-stocks.json`, (*newStage)[`class stocks`]}
	// for i := range outputList {
	// 	out, _ := json.MarshalIndent(outputList[i].object, "", "")
	// 	f, err = os.Create(`./dump/` + outputList[i].filename)
	// 	if err != nil {
	// 		utils.TraceErrorf("Error %v creating download file %v", err, outputList[i].filename)
	// 		return
	// 	}
	// 	defer f.Close()
	// 	_, err = f.Write(out)
	// 	if err != nil {
	// 		utils.TraceErrorf("Error %v downloading to file %s", err, outputList[i].filename)
	// 		return
	// 	}
	// }
}

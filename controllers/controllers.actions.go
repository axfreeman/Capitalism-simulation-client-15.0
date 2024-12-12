// display.actions.go
// This module processes the actions that take the simulation through
// a circuit - Demand, Supply, Trade, Produce, Consume, Invest

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"simulation-client/api"
	"simulation-client/config"
	"simulation-client/logging"
	"simulation-client/views"
	"strconv"

	"github.com/gorilla/mux"
)

// Simplified message type to pass into templates
// without calculating Views
type messageData struct {
	Message  string
	Username string
}

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
	logging.TraceInfof(logging.Green, "Processing action for user %s", user.UserName)

	// Find the requested action
	if action, ok = mux.Vars(r)["action"]; !ok {
		ReportError(user, w, "Poorly specified action in the URL")
		return
	}
	logging.TraceInfof(logging.Green, "User requested action %s", action)

	if action == `setpricesform` {
		// Just display the form
		// Don't do anything else
		// The action is processed by setpriceshandle when user submits
		SetPricesFormDisplay(w, r)
		return
	} else if action == `setpriceshandle` {
		// Special case, because user can change prices at any point in the simulation
		// but the price change is nevertheless registered as a stage, so that it
		// can be inspected and traced

		// TODO some code here to fix the lastvistedpage display
		// TODO this isn't easy so at this point in development
		// we just have a scaffold sufficient to test the effects
		// of a price change
		SetPricesPostHandler(w, r)
	} else {
		// Tell the API server to perform the action
		if _, err = api.UserGetRequest(user.ApiKey, `/action/`+action); err != nil {
			ReportError(user, w, "The server could not complete the action")
			return
		}
	}
	// The action worked. Now retrieve the results from the API and store them locally.
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
	logging.TraceInfof(logging.Green, "Fetched a new set of tables")

	// Fetch the trace table
	// TODO this could get very big. Can we do an incremental fetch?
	if err = api.Fetch(user.ApiKey, simulation.Trace); err != nil {
		logging.TraceErrorf("Could not retrieve trace data for simulation with id %d using apikey %s", user.CurrentSimulationID, user.ApiKey)
		ReportError(user, w, "oops")
		return
	}
	logging.TraceInfof(logging.Green, "Refreshed the trace table")

	// Convert the data to add pointers in place of Id field
	api.ConvertStage(user.GetCurrentStage(), manager)

	// Set the state so that the simulation can proceed to the next action.
	user.SetCurrentState(nextStates[action])

	// Choose which page to display, depending on what the user was looking at
	logging.TraceInfof(logging.Green, "The last page this user visited was %v ", user.CurrentPage.Url)

	if useLastVisited(user.CurrentPage.Url) && action != `setprices` {
		views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, ""))
	} else {
		views.Tpl.ExecuteTemplate(w, "user-dashboard.html", views.CreateTemplateData(user, ""))
	}
}

// View the previous stage of the simulation
// Comparator stays one step behind Viewed
// Later we can develop more sophisticated logic
// If Comparator is at the start, do nothing
func Back(w http.ResponseWriter, r *http.Request) {
	logging.TraceInfo(logging.Green, "Back was requested")
	u := CurrentUser(r)
	m := &CurrentUser(r).GetCurrentSimulation().Manager

	// View one earlier stage and compare it with the preceding
	if m.ComparatorTimeStamp > 0 {
		m.ComparatorTimeStamp--
		m.ViewedTimeStamp--
	}
	logging.TraceInfof(logging.Green, "Viewing timeStamp %d with comparator %d", m.ViewedTimeStamp, m.ComparatorTimeStamp)

	// Display appropriate page depending what the user was looking at
	if useLastVisited(u.CurrentPage.Url) {
		views.Tpl.ExecuteTemplate(w, u.CurrentPage.Url, views.CreateTemplateData(u, ""))
	} else {
		views.Tpl.ExecuteTemplate(w, "index.html", views.CreateTemplateData(u, ""))
	}
}

// View the previous stage of the simulation
// Comparator stays at least one step behind Viewed
// Later we can develop more sophisticated logic
// If Viewed is at the current stage, do nothing
func Forward(w http.ResponseWriter, r *http.Request) {
	logging.TraceInfo(logging.Green, "Forward was requested")
	u := CurrentUser(r)
	m := &CurrentUser(r).GetCurrentSimulation().Manager

	if m.ViewedTimeStamp < m.TimeStamp {
		m.ViewedTimeStamp++
		m.ComparatorTimeStamp++
	}

	logging.TraceInfof(logging.Green, "Viewing %d with comparator %d", m.ViewedTimeStamp, m.ComparatorTimeStamp)
	if useLastVisited(u.CurrentPage.Url) {
		views.Tpl.ExecuteTemplate(w, u.CurrentPage.Url, views.CreateTemplateData(u, ""))
	} else {
		views.Tpl.ExecuteTemplate(w, "index.html", views.CreateTemplateData(u, ""))
	}
}

// Set the DisplayDimension of this simulation
//
//	Does not check for validity. Just don't be bloody stupid.
//
//	DisplayDimension: the dimension to set - either `Size`, `Value`, or `Price`.
func SetDisplayDimension(w http.ResponseWriter, r *http.Request, displayDimension string) {
	logging.TraceInfof(logging.Green, "Set Display Dimension to %s was requested", displayDimension)
	u := CurrentUser(r)
	m := &CurrentUser(r).GetCurrentSimulation().Manager

	logging.TraceInfof(logging.Green, "Display dimension will be changed from %s to %s", m.DisplayDimension, displayDimension)
	m.DisplayDimension = displayDimension

	if useLastVisited(u.CurrentPage.Url) {
		views.Tpl.ExecuteTemplate(w, u.CurrentPage.Url, views.CreateTemplateData(u, ""))
	} else {
		views.Tpl.ExecuteTemplate(w, "index.html", views.CreateTemplateData(u, ""))
	}
}

// Set the DisplayDimension of this simulation so that use values are displayed
func DisplaySize(w http.ResponseWriter, r *http.Request) {
	SetDisplayDimension(w, r, `Size`)
}

// Set the DisplayDimension of this simulation so that values are displayed
func DisplayValue(w http.ResponseWriter, r *http.Request) {
	SetDisplayDimension(w, r, `Value`)
}

// Set the DisplayDimension of this simulation so that prices are displayed
func DisplayPrice(w http.ResponseWriter, r *http.Request) {
	SetDisplayDimension(w, r, `Price`)
}

// TODO not working yet
func SwitchSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, "Sorry, Switching Simulations is not ready yet"))
}

// TODO not working yet
func DeleteSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, "Sorry, Deleting a Simulation is not ready yet"))

}

// TODO not working yet
func RestartSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, "Sorry, Restarting a Simulation is not ready yet"))
}

// Quick and Dirty download method
// TODO rewrite
func Download(w http.ResponseWriter, r *http.Request) {
	// user := CurrentUser(r)
	// newStage := api.FetchStage(user)
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
	// outputList[3] = listItem{`industry-stocks.json`, (*newStage)[`industry_stocks`]}
	// outputList[4] = listItem{`class-stocks.json`, (*newStage)[`class stocks`]}
	// for i := range outputList {
	// 	out, _ := json.MarshalIndent(outputList[i].object, "", "")
	// 	f, err = os.Create(`./dump/` + outputList[i].filename)
	// 	if err != nil {
	// 		logging.TraceErrorf("Error %v creating download file %v", err, outputList[i].filename)
	// 		return
	// 	}
	// 	defer f.Close()
	// 	_, err = f.Write(out)
	// 	if err != nil {
	// 		logging.TraceErrorf("Error %v downloading to file %s", err, outputList[i].filename)
	// 		return
	// 	}
	// }
}

// Process post request to change prices, submitted from form
// TODO underdevelopment
func SetPricesPostHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var req *http.Request
	var res *http.Response
	var price float64

	user := CurrentUser(r)

	// TODO validate the form
	if r.ParseForm() != nil {
		views.Tpl.ExecuteTemplate(w, "Commodity.html", views.CreateTemplateData(user, "Incorrect details. Try again"))
	}

	form := r.Form
	logging.TraceInfof(logging.Purple, "User %s submitted a price change form containing %v", user.UserName, form)

	type PriceRequest struct {
		CommodityId  int     `json:"commodityId"`
		SimulationId int     `json:"simulationId"`
		UnitPrice    float64 `json:"unitPrice"`
	}

	l := len(form)
	PricesRequest := make([]*PriceRequest, l)
	i := 0
	for k, v := range form {
		price, err = strconv.ParseFloat(v[0], 64)
		if err != nil {
			logging.TraceErrorf("Non-numeric price submitted %v", err)
			// TODO flag the error
			return
		}
		n, _ := strconv.Atoi(k) // this is set by the client so should always be valid...

		priceRequest := PriceRequest{
			CommodityId:  n,
			SimulationId: user.CurrentSimulationID,
			UnitPrice:    price,
		}
		PricesRequest[i] = &priceRequest
		i = i + 1
	}

	body, err := json.Marshal(&PricesRequest)
	// body, err := json.Marshal(&PricesRequest)
	if err != nil {
		log.Printf("Failed to marshal body: %s", err)
		return
	}

	req, rerr := http.NewRequest("POST", config.Config.ApiSource+"/action/setprices", bytes.NewBuffer(body))

	if rerr != nil {
		logging.TraceErrorf("Error constructing server request: %v", err)
		views.Tpl.ExecuteTemplate(w, "register.html", messageData{Message: fmt.Sprintf("Error constructing server request:%v", err), Username: "admin"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("x-api-key", user.ApiKey)
	client := &http.Client{}

	res, err = client.Do(req)
	if err != nil {
		logging.TraceErrorf("Server returned error:%v", err)
		views.Tpl.ExecuteTemplate(w, "errors.html", messageData{Message: fmt.Sprintf("Server returned error:%v", err), Username: "admin"})
		return
	}
	// respBody, _ := io.ReadAll(res.Body)

	defer res.Body.Close()

	// logging.TraceInfof(logging.BrightGreen, "Server returned status %d and said:%s", res.StatusCode, string(respBody))
	if res.StatusCode != http.StatusOK {
		logging.TraceInfof(logging.BrightGreen, "The server didn't like this and returned code %d", res.StatusCode)
		return
	}

	// Fetch the trace table (because it has been modified by a route other than performing an action)
	// TODO this could get very big. Can we do an incremental fetch?
	simulation := user.GetCurrentSimulation()
	if err = api.Fetch(user.ApiKey, simulation.Trace); err != nil {
		logging.TraceErrorf("Could not retrieve trace data for simulation with id %d using apikey %s", user.CurrentSimulationID, user.ApiKey)
		ReportError(user, w, "oops")
		return
	}
	logging.TraceInfof(logging.Green, "Refreshed the trace table")

	// TODO display the last-used page
}

// Display the setprices form
func SetPricesFormDisplay(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	logging.TraceInfof(logging.BrightGreen, "User %s entered SetPricesAuthHandler", user.UserName)
	views.Tpl.ExecuteTemplate(w, "set-prices.html", views.CreateTemplateData(user, ""))
}

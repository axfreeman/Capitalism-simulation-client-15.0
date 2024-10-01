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
func Fetch(apiKey string, d *models.Table) error {
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

// Replace Id fields with pointers. This makes for  legible code and faster access.
//
//	newStage: a Stage, which has been populated by FetchStage
func ConvertStage(stage *models.Stage) {
	// fmt.Printf("Entering ConvertStage with stage\n %v\n", stage)
	industries := *(*stage)[`industries`].Table.(*[]models.Industry)
	industryStocks := *(*stage)[`industry stocks`].Table.(*[]models.IndustryStock)
	classes := *(*stage)[`classes`].Table.(*[]models.Class)
	classStocks := *(*stage)[`class stocks`].Table.(*[]models.ClassStock)
	commodities := *(*stage)[`commodities`].Table.(*[]models.Commodity)

	// set the Commodity, Sales Stock, Money stock, Industrial stocks (=Constant capital) and Social stock (=Variable Capital) of every industry
	for ind := range industries {
		// fmt.Printf("Convert Industries is Processing industry %s\n", industries[ind].Name)
		industries[ind].Constant = make([]*models.IndustryStock, 0)
		for i := range industryStocks {
			if industryStocks[i].IndustryId == industries[ind].Id {
				// fmt.Printf("Convert Industries is Processing stock %s\n", industryStocks[i].Name)
				industryStocks[i].IndustryAddress = &industries[ind]
				industryStocks[i].IndustryName = industries[ind].Name
				switch industryStocks[i].UsageType {
				case `Money`:
					// fmt.Println("Money Stock")
					industries[ind].Money = &(industryStocks[i])
				case `Production`:
					// fmt.Println("Production Stock")
					if industryStocks[i].Origin == `SOCIAL` {
						// fmt.Println("Social Origin Stock")
						industries[ind].Variable = &(industryStocks[i])
					} else {
						// fmt.Println("Production Origin Stock")
						industries[ind].Constant = append(industries[ind].Constant, &(industryStocks[i]))
					}
				case `Sales`:
					// fmt.Println("Sales Stock")
					industries[ind].Sales = &(industryStocks[i])
				default:
				}
			}
		}
	}

	// set the Sales Stock, Money stock and Consumption stocks of every class
	for c := range classes {
		for s := range classStocks {
			if classStocks[s].ClassId == classes[c].Id {
				classStocks[s].ClassAddress = &classes[c]
				classStocks[s].ClassName = classes[c].Name
				switch classStocks[s].UsageType {
				case `Money`:
					classes[c].Money = &(classStocks[s])
				case `Consumption`:
					classes[c].Consumption = append(classes[c].Consumption, &(classStocks[s]))
				case `Sales`:
					classes[c].Sales = &(classStocks[s])
				default:
					utils.TraceErrorf("Industry stock of unknown type %s and id %d detected", industryStocks[s].UsageType, industryStocks[s].Id)
				}
			}
		}
	}

	// create direct pointers to commodities in the stock, industry and class objects.
	for com := range commodities {
		commodityId := commodities[com].Id
		for is := range industryStocks {
			if industryStocks[is].CommodityId == commodityId {
				industryStocks[is].Commodity = &commodities[com]
				industryStocks[is].CommodityName = commodities[com].Name
			}
		}
		for cs := range classStocks {
			if classStocks[cs].CommodityId == commodityId {
				classStocks[cs].Commodity = &commodities[com]
				classStocks[cs].CommodityName = commodities[com].Name
			}
		}
		// create direct pointers to commodities in the industry objects
		for ind := range industries {
			if industries[ind].Sales.CommodityId == commodityId {
				industries[ind].Commodity = &commodities[com]
			}
		}
		// create direct pointers to commodities in the class objects
		for class := range classes {
			if classes[class].Sales.CommodityId == commodityId {
				classes[class].Commodity = &commodities[com]
			}
		}

	}
}

// Fetch the tables representing one Stage in a simulation from the api server.
// The server works out who the user is from the apiKey.
// The server knows the simulationID because it knows about the user
//
// Do not convert database (Id-based) references from the AIP into pointers.
// The separate function 'ConvertStage' does this
//
// NOTE: there are many Stages for each Manager. Therefore, we do not
// access or modify the Manager from within this function. This has to
// be done externally, by setting the timestamps
//
//		user:
//	      supplies apiKey and simulationID that uniquely identify the simulation
//	      supplies the Manager for this Stage
//		returns:
//				err if anything goes wrong
func FetchStage(user *models.User) error {
	var err error
	simulationID := user.CurrentSimulationID
	utils.TraceInfof(utils.BrightCyan, "User %s is creating a new simulation with Id %d", user.UserName, simulationID)

	// Create a receiver for the data
	newStage := models.NewStage()

	// ask the API server for the data
	for key, value := range newStage {
		err = Fetch(user.ApiKey, &value)
		if err != nil {
			utils.TraceErrorf("Could not retrieve server data with key %s because of error %s", key, err.Error())
			return err
		}
	}

	// get the simulation
	simulation := user.Simulations[simulationID]

	// add the stage
	simulation.Stages = append(simulation.Stages, &newStage)

	b, _ := json.MarshalIndent(simulation, " ", " ")
	utils.TraceLogf(utils.BrightCyan, "Simulation was %s", string(b))

	return nil
}

// to be renamed 'FetchManager' once the  code is working
//
// Fetch a Manager object from the API server for the given user.
// The API server thinks this is called a 'Simulation'.
// We rename it because the API server only knows about the current stage
// but we keep a record of all stages.
//
// Contains a clumsy workaround because we haven't built a function
// to retrieve a single object. Everything else is retrieved as a table.
//
//		user: the user who will receive the new object
//	  id: the id of the required Manager
func FetchManager(user *models.User, id int) (*models.Manager, error) {
	tableContainer := models.Table{
		ApiUrl: `/simulations`,
		Table:  new([]models.Manager),
		Name:   "Simulations",
	}
	utils.TraceInfof(utils.BrightCyan, "Fetching a manager with id %d", id)

	// Get all the Managers for this user from the API
	err := Fetch(user.ApiKey, &tableContainer)
	if err != nil {
		return nil, err
	}

	table := tableContainer.Table.(*[]models.Manager)

	// Find the manager we actually want and return a pointer to it
	for i := range *table {
		if (*table)[i].Id == id {
			return &(*table)[i], nil
		}
	}
	return nil, fmt.Errorf("no simulation with id %d was found", id)
}

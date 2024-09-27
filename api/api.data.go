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

// experimental document soon
// build a tableset from the API, but don't do anything with it yet
func FetchTableSet(user *models.User) (*models.TableSet, error) {

	// Fetch all the simulations for this user (regardless of ID)
	// NOTE unlike all other objects there is only simulation object
	// and new objects are not created when the simulation steps forward.
	// HOWEVER when the simulation steps forward the States map is updated.
	// The states map is created in the 'ConvertTableSet()' function
	err := Fetch(user.ApiKey, &user.Simulations)
	if err != nil {
		return nil, err
	}

	simulations := (*user).Simulations.Table.(*[]models.Simulation)
	for i := range *simulations {
		utils.TraceInfof(utils.BrightCyan, "Setting up State table for the simulation %s", (*simulations)[i].Name)
		(*simulations)[i].States = make(map[int]string)
	}

	newTableSet := models.NewTableSet()
	for key, value := range newTableSet {
		err = Fetch(user.ApiKey, &value)
		if err != nil {
			utils.TraceErrorf("Could not retrieve server data with key %s because of error %s", key, err.Error())
		}
	}
	return &newTableSet, nil
}

// Replace Id fields with pointers. This makes for  legible code and faster access.
//
//	newTableSet: a TableSet, which has been populated by FetchTableSet
func ConvertTableSet(tableSet *models.TableSet) {
	industries := *(*tableSet)[`industries`].Table.(*[]models.Industry)
	industryStocks := *(*tableSet)[`industry stocks`].Table.(*[]models.IndustryStock)
	classes := *(*tableSet)[`classes`].Table.(*[]models.Class)
	classStocks := *(*tableSet)[`class stocks`].Table.(*[]models.ClassStock)
	commodities := *(*tableSet)[`commodities`].Table.(*[]models.Commodity)

	// set the Commodity, Sales Stock, Money stock, Industrial stocks (=Constant capital) and Social stock (=Variable Capital) of every industry
	for ind := range industries {
		industries[ind].Constant = make([]*models.IndustryStock, 0)
		for i := range industryStocks {
			if industryStocks[i].IndustryId == industries[ind].Id {
				industryStocks[i].IndustryAddress = &industries[ind]
				industryStocks[i].IndustryName = industries[ind].Name
				switch industryStocks[i].UsageType {
				case `Money`:
					industries[ind].Money = &(industryStocks[i])
				case `Production`:
					if industryStocks[i].Origin == `SOCIAL` {
						industries[ind].Variable = &(industryStocks[i])
					} else {
						industries[ind].Constant = append(industries[ind].Constant, &(industryStocks[i]))
					}
				case `Sales`:
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

// Fetches a simulation and associated tables from the api server.
// NOTE the server works out who the user is from the apiKey
// NOTE the server knows the simulationID because it knows about the user
//
// Replace Id fields with pointers. This makes for speed of access and
// visibility of code.
//
// user: supplies apiKey and simulationID that uniquely identify the simulation
//
//	 returns:
//		err if anything goes wrong
func CreateTableSet(user *models.User) error {
	// Fetch all the simulations for this user (regardless of ID)
	err := Fetch(user.ApiKey, &user.Simulations)
	if err != nil {
		return err
	}

	newTableSet, err := FetchTableSet(user)
	if err != nil {
		return err
	}

	ConvertTableSet(newTableSet)

	user.TableSets = append(user.TableSets, newTableSet)
	return nil
}

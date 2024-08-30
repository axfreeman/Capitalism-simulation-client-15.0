package models

import (
	"encoding/json"
	"gorilla-client/utils"
)

// Commonly-used Views and Tables, to pass into templates
type DisplayData struct {
	Title          string
	Simulations    *[]Simulation
	Templates      *[]Simulation
	CommodityViews *[]CommodityViewer
	IndustryViews  *[]IndustryViewer
	ClassViews     *[]ClassViewer
	IndustryStocks *[]IndustryStock
	ClassStocks    *[]ClassStock
	Trace          *[]Trace
	Count          int
	Username       string
	State          string
	Message        string
}

// Embedded data for a single commodity, to pass into templates
type CommodityData struct {
	DisplayData
	Commodity Commodity
}

// Embedded data for a single class, to pass into templates
type ClassData struct {
	DisplayData
	Class Class
}

// Embedded data for a single industry, to pass into templates
type IndustryData struct {
	DisplayData
	Industry Industry
}

// Defines Table to be synchronised with the server
//
//	ApiUrl:the endpoint on the server which fetches the Table
//	Table: one of Commodity, Industry, Class, etc etc
//	Name: convenience field for diagnostics
type Tabler struct {
	ApiUrl string      //Url to use when requesting data from the server
	Table  interface{} //All the data for one Table (eg Commodity, Industry, etc)
	Name   string      //The name of the table (for convenience, may be redundant)
}

// Contains all the tables in one stage of one simulation
// Indexed by the name of the table (commodity, industry, etc)
type TableSet map[string]Tabler

// Constructor for a TableSet. This contains all the Tables in one stage
// required for one stage of one simulation. Tables are "commodities",
// "industries", etc
func NewTableSet() TableSet {
	return map[string]Tabler{
		"commodities": {
			ApiUrl: `/commodity`,
			Table:  new([]Commodity),
			Name:   `Commodity`,
		},
		"industries": {
			ApiUrl: `/industry`,
			Table:  new([]Industry),
			Name:   `Industry`,
		},
		"classes": {
			ApiUrl: `/classes`,
			Table:  new([]Class),
			Name:   `Class`,
		},
		"industry stocks": {
			ApiUrl: `/stocks/industry`,
			Table:  new([]IndustryStock),
			Name:   `IndustryStock`,
		},
		"class stocks": {
			ApiUrl: `/stocks/class`,
			Table:  new([]ClassStock),
			Name:   `ClassStock`,
		},
		// TODO this is very verbose. Restore it later
		// "trace": {
		// 	ApiUrl: `/trace`,
		// 	Table:  new([]Trace),
		// 	Name:   `Trace`,
		// },
	}
}

// Supplies data to pass into Templates for display
//
//		u: a user
//
//		returns:
//	     if the user has no simulations, just the template list
//	     otherwise, the output data the users current simulation
func (u *User) TemplateDisplayData(message string) DisplayData {
	slist := u.SimulationsList()
	state := u.GetCurrentState()
	utils.TraceInfof(utils.BrightYellow, "Entering TemplateData for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)
	if u.CurrentSimulationID == 0 {
		utils.TraceInfo(utils.BrightYellow, "User has no simulations")
		return DisplayData{
			Title:          "Hello",
			Simulations:    nil,
			Templates:      &TemplateList,
			Count:          0,
			Username:       u.UserName,
			State:          state,
			CommodityViews: nil,
			IndustryViews:  nil,
			ClassViews:     nil,
			IndustryStocks: nil,
			ClassStocks:    nil,
			Trace:          nil,
			Message:        message,
		}
	}
	utils.TraceInfof(utils.BrightYellow, "TemplateData is retrieving data for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)
	commodityView := u.CommodityViews()
	commodityViewAsString, _ := json.MarshalIndent(commodityView, " ", " ")
	utils.TraceLogf(utils.BrightYellow, "CommodityViews returned %s", string(commodityViewAsString))
	return DisplayData{
		Title:          "Hello",
		Simulations:    slist,
		Templates:      &TemplateList,
		Count:          len(*slist),
		Username:       u.UserName,
		State:          state,
		CommodityViews: u.CommodityViews(),
		IndustryViews:  u.IndustryViews(),
		ClassViews:     u.ClassViews(),
		IndustryStocks: u.IndustryStocks(*u.GetViewedTimeStamp()),
		ClassStocks:    u.ClassStocks(*u.GetViewedTimeStamp()),
		Trace:          u.Traces(*u.GetViewedTimeStamp()),
		Message:        message,
	}
}
func (u User) CommodityDisplayData(message string, id int) CommodityData {
	return CommodityData{
		u.TemplateDisplayData(message),
		*u.Commodity(id),
	}
}

// Get a ClassData to display a single social class in the class.html template
//
//	u: the user
//	message: any message
//	id: the id of the social class item to be displayed
//
//	returns: classData which references this class, and embeds an OutputData
func (u User) ClassDisplayData(message string, id int) ClassData {
	return ClassData{
		u.TemplateDisplayData(message),
		*u.Class(id),
	}
}

// Get an IndustryData to display a single industry in the industry.html template
//
//	u: the user
//	message: any message
//	id: the id of the industry item to be displayed
//
//	returns: industryData which references this industry, and embeds an OutputData
func (u User) IndustryDisplayData(message string, id int) IndustryData {
	return IndustryData{
		u.TemplateDisplayData(message),
		*u.Industry(id),
	}
}

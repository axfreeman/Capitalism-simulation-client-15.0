package models

// Commonly-used Views and Tables, to pass into templates
type OutputData struct {
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
	OutputData
	Commodity Commodity
}

// Embedded data for a single class, to pass into templates
type ClassData struct {
	OutputData
	Class Class
}

// Embedded data for a single industry, to pass into templates
type IndustryData struct {
	OutputData
	Industry Industry
}

type TableItem interface {
	// we will fill out the methods later as this mod develops
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

package models

//TODO Table should maybe use generics instead of interface{}?
//TODO user should get the timestamp from the current simulation

// Defines Table to be synchronised with the server
//
//	ApiUrl:the endpoint on the server which fetches the Table
//	Table: one of Commodity, Industry, Class, etc etc
//	Name: convenience field for diagnostics
type Table struct {
	ApiUrl string      //Url to use when requesting data from the server
	Name   string      //The name of the table (for convenience, may be redundant)
	Table  interface{} //All the data for one Table (eg Commodity, Industry, etc)
}

// Contains all the tables in one stage of one simulation
// Indexed by the name of the table (commodity, industry, etc)
type Stage map[string]Table

// Constructor for a Stage. This contains all the Tables required for one stage of one simulation. Tables are "commodities",
// "industries", etc
func NewStage() Stage {
	return map[string]Table{
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
		"trace": {
			ApiUrl: `/trace`,
			Table:  new([]Trace),
			Name:   `Trace`,
		},
	}
}

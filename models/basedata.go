package models

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
	Name   string      //The name of the table (for convenience, may be redundant)
	Table  interface{} //All the data for one Table (eg Commodity, Industry, etc)
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
		// TODO the below produces a very verbose output. Restore it later after adding a trace viewer (with concertina)
		// "trace": {
		// 	ApiUrl: `/trace`,
		// 	Table:  new([]Trace),
		// 	Name:   `Trace`,
		// },
	}
}

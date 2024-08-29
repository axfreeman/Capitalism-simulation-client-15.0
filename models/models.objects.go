//models.objects.go
//describes the objects of the simulation itself
//functionally, these play two roles
// (1) they define how this front end communicates with the API of the backend
// (2) they define how this front end communicates with the user
// that is, the purpose is to intermediate between the simulation itself and the display of its results

package models

// A Simulation
//
// For a detailed description of the data model, consult the api.
//
// In the api, a simulation is database table which plays a key role.
// All data objects (commodity, etc) link to a record in the simulation table
//
// In this frontend, this relational structure is not used.
// Instead, the objects and the simulation are indexed by the user's simulationID
//
// UserName is a convenience field added by this frontend
// after retrieving the data from the server. It is probably redundant now.
type Simulation struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	TimeStamp            int
	ViewedTimeStamp      int
	ComparatorTimeStamp  int
	UserName             string  `json:"username"`
	State                string  `json:"state"`
	PeriodsPerYear       float32 `json:"periods_per_year"`
	PopulationGrowthRate float32 `json:"population_growth_rate"`
	InvestmentRatio      float32 `json:"investment_ratio"`
	LabourSupplyDemand   string  `json:"labour_supply_response"`
	PriceResponseType    string  `json:"price_response_type"`
	MeltResponseType     string  `json:"melt_response_type"`
	CurrencySymbol       string  `json:"currency_symbol"`
	QuantitySymbol       string  `json:"quantity_symbol"`
	Melt                 float32 `json:"melt"`
	User                 int32   `json:"user_id"`
}

type Commodity struct {
	Id                        int    `json:"id"`
	Name                      string `json:"name"`
	SimulationId              int32  `json:"simulation_id"`
	TimeStamp                 int32
	UserName                  string  `json:"username"`
	Origin                    string  `json:"origin"`
	Usage                     string  `json:"usage"`
	Size                      float32 `json:"size"`
	TotalValue                float32 `json:"total_value"`
	TotalPrice                float32 `json:"total_price"`
	UnitValue                 float32 `json:"unit_value"`
	UnitPrice                 float32 `json:"unit_price"`
	TurnoverTime              float32 `json:"turnover_time"`
	Demand                    float32 `json:"demand"`
	Supply                    float32 `json:"supply"`
	AllocationRatio           float32 `json:"allocation_ratio"`
	DisplayOrder              float32 `json:"display_order"`
	ImageName                 string  `json:"image_name"`
	Tooltip                   string  `json:"tooltip"`
	MonetarilyEffectiveDemand float32 `json:"monetarily_effective_demand"`
	InvestmentProportion      float32 `json:"investment_proportion"`
}

type Industry struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	SimulationId     int32  `json:"simulation_id"`
	TimeStamp        int
	UserName         string  `json:"username"`
	Output           string  `json:"output"`
	OutputScale      float32 `json:"output_scale"`
	OutputGrowthRate float32 `json:"output_growth_rate"`
	InitialCapital   float32 `json:"initial_capital"`
	WorkInProgress   float32 `json:"work_in_progress"`
	CurrentCapital   float32 `json:"current_capital"`
	Profit           float32 `json:"profit"`
	ProfitRate       float32 `json:"profit_rate"`
}

type Class struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	SimulationId       int32  `json:"simulation_id"`
	TimeStamp          int
	UserName           string  `json:"username"`
	Population         float32 `json:"population"`
	ParticipationRatio float32 `json:"participation_ratio"`
	ConsumptionRatio   float32 `json:"consumption_ratio"`
	Revenue            float32 `json:"revenue"`
	Assets             float32 `json:"assets"`
}

type IndustryStock struct {
	Id           int     `json:"id"`
	SimulationId int     `json:"simulation_id" `
	IndustryId   int     `json:"industry_id"`
	CommodityId  int     `json:"commodity_id" `
	UserName     string  `json:"username"`
	Name         string  `json:"name" `
	UsageType    string  `json:"usage_type" `
	Size         float32 `json:"size" `
	Value        float32 `json:"value" `
	Price        float32 `json:"price" `
	Requirement  float32 `json:"requirement" `
	Demand       float32 `json:"demand" `
}

type ClassStock struct {
	Id           int     `json:"id"`
	SimulationId int     `json:"simulation_id" `
	ClassId      int     `json:"class_id"`
	CommodityId  int     `json:"commodity_id"`
	UserName     string  `json:"username"`
	Name         string  `json:"name" `
	UsageType    string  `json:"usage_type" `
	Size         float32 `json:"size" `
	Value        float32 `json:"value" `
	Price        float32 `json:"price" `
	Requirement  float32 `json:"requirement"`
	Demand       float32 `json:"demand" `
}

// This contains a record, generated by the server, of the results of the actions
type Trace struct {
	Id            int `json:"id"`
	Simulation_id int `json:"simulation_id"`
	TimeStamp     int
	UserName      string `json:"username"`
	Level         int    `json:"level"`
	Message       string `json:"message"`
}

// This list of templates is common to all users.
// It would normally change only when the database is reset from
// immutable fixtures using Refresh().
// It is initialized when this frontend restarts.
// In future there should be some procedure for adding new templates
// or editing existing ones.
var TemplateList []Simulation

//models.objects.go
//describes the objects of the simulation itself
//functionally, these play two roles
// (1) they define how this front end communicates with the API of the backend
// (2) they define how this front end communicates with the user
// that is, the purpose is to intermediate between the simulation itself and the display of its results

// TODO Industry should have a pointer to its Commodity
// TODO IndustryStock should have a pointer to its Industry
// TODO ClassStock should have a pointer to its Class

package models

import (
	"encoding/json"
	"fmt"
)

// A Manager
//
// For a detailed description of the data model, consult the api.
//
// In the api, a simulation is a database table which plays a key role.
// All data objects (commodity, etc) link to a record in the simulation table
//
// In this frontend, this relational structure is not used.
// Instead, the objects are stored in local memory for speed.
// Each step in the simulation generates a new Stage, which represents
// this step in a form that can be passed into the Templates.
//
// NOTE UserName is a convenience field added by this frontend
// after retrieving the data from the server. It is probably redundant now.
type Manager struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	TimeStamp            int
	ViewedTimeStamp      int
	ComparatorTimeStamp  int
	UserName             string         `json:"username"`
	State                string         `json:"state"`
	States               map[int]string // use a map not a slice for efficiency
	PeriodsPerYear       float32        `json:"periods_per_year"`
	PopulationGrowthRate float32        `json:"population_growth_rate"`
	InvestmentRatio      float32        `json:"investment_ratio"`
	LabourSupplyDemand   string         `json:"labour_supply_response"`
	PriceResponseType    string         `json:"price_response_type"`
	MeltResponseType     string         `json:"melt_response_type"`
	CurrencySymbol       string         `json:"currency_symbol"`
	QuantitySymbol       string         `json:"quantity_symbol"`
	Melt                 float32        `json:"melt"`
	User                 int32          `json:"user_id"`
}

type Commodity struct {
	Id                        int     `json:"id"`
	Name                      string  `json:"name"`
	SimulationId              int32   `json:"simulation_id"`
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
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	SimulationId     int32   `json:"simulation_id"`
	UserName         string  `json:"username"`
	Output           string  `json:"output"`
	OutputScale      float32 `json:"output_scale"`
	OutputGrowthRate float32 `json:"output_growth_rate"`
	InitialCapital   float32 `json:"initial_capital"`
	WorkInProgress   float32 `json:"work_in_progress"`
	CurrentCapital   float32 `json:"current_capital"`
	Profit           float32 `json:"profit"`
	ProfitRate       float32 `json:"profit_rate"`
	Commodity        *Commodity
	Constant         []*IndustryStock
	Variable         *IndustryStock // For now, only one social stock, being Labour Power
	Money            *IndustryStock
	Sales            *IndustryStock
}

type Class struct {
	Id                 int     `json:"id"`
	Name               string  `json:"name"`
	SimulationId       int32   `json:"simulation_id"`
	Output             string  `json:"output"`
	UserName           string  `json:"username"`
	Population         float32 `json:"population"`
	ParticipationRatio float32 `json:"participation_ratio"`
	ConsumptionRatio   float32 `json:"consumption_ratio"`
	Revenue            float32 `json:"revenue"`
	Assets             float32 `json:"assets"`
	Commodity          *Commodity
	Consumption        []*ClassStock
	Money              *ClassStock
	Sales              *ClassStock
}

// Custom MarshalJSON to prevent jsonMarshal following pointer to Industry
func (i *Industry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Address string `json:"address"`
	}{
		Address: fmt.Sprintf("%v", i),
	})
}

type IndustryStock struct {
	Id              int     `json:"id"`
	SimulationId    int     `json:"simulation_id" `
	IndustryId      int     `json:"industry_id"`
	CommodityId     int     `json:"commodity_id" `
	UserName        string  `json:"username"`
	Name            string  `json:"name" `
	UsageType       string  `json:"usage_type" `
	Origin          string  `json:"origin"`
	Size            float32 `json:"size" `
	Value           float32 `json:"value" `
	Price           float32 `json:"price" `
	Requirement     float32 `json:"requirement" `
	Demand          float32 `json:"demand" `
	IndustryName    string
	CommodityName   string
	Commodity       *Commodity
	IndustryAddress *Industry
}

// Custom MarshalJSON to prevent jsonMarshal following pointer to Class
func (c *Class) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Address string `json:"address"`
	}{
		Address: fmt.Sprintf("%v", c),
	})
}

type ClassStock struct {
	Id            int     `json:"id"`
	SimulationId  int     `json:"simulation_id" `
	ClassId       int     `json:"class_id"`
	CommodityId   int     `json:"commodity_id"`
	UserName      string  `json:"username"`
	Name          string  `json:"name" `
	UsageType     string  `json:"usage_type" `
	Size          float32 `json:"size" `
	Value         float32 `json:"value" `
	Price         float32 `json:"price" `
	Requirement   float32 `json:"requirement"`
	Demand        float32 `json:"demand" `
	ClassName     string
	CommodityName string
	Commodity     *Commodity
	ClassAddress  *Class
}

// This list of templates is common to all users.
// It would normally change only when the database is reset from
// immutable fixtures using Refresh().
// It is initialized when this frontend restarts.
// In future there should be some procedure for adding new templates
// or editing existing ones.
var TemplateList []Manager

// A default Industry_stock returned if any condition is not met (that is, if the predicated stock does not exist)
// Used to signal to the user that there has been a programme error
var NotFoundIndustryStock = IndustryStock{
	Id:           0,
	SimulationId: 0,
	CommodityId:  0,
	Name:         "NOT FOUND",
	UsageType:    "PROGRAMME ERROR",
	Size:         -1,
	Value:        -1,
	Price:        -1,
	Requirement:  -1,
	Demand:       -1,
}

// A default Industry_stock returned if any condition is not met (that is, if the predicated stock does not exist)
// Used to signal to the user that there has been a programme error
var NotFoundClassStock = ClassStock{
	Id:           0,
	SimulationId: 0,
	CommodityId:  0,
	Name:         "NOT FOUND",
	UsageType:    "PROGRAMME ERROR",
	Size:         -1,
	Value:        -1,
	Price:        -1,
	Demand:       -1,
}

var NotFoundCommodity = Commodity{
	Id:                        0,
	Name:                      "NOT FOUND",
	SimulationId:              0,
	Origin:                    "UNDEFINED",
	Usage:                     "UNDEFINED",
	Size:                      0,
	TotalValue:                0,
	TotalPrice:                0,
	UnitValue:                 0,
	UnitPrice:                 0,
	TurnoverTime:              0,
	Demand:                    0,
	Supply:                    0,
	AllocationRatio:           0,
	DisplayOrder:              0,
	ImageName:                 "UNDEFINED",
	Tooltip:                   "UNDEFINED",
	MonetarilyEffectiveDemand: 0,
	InvestmentProportion:      0,
}

// A default Class returned if any condition is not met (that is, if the class does not exist)
// Used to signal to the user that there has been a programme error
var NotFoundClass = Class{
	Id:                 0,
	Name:               "NOT FOUND",
	SimulationId:       0,
	UserName:           "UNDEFINED",
	Population:         0,
	ParticipationRatio: 0,
	ConsumptionRatio:   0,
	Revenue:            0,
	Assets:             0,
}

var NotFoundIndustry = Industry{
	Id:           0,
	Name:         "NOT FOUND",
	SimulationId: 0,
	UserName:     "UNDEFINED",
}

func (c Commodity) GetId() int {
	return c.Id
}

func (i Industry) GetId() int {
	return i.Id
}

func (c Class) GetId() int {
	return c.Id
}

func (is IndustryStock) GetId() int {
	return is.Id
}

func (cs ClassStock) GetId() int {
	return cs.Id
}

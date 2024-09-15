package models

import (
	"fmt"
	"html/template"
	"strconv"
)

//METHODS OF INDUSTRIES
//METHODS OF INDUSTRIES

// crude searches without database implementation
// justified because we can avoid the complications of a database implementation
// and the size of the tables is not large, because they are provided on a per-user basis
// However as the simulations get large, this may become more problematic (let's find out pragmatically)
// In that case some more sophisticated system, such as a local database, may be needed
// A simple solution would be to add direct links to related objects in the models
// perhaps populated by an asynchronous process in the background
// TODO remove boilerplate with generic search OR use maps instead of slices.
// TODO a generic function for commodity is available. Problem is to make it generic.

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

func (p Pair) Format() template.HTML {
	var htmlString string
	if p.Viewed == p.Compared {
		htmlString = fmt.Sprintf("<td style=\"text-align:right\">%0.2f</td>", p.Viewed)
	} else {
		htmlString = fmt.Sprintf("<td style=\"text-align:right; color:red\">%0.2f</td>", p.Viewed)
	}
	return template.HTML(htmlString)
}

func (p Pair) FormatRounded() template.HTML {
	var htmlString string
	if p.Viewed == p.Compared {
		htmlString = fmt.Sprintf("<td style=\"text-align:center\">%.0f</td>", p.Viewed)
	} else {
		htmlString = fmt.Sprintf("<td style=\"text-align:center; color:red\">%.0f</td>", p.Viewed)
	}
	return template.HTML(htmlString)
}

// returns the Labour Power stock of the given industry
func (industry Industry) VariableCapital() IndustryStock {
	s := industry.Variable
	return *s
}

// returns the commodity that an industry produces
func (industry Industry) OutputCommodity() *Commodity {
	return industry.Sales.Commodity()
}

// returns the sales stock of the given class
func (class Class) MoneyStock() ClassStock {
	stockList := *class.Stocks
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.ClassId == class.Id) && (s.UsageType == `Money`) {
			return *s
		}
	}
	return NotFoundClassStock
}

// returns the sales stock of the given class
func (class Class) SalesStock() ClassStock {
	stockList := *class.Stocks
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.ClassId == class.Id) && (s.UsageType == `Sales`) {
			return *s
		}
	}
	return NotFoundClassStock
}

// returns the consumption stock of the given class
// under development - at present assumes there is only one
func (class Class) ConsumerGood() ClassStock {
	stockList := *class.Stocks
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.ClassId == class.Id) && (s.UsageType == `Consumption`) {
			return *s
		}
	}
	return NotFoundClassStock
}

// Get all consumption stocks of a given class. Under Development
//
//	returns:
//	 slice of stocks of usageType "Consumption" owned by the class
func (class Class) ConsumerGoods() *[]ClassStock {
	user := LoggedInUsers[class.UserName]
	partialStockList := make([]ClassStock, 0)

	fullStockList := user.ClassStocks(*user.GetTimeStamp())
	for i := range *fullStockList {
		s := (*fullStockList)[i]
		if s.UsageType == `Consumption` && s.ClassId == class.Id {
			partialStockList = append(partialStockList, (*fullStockList)[i])
		}
	}
	return &partialStockList
}

// return the name of the commodity that the given Industry_Stock consists of
//
//	Deprecated - just use s.Commodity().Name
func (s IndustryStock) CommodityName() string {
	return s.Commodity().Name
}

// return the Commodity that the given stock consists of
func (s IndustryStock) Commodity() *Commodity {
	return LoggedInUsers[s.UserName].Commodity(s.CommodityId)
}

// return the Commodity that the given stock consists of
func (s ClassStock) Commodity() *Commodity {
	return LoggedInUsers[s.UserName].Commodity(s.CommodityId)
}

// (Experimental) Creates a url to link to this simulation, to be used in templates such as dashboard
// In this way all the URL naming is done in native Golang, not in the template
// We may also use such methods in the Trace function to improve usability
func (s Simulation) Link() string {
	return `/user/create/` + strconv.Itoa(s.Id)
}

// fetches the industry that owns this industry stock
// If it has none (an error, but we need to diagnose it) return nil.
func (s IndustryStock) Industry() *Industry {
	return LoggedInUsers[s.UserName].Industry(s.IndustryId)
}

// fetches the class that owns this Class_stock
// If it has none (an error, but we need to diagnose it) return nil.
func (s ClassStock) Class() *Class {
	return LoggedInUsers[s.UserName].Class(s.ClassId)
}

// Convenience named functions for PopulateView
// TODO rationalise - still in development
// Deprecated - get rid ASAP

// Named as convenience for the PopulateView function to use
func (i Industry) ConstantCapitalSize() float32 {
	return i.Constant[0].Size
}

// Named as convenience for the PopulateView function to use
func (i Industry) ConstantCapitalValue() float32 {
	return i.Constant[0].Value
}

// Named as convenience for the PopulateView function to use
func (i Industry) ConstantCapitalPrice() float32 {
	return i.Constant[0].Price
}

// Named as convenience for the PopulateView function to use
func (i Industry) VariableCapitalSize() float32 {
	return i.Variable.Size
}

// Named as convenience for the PopulateView function to use
func (i Industry) VariableCapitalValue() float32 {
	return i.Variable.Value
}

// Named as convenience for the PopulateView function to use
func (i Industry) VariableCapitalPrice() float32 {
	return i.Variable.Price
}

// Named as convenience for the PopulateView function to use
func (i Industry) SalesStockSize() float32 {
	return i.Sales.Size
}

// Named as convenience for the PopulateView function to use
func (i Industry) SalesStockValue() float32 {
	return i.Sales.Value
}

// Named as convenience for the PopulateView function to use
func (i Industry) SalesStockPrice() float32 {
	return i.Sales.Price
}

// Named as convenience for the PopulateView function to use
func (i Class) MoneyStockSize() float32 {
	return i.MoneyStock().Size
}

// Named as convenience for the PopulateView function to use
func (i Class) MoneyStockValue() float32 {
	return i.MoneyStock().Value
}

// Named as convenience for the PopulateView function to use
func (i Class) MoneyStockPrice() float32 {
	return i.MoneyStock().Price
}

// Named as convenience for the PopulateView function to use
func (i Class) SalesStockSize() float32 {
	return i.SalesStock().Size
}

// Named as convenience for the PopulateView function to use
func (i Class) SalesStockValue() float32 {
	return i.SalesStock().Value
}

// Named as convenience for the PopulateView function to use
func (i Class) SalesStockPrice() float32 {
	return i.SalesStock().Price
}

// Named as convenience for the PopulateView function to use
func (i Class) ConsumerGoodSize() float32 {

	return i.ConsumerGood().Size
}

// Named as convenience for the PopulateView function to use
func (i Class) ConsumerGoodValue() float32 {
	return i.ConsumerGood().Value
}

// Named as convenience for the PopulateView function to use
func (i Class) ConsumerGoodPrice() float32 {
	return i.ConsumerGood().Price
}

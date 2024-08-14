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
	TimeStamp:                 0,
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
	TimeStamp:          0,
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
	TimeStamp:    0,
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

// returns the money stock of the given industry
func (industry Industry) MoneyStock(timeStamp int) IndustryStock {
	username := industry.UserName
	stockList := *LoggedInUsers[username].IndustryStocks(timeStamp)
	for i := 0; i < len(stockList); i++ {
		s := stockList[i]
		if (s.IndustryId == industry.Id) && (s.UsageType == `Money`) {
			return s
		}
	}
	return NotFoundIndustryStock
}

// returns the sales stock of the given industry
func (industry Industry) SalesStock(timeStamp int) IndustryStock {
	username := industry.UserName
	stockList := *LoggedInUsers[username].IndustryStocks(timeStamp)
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.IndustryId == industry.Id) && (s.UsageType == `Sales`) {
			return *s
		}
	}
	return NotFoundIndustryStock
}

// returns the Labour Power stock of the given industry
// bit of a botch to use the name of the commodity as a search term
func (industry Industry) VariableCapital(timeStamp int) IndustryStock {
	username := industry.UserName
	stockList := *LoggedInUsers[username].IndustryStocks(timeStamp)
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.IndustryId == industry.Id) && (s.UsageType == `Production`) && (s.CommodityName() == "Labour Power") {
			return *s
		}
	}
	return NotFoundIndustryStock
}

// returns the commodity that an industry produces
func (industry Industry) OutputCommodity(timeStamp int) *Commodity {
	return industry.SalesStock(timeStamp).Commodity()
}

// return the productive capital stock of the given industry
// under development - at present assumes there is only one
func (industry Industry) ConstantCapital(timeStamp int) IndustryStock {
	username := industry.UserName
	stockList := *LoggedInUsers[username].IndustryStocks(timeStamp)
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.IndustryId == industry.Id) && (s.UsageType == `Production`) && (s.CommodityName() == "Means of Production") {
			return *s
		}
	}
	return NotFoundIndustryStock
}

// returns all the constant capitals of a given industry.
// Under development.
// func (industry Industry) ConstantCapitals() []Stock {
// 	return &stocks [Programming error here]
// }

// returns the sales stock of the given class
func (class Class) MoneyStock(timeStamp int) ClassStock {
	username := class.UserName
	stockList := *LoggedInUsers[username].ClassStocks(timeStamp)

	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.ClassId == class.Id) && (s.UsageType == `Money`) {
			return *s
		}
	}
	return NotFoundClassStock
}

// returns the sales stock of the given class
func (class Class) SalesStock(timeStamp int) ClassStock {
	username := class.UserName
	stockList := *LoggedInUsers[username].ClassStocks(timeStamp)
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
func (class Class) ConsumerGood(timeStamp int) ClassStock {
	username := class.UserName
	stockList := *LoggedInUsers[username].ClassStocks(timeStamp)

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

// under development
// will eventually be parameterised to yield value, price or quantity depending on a 'display' parameter
func (stock IndustryStock) DisplaySize(mode string) float32 {
	switch mode {
	case `prices`:
		return stock.Size
	case `quantities`:
		return stock.Size // switch in price once this is in the model
	default:
		panic(`unknown display mode requested`)
	}
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

package models

import (
	"encoding/json"
	"gorilla-client/utils"
)

// List of the user's Simulations.
//
//	u: the user
//	returns:
//	 Slice of SimulationsList
//	 If the user has no simulations, an empty slice
func (u User) SimulationsList() *[]Simulation {
	list := u.Simulations.Table.(*[]Simulation)
	if len(*list) == 0 {
		var fakeList []Simulation = *new([]Simulation)
		return &fakeList
	}
	return list
}

// Create a CommodityView for display in a template
// taking data from two Commodity objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func NewCommodityView(v *Commodity, c *Commodity) *CommodityView {
	newCommodityView := CommodityView{
		Id:                        v.Id,
		Name:                      v.Name,
		Origin:                    v.Origin,
		Usage:                     v.Usage,
		Size:                      Pair{Viewed: v.Size, Compared: c.Size},
		TotalValue:                Pair{Viewed: (v.TotalValue), Compared: (c.TotalValue)},
		TotalPrice:                Pair{Viewed: (v.TotalPrice), Compared: (c.TotalPrice)},
		UnitValue:                 Pair{Viewed: (v.UnitValue), Compared: (c.UnitValue)},
		UnitPrice:                 Pair{Viewed: (v.UnitPrice), Compared: (c.UnitPrice)},
		TurnoverTime:              Pair{Viewed: v.TurnoverTime, Compared: c.TurnoverTime},
		Demand:                    Pair{Viewed: v.Demand, Compared: c.Demand},
		Supply:                    Pair{Viewed: v.Supply, Compared: c.Supply},
		AllocationRatio:           Pair{Viewed: v.AllocationRatio, Compared: c.AllocationRatio},
		MonetarilyEffectiveDemand: v.MonetarilyEffectiveDemand,
		InvestmentProportion:      v.InvestmentProportion,
	}
	return &newCommodityView
}

func NewCommodityViews(v *[]Commodity, c *[]Commodity) *[]CommodityView {
	var newViews = make([]CommodityView, len(*v))
	for i := range *v {
		newView := NewCommodityView(&(*v)[i], &(*c)[i])
		newViews[i] = *newView
	}
	return &newViews
}

// Create an IndustryView for display in a template
// taking data from two Industry objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'.
//
// We load up all the 'calculated magnitudes' such as ConstantCapitalValue
// so that when the user is scanning the simulation results, the retrieval
// time is as small as can be.
//
//		v the viewed industry
//		c the comparator industry
//		vTimeStamp the viewed TimeStamp
//		cTimeStamp the comparator TimeStamp
//
//	 Returns: a new IndustryView
func NewIndustryView(vTimeStamp int, cTimeStamp int, v *Industry, c *Industry) *IndustryView {
	newView := IndustryView{
		Id:                   v.Id,
		Name:                 v.Name,
		Output:               v.Output,
		OutputCommodityId:    v.OutputCommodity(vTimeStamp).Id, // TODO check if this causes any problems
		OutputScale:          Pair{Viewed: (v.OutputScale), Compared: (c.OutputScale)},
		OutputGrowthRate:     Pair{Viewed: (v.OutputGrowthRate), Compared: (c.OutputGrowthRate)},
		InitialCapital:       Pair{Viewed: (v.InitialCapital), Compared: (c.InitialCapital)},
		WorkInProgress:       Pair{Viewed: (v.WorkInProgress), Compared: (c.WorkInProgress)},
		CurrentCapital:       Pair{Viewed: (v.CurrentCapital), Compared: (c.CurrentCapital)},
		ConstantCapitalSize:  Pair{Viewed: (v.ConstantCapital(vTimeStamp).Size), Compared: (c.ConstantCapital(cTimeStamp).Size)},
		ConstantCapitalValue: Pair{Viewed: (v.ConstantCapital(vTimeStamp).Value), Compared: (c.ConstantCapital(cTimeStamp).Value)},
		ConstantCapitalPrice: Pair{Viewed: (v.ConstantCapital(vTimeStamp).Price), Compared: (c.ConstantCapital(cTimeStamp).Price)},
		VariableCapitalSize:  Pair{Viewed: (v.VariableCapital(vTimeStamp).Size), Compared: (c.VariableCapital(cTimeStamp).Size)},
		VariableCapitalValue: Pair{Viewed: (v.VariableCapital(vTimeStamp).Value), Compared: (c.VariableCapital(cTimeStamp).Value)},
		VariableCapitalPrice: Pair{Viewed: (v.VariableCapital(vTimeStamp).Price), Compared: (c.VariableCapital(cTimeStamp).Price)},
		MoneyStockSize:       Pair{Viewed: (v.MoneyStock(vTimeStamp).Size), Compared: (c.MoneyStock(cTimeStamp).Size)},
		MoneyStockValue:      Pair{Viewed: (v.MoneyStock(vTimeStamp).Value), Compared: (c.MoneyStock(cTimeStamp).Value)},
		MoneyStockPrice:      Pair{Viewed: (v.MoneyStock(vTimeStamp).Price), Compared: (c.MoneyStock(cTimeStamp).Price)},
		SalesStockSize:       Pair{Viewed: (v.SalesStock(vTimeStamp).Size), Compared: (c.SalesStock(cTimeStamp).Size)},
		SalesStockValue:      Pair{Viewed: (v.SalesStock(vTimeStamp).Value), Compared: (c.SalesStock(cTimeStamp).Value)},
		SalesStockPrice:      Pair{Viewed: (v.SalesStock(vTimeStamp).Price), Compared: (c.SalesStock(cTimeStamp).Price)},
		Profit:               Pair{Viewed: (v.Profit), Compared: (c.Profit)},
		ProfitRate:           Pair{Viewed: (v.ProfitRate), Compared: (c.ProfitRate)},
	}

	// newViewAsString, _ := json.MarshalIndent(newView, " ", " ")
	// utils.Trace(utils.BrightCyan, "  Industry view is\n"+string(newViewAsString))
	return &newView
}

// Creates a slice of IndustryViews which provide pairs
// of Industry objects corresponding to two points in the
// simulation - viewed and compared.
// This allows us to display, visually, changes that have
// taken place between any two steps in the simulation.
//
//	vTimeStamp: the viewed TimeStamp.
//	cTimeStamp: the comparator TimeStamp.
//	v: a snapshot industry array (Department I, Department II, etc) at state vTimeStamp.
//	v: a snapshot industry array (Department I, Department II, etc) at state cTimeStamp.
//	returns: a slice of IndustryViews.
func NewIndustryViews(vTimeStamp int, cTimeStamp int, v *[]Industry, c *[]Industry) *[]IndustryView {
	var newViews = make([]IndustryView, len(*v))
	for i := range *v {
		newView := NewIndustryView(vTimeStamp, cTimeStamp, &(*v)[i], &(*c)[i])
		newViews[i] = *newView
	}
	return &newViews
}

func NewClassView(vTimeStamp int, cTimeStamp int, v *Class, c *Class) *ClassView {
	newView := ClassView{
		Id:                    v.Id,
		Name:                  v.Name,
		SimulationId:          v.SimulationId,
		TimeStamp:             v.TimeStamp,
		UserName:              v.UserName,
		Population:            Pair{Viewed: (v.Population), Compared: (c.Population)},
		ParticipationRatio:    v.ParticipationRatio,
		ConsumptionRatio:      v.ConsumptionRatio,
		Revenue:               Pair{Viewed: (v.Revenue), Compared: (c.Revenue)},
		Assets:                Pair{Viewed: (v.Assets), Compared: (c.Assets)},
		ConsumptionStockSize:  Pair{Viewed: (v.ConsumerGood(vTimeStamp).Size), Compared: (c.ConsumerGood(vTimeStamp).Size)},
		ConsumptionStockValue: Pair{Viewed: (v.ConsumerGood(vTimeStamp).Value), Compared: (c.ConsumerGood(vTimeStamp).Value)},
		ConsumptionStockPrice: Pair{Viewed: (v.ConsumerGood(vTimeStamp).Price), Compared: (c.ConsumerGood(vTimeStamp).Price)},
		MoneyStockSize:        Pair{Viewed: (v.MoneyStock(vTimeStamp).Size), Compared: (c.MoneyStock(vTimeStamp).Size)},
		MoneyStockValue:       Pair{Viewed: (v.MoneyStock(vTimeStamp).Value), Compared: (c.MoneyStock(vTimeStamp).Value)},
		MoneyStockPrice:       Pair{Viewed: (v.MoneyStock(vTimeStamp).Price), Compared: (c.MoneyStock(vTimeStamp).Price)},
		SalesStockSize:        Pair{Viewed: (v.SalesStock(vTimeStamp).Size), Compared: (c.SalesStock(vTimeStamp).Size)},
		SalesStockValue:       Pair{Viewed: (v.SalesStock(vTimeStamp).Value), Compared: (c.SalesStock(vTimeStamp).Value)},
		SalesStockPrice:       Pair{Viewed: (v.SalesStock(vTimeStamp).Price), Compared: (c.SalesStock(vTimeStamp).Price)},
	}
	return &newView
}

// Creates a slice of ClassViews which provide pairs
// of Class objects corresponding to two points in the
// simulation - viewed and compared.
// This allows us to display, visually, changes that have
// taken place between any two steps in the simulation.
//
//	vTimeStamp: the viewed TimeStamp.
//	cTimeStamp: the comparator TimeStamp.
//	v: a snapshot Class array (Department I, Department II, etc) at state vTimeStamp.
//	v: a snapshot Class array (Department I, Department II, etc) at state cTimeStamp.
//	returns: a slice of ClassViews.
func NewClassViews(vTimeStamp int, cTimeStamp int, v *[]Class, c *[]Class) *[]ClassView {
	var newViews = make([]ClassView, len(*v))
	for i := range *v {
		newView := NewClassView(vTimeStamp, cTimeStamp, &(*v)[i], &(*c)[i])
		newViews[i] = *newView
	}
	return &newViews
}

// supplies outputData to be passed into Templates for display
//
//		u: a user
//
//		returns:
//	     if the user has no simulations, just the template list
//	     otherwise, the output data the users current simulation
func (u *User) TemplateData(message string) OutputData {
	slist := u.SimulationsList()
	state := u.GetCurrentState()
	utils.TraceInfof(utils.BrightYellow, "Entering TemplateData for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)
	if u.CurrentSimulationID == 0 {
		utils.TraceInfo(utils.BrightYellow, "User has no simulations")
		return OutputData{
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
	return OutputData{
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
func (u User) OutputCommodityData(message string, id int) CommodityData {
	return CommodityData{
		u.TemplateData(message),
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
func (u User) OutputClassData(message string, id int) ClassData {
	return ClassData{
		u.TemplateData(message),
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
func (u User) OutputIndustryData(message string, id int) IndustryData {
	return IndustryData{
		u.TemplateData(message),
		*u.Industry(id),
	}
}

func (u *User) LogTemplateData() string {
	output := u.TemplateData("hello")
	outputAsString, _ := json.MarshalIndent(output, " ", " ")
	return string(outputAsString)
}

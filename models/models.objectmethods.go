// methods.simulation.go
// class methods of the objects specified in models.simulation.go
package models

import (
	"gorilla-client/utils"
)

// Retrieve the state of the current simulation
//
//		returns:
//	   if successful, one of "DEMAND", "TRADE",  ...(the stages of the cycle)
//	   if unsuccessful "UNKNOWN"
func (u User) GetCurrentState() string {
	var s *Simulation
	if s = u.Simulation(u.CurrentSimulationID); s != nil {
		return s.State
	}
	return "UNKNOWN"
}

// Set the state of the current simulation
//
//	new_state: one of "DEMAND", "TRADE",  ... (the stages of the cycle)
//	returns: does not report any error. It probably should.
func (u User) SetCurrentState(new_state string) {
	utils.TraceInfof(utils.Purple, "Set the state of simulation with id %d to %s", u.CurrentSimulationID, new_state)
	var s *Simulation
	if s = u.Simulation(u.CurrentSimulationID); s == nil {
		utils.TraceError("Attempt to access non-existent simulation")
		return
	}
	s.State = new_state
}

func (u User) Commodities() *[]Commodity {
	return (*u.TableSets[*u.GetViewedTimeStamp()])["commodities"].Table.(*[]Commodity)
}

func (u User) CommodityViews() *[]CommodityView {
	utils.TraceLogf(utils.BrightRed, "Entered CommodityViews with time stamp %d and comparator %d", *u.GetViewedTimeStamp(), *u.GetComparatorTimeStamp())
	v := (*u.TableSets[*u.GetViewedTimeStamp()])["commodities"].Table.(*[]Commodity)
	c := (*u.TableSets[*u.GetComparatorTimeStamp()])["commodities"].Table.(*[]Commodity)
	return VeryNewCommodityViews(v, c)
}

func (u User) Industries() *[]Industry {
	return (*u.TableSets[*u.GetViewedTimeStamp()])["industries"].Table.(*[]Industry)
}

func (u User) IndustryViews() *[]IndustryView {
	v := (*u.TableSets[*u.GetViewedTimeStamp()])["industries"].Table.(*[]Industry)
	c := (*u.TableSets[*u.GetComparatorTimeStamp()])["industries"].Table.(*[]Industry)
	return VeryNewIndustryViews(v, c)
}

func (u User) Classes() *[]Class {
	return (*u.TableSets[*u.GetViewedTimeStamp()])["classes"].Table.(*[]Class)
}

func (u User) ClassViews() *[]ClassView {
	v := (*u.TableSets[*u.GetViewedTimeStamp()])["classes"].Table.(*[]Class)
	c := (*u.TableSets[*u.GetComparatorTimeStamp()])["classes"].Table.(*[]Class)
	return VeryNewClassViews(v, c)
}

// Wrapper for the IndustryStockList
func (u User) IndustryStocks(timeStamp int) *[]IndustryStock {
	return (*u.TableSets[timeStamp])["industry stocks"].Table.(*[]IndustryStock)
}

// Wrapper for the ClassStockList
func (u User) ClassStocks(timeStamp int) *[]ClassStock {
	return (*u.TableSets[timeStamp])["class stocks"].Table.(*[]ClassStock)
}

// Wrapper for the TraceList
func (u User) Traces(timeStamp int) *[]Trace {
	if len(u.TableSets) == 0 {
		return nil
	}
	table, ok := (*u.TableSets[timeStamp])["trace"]
	if !ok {
		return nil
	}
	return table.Table.(*[]Trace)
}

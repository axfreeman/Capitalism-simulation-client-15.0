package models

import (
	"gorilla-client/utils"
)

// TODO rationalise traces and simulations

// Retrieve the current state of the current simulation
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

// Set the state of the current simulation. Make a record of this state
// in the 'States' map so it can be retrieved as a comparator state
//
//	new_state: one of "DEMAND", "TRADE",  ... (the stages of the cycle)
//	returns: does not report any error. It probably should.
func (u User) SetCurrentState(new_state string) {
	utils.TraceInfof(utils.Green, "Set the state of simulation with id %d to %s", u.CurrentSimulationID, new_state)
	var s *Simulation
	if s = u.Simulation(u.CurrentSimulationID); s == nil {
		utils.TraceError("Attempt to access non-existent simulation")
		return
	}
	s.State = new_state
	s.States[u.TimeStamp] = new_state
	utils.TraceInfof(utils.Green, "Setting new state %s. States map now has %d elements", new_state, len(s.States))
	for i := range s.States {
		utils.TraceInfof(utils.BrightGreen, "State %d is %s", i, s.States[i])
	}
}

type Object interface {
	Commodity | Industry | Class | IndustryStock | ClassStock | Simulation | Trace
	GetId() int
}

func ViewedObjects[T Object](u User, objectType string) *[]T {
	return (*u.TableSets[*u.GetViewedTimeStamp()])[objectType].Table.(*[]T)
}

func ComparedObjects[T Object](u User, objectType string) *[]T {
	return (*u.TableSets[*u.GetComparatorTimeStamp()])[objectType].Table.(*[]T)
}

func ViewedObject[T Object](u User, objectType string, id int) *T {
	objectList := (*u.TableSets[*u.GetViewedTimeStamp()])[objectType].Table.(*[]T)
	for i := 0; i < len(*objectList); i++ {
		o := (*objectList)[i]
		if id == o.GetId() {
			return &o
		}
	}
	return nil
}

func ComparedObject[T Object](u User, objectType string, id int) *T {
	objectList := (*u.TableSets[*u.GetComparatorTimeStamp()])[objectType].Table.(*[]T)
	for i := 0; i < len(*objectList); i++ {
		o := (*objectList)[i]
		if id == o.GetId() {
			return &o
		}
	}
	return nil
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

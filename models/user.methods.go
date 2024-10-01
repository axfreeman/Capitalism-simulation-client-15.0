package models

import (
	"gorilla-client/utils"
	"log"
)

// TODO handle programme errors more systematically

// Retrieve the current simulation
func (u *User) GetCurrentSimulation() *Simulation {
	s, ok := u.Simulations[u.CurrentSimulationID]
	if !ok {
		utils.TraceErrorf("could not retrieve the simulation with id %d", u.CurrentSimulationID)
		log.Fatalf("could not retrieve the simulation with id %d", u.CurrentSimulationID) //TODO very temporary
	}
	return s
}

// Retrieve the current stage of the simulation.
func (u *User) GetCurrentStage() *Stage {
	manager := &u.GetCurrentSimulation().Manager
	return u.GetCurrentSimulation().Stages[manager.TimeStamp]
}

// Retrieve the viewed stage of the simulation. The Comparator Stage
// is normally one step behind, but this may change in the future
func (u *User) GetViewedStage() *Stage {
	manager := &u.GetCurrentSimulation().Manager
	return u.GetCurrentSimulation().Stages[manager.ViewedTimeStamp]
}

// Retrieve the comparator stage of the simulation. This is normally one
// step behind the ViewedStage, but we may change this later, for example
// to view the difference between one period and the next.
func (u *User) GetComparatorStage() *Stage {
	manager := &u.GetCurrentSimulation().Manager
	return u.GetCurrentSimulation().Stages[manager.ComparatorTimeStamp]
}

// Retrieve the current state of the current simulation
//
//		returns:
//	   if successful, one of "DEMAND", "TRADE",  ...(the stages of the cycle)
//	   if unsuccessful "UNKNOWN"
func (u User) GetCurrentState() string {
	manager := &u.GetCurrentSimulation().Manager
	return manager.State
}

// Set the state of the current simulation. Make a record of this state
// in the 'States' map so it can be retrieved as a comparator state
//
//	new_state: one of "DEMAND", "TRADE",  ... (the stages of the cycle)
//	returns: does not report any error. It probably should.
func (u User) SetCurrentState(new_state string) {
	utils.TraceInfof(utils.Green, "Set the state of simulation with id %d to %s", u.CurrentSimulationID, new_state)
	s, ok := u.Simulations[u.CurrentSimulationID]
	if !ok {
		utils.TraceErrorf("attempt to set the state of a non-existent simulation using id %d", u.CurrentSimulationID)
		return
	}
	m := &s.Manager
	m.State = new_state
	m.States[u.TimeStamp] = new_state
	utils.TraceInfof(utils.Green, "Setting new state %s. States map now has %d elements", new_state, len(m.States))
	for i := range m.States {
		utils.TraceInfof(utils.BrightGreen, "State %d is %s", i, m.States[i])
	}
}

type Object interface {
	Commodity | Industry | Class | IndustryStock | ClassStock | Manager | Trace
	GetId() int
}

func ViewedObjects[T Object](u User, objectType string) *[]T {
	return (*u.GetViewedStage())[objectType].Table.(*[]T)
}

func ComparedObjects[T Object](u User, objectType string) *[]T {
	return (*u.GetComparatorStage())[objectType].Table.(*[]T)
}

func ViewedObject[T Object](u User, objectType string, id int) *T {
	objectList := (*u.GetViewedStage())[objectType].Table.(*[]T)
	for i := 0; i < len(*objectList); i++ {
		o := (*objectList)[i]
		if id == o.GetId() {
			return &o
		}
	}
	return nil
}

func ComparedObject[T Object](u User, objectType string, id int) *T {
	objectList := (*u.GetComparatorStage())[objectType].Table.(*[]T)
	for i := 0; i < len(*objectList); i++ {
		o := (*objectList)[i]
		if id == o.GetId() {
			return &o
		}
	}
	return nil
}

// TODO deprecate and remove?

// List of the user's Simulations.
//
//	u: the user
//	returns:
//	 Slice of SimulationsList
//	 If the user has no simulations, an empty slice
// func (u User) SimulationsList() *[]Manager {
// 	list := u.Managers.Table.(*[]Manager)
// 	if len(*list) == 0 {
// 		var fakeList []Manager = *new([]Manager)
// 		return &fakeList
// 	}
// 	return list
// }

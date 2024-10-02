package models

import (
	"gorilla-client/utils"
	"log"
)

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
	// fmt.Printf("GetViewedStage called with stamp %d\n", manager.ViewedTimeStamp)
	return u.GetCurrentSimulation().Stages[manager.ViewedTimeStamp]
}

// Retrieve the comparator stage of the simulation. This is normally one
// step behind the ViewedStage, but we may change this later, for example
// to view the difference between one period and the next.
func (u *User) GetComparatorStage() *Stage {
	manager := &u.GetCurrentSimulation().Manager
	// fmt.Printf("GetComparatorStage called with stamp %d\n", manager.ComparatorTimeStamp)
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

// Retrieve the viewed state of the current simulation
//
//		returns:
//	   if successful, one of "DEMAND", "TRADE",  ...(the stages of the cycle)
//	   if unsuccessful "UNKNOWN"
func (u User) ViewedState() string {
	manager := &u.GetCurrentSimulation().Manager
	viewedStamp := manager.ViewedTimeStamp
	viewedState := manager.States[viewedStamp]
	return viewedState
}

// Retrieve the comparator state of the current simulation
//
//		returns:
//	   if successful, one of "DEMAND", "TRADE",  ...(the stages of the cycle)
//	   if unsuccessful "UNKNOWN"
func (u User) ComparatorState() string {
	manager := &u.GetCurrentSimulation().Manager
	viewedStamp := manager.ComparatorTimeStamp
	viewedState := manager.States[viewedStamp]
	return viewedState
}

// Set the state of the current simulation. Make a record of this state
// in the 'States' map so it can be retrieved as a comparator state
//
//	new_state: one of "DEMAND", "TRADE",  ... (the stages of the cycle)
//	returns: does not report any error. It probably should.
func (u User) SetCurrentState(new_state string) {
	s := u.GetCurrentSimulation()
	m := &s.Manager
	utils.TraceInfof(utils.Green, "Set the state of simulation %d to %s at timestamp %d", u.CurrentSimulationID, new_state, m.TimeStamp)
	m.State = new_state
	m.States[m.TimeStamp] = new_state
	utils.TraceInfof(utils.Green, "Setting new state %s. States map now has %d elements", new_state, len(m.States))
	for i := range m.States {
		utils.TraceInfof(utils.BrightGreen, "State %d is %s", i, m.States[i])
	}
}

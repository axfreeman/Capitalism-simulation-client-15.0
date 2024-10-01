package models

//TODO replace object finders with generics and maybe they don't belong in this file anyhow

import (
	"encoding/json"
)

// A record describing what page the user was visiting
// together with the information needed to display the page
type CurrentPageType struct {
	Url string
	Id  int
}

// A Simulation object completely describes one simulation
// One Simulation contains:
//
//	Manager, which controls access to the data
//	Stages, a slice of Stage, which contains the data
//
// A Stage is one step in the simulation (Trade, Produce, etc)
// A Stage contains all Tables of this step (Commoditis, Industries, etc)
// Each Table contains all the Objects of a single type. Eg in a two-sector
// model, the Industries Table contains two elements, DI and DII
type Simulation struct {
	Manager Manager  // Manager for the Stages of this simulation
	Stages  []*Stage // All Stages generated during one simulation
}

// Constructor for a Simulation with empty Stages and nil Manager elements
//
//	returns:	a pointer to a new empty Simulation
func NewSimulation() *Simulation {
	stages := make([]*Stage, 0)
	simulation := Simulation{Stages: stages}
	return &simulation
}

// A User record contains everything relevant to the simulations of a single logged in user
type User struct {
	UserName            string              `json:"username"` // Repeats the key in the map,for ease of use
	Email               string              `json:"email"`
	ApiKey              string              `json:"api_key"` // The api key allocated to this user
	Password            string              `json:"password"`
	Role                string              `json:"role"`
	CurrentSimulationID int                 `json:"current_simulation_id"` // the id of the simulation that this user is currently using
	CurrentPage         CurrentPageType     // more information about what the user was looking at (under development)
	TimeStamp           int                 // Indexes Datasets. Selects the stage that the simulation has reached
	ViewedTimeStamp     int                 // Indexes Datasets. Selects what the user is viewing
	ComparatorTimeStamp int                 // Indexes Datasets. Selects what Viewed items are compared with.
	Simulations         map[int]*Simulation // Simulations, indexed by SimulationId
}

// Constructor for a standard initial User.
func NewUser(username string) *User {
	newUser := User{
		UserName:            username,
		Password:            "",
		ApiKey:              "",
		CurrentSimulationID: 0,
		CurrentPage:         CurrentPageType{"", 0},
		TimeStamp:           0,
		ViewedTimeStamp:     0,
		ComparatorTimeStamp: 0,
		Simulations:         make(map[int]*Simulation, 0), // TODO under development
	}
	return &newUser
}

// List of LoggedInUsers
var LoggedInUsers = make(map[string]*User) // Every user's simulation data

// A RegisteredUser is used for local authentication
// A User is a logged-in RegisteredUser
type RegisteredUser struct {
	UserName string
	ApiKey   string `json:"api_key"` // The api key will be retrieved from the server
	Password string // hashed
	Cookie   string // TODO NOT USED DEPRECATE
}

// A RegisteredUserServerRequest is used to send a RegisteredUser to the server
type RegisteredUserServerRequest struct {
	UserName string `json:"username"` // Only this field is sent to the server, for security reasons
}

func (u RegisteredUser) Write() string {
	result, _ := json.MarshalIndent(u, " ", " ")
	return string(result)
}

func NewRegisteredUser(username string, password string, apikey string) *RegisteredUser {
	new_RegisteredUser := RegisteredUser{
		UserName: username,
		Password: password,
		ApiKey:   apikey,
		Cookie:   "",
	}
	return &new_RegisteredUser
}

// TODO restore or deprecate
// Find the simulation with a given id.
//
//	u: the user to whom the simulation belongs
//	Return: pointer to the simulation if it found
//	Return: nil if not found.
// func (u *User) Simulation(id int) *Manager {
// 	// fmt.Println("Looking for a simulation with id", id, "Length of list is", len(*u.SimulationsList()))
// 	simulationList := u.SimulationsList()
// 	for i := 0; i < len(*simulationList); i++ {
// 		s := (*simulationList)[i]
// 		// fmt.Printf("Looking for simulationId %d and checking the simulation whose Id is %d\n", id, s.Id)
// 		if id == s.Id {
// 			return &s
// 		}
// 	}
// 	return nil
// }

// Return a pointer to the TimeStamp of the user's current simulation
// Temporary stepping stone
func (u *User) GetTimeStamp() *int {
	// s := u.Simulation(u.CurrentSimulationID)
	return &u.TimeStamp
}

// Return a pointer to the viewed TimeStamp of the user's current simulation
// Temporary stepping stone
func (u *User) GetViewedTimeStamp() *int {
	// s := u.Simulation(u.CurrentSimulationID)
	return &u.ViewedTimeStamp
}

// Return a pointer to the comparator TimeStamp of the user's current simulation
// Temporary stepping stone
func (u *User) GetComparatorTimeStamp() *int {
	// s := u.Simulation(u.CurrentSimulationID)
	return &u.ComparatorTimeStamp
}

package models

//TODO replace object finders with generics and maybe they don't belong in this file anyhow

import (
	"encoding/json"
)

// A User record contains everything relevant to the simulations of a single logged in user
type User struct {
	UserName            string              `json:"username"` // Repeats the key in the map,for ease of use
	Email               string              `json:"email"`
	ApiKey              string              `json:"api_key"` // The api key allocated to this user
	Password            string              `json:"password"`
	Role                string              `json:"role"`
	CurrentSimulationID int                 `json:"current_simulation_id"` // the id of the simulation that this user is currently using
	CurrentPage         CurrentPageType     // more information about what the user was looking at
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
		Simulations:         make(map[int]*Simulation, 0),
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

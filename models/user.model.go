package models

//TODO replace object finders with generics and maybe they don't belong in this file anyhow

import (
	"encoding/json"
)

// A record describing what page the user was visiting
// together with the information needed to display the page
type CurrentPager struct {
	Url string
	Id  int
}

// A User record contains everything relevant to the simulations of a single logged in user
type User struct {
	UserName            string       `json:"username"` // Repeats the key in the map,for ease of use
	Email               string       `json:"email"`
	ApiKey              string       `json:"api_key"` // The api key allocated to this user
	Password            string       `json:"password"`
	Role                string       `json:"role"`
	CurrentSimulationID int          `json:"current_simulation_id"` // the id of the simulation that this user is currently using
	CurrentPage         CurrentPager // more information about what the user was looking at (under development)
	TimeStamp           int          // Indexes Datasets. Selects the stage that the simulation has reached
	ViewedTimeStamp     int          // Indexes Datasets. Selects what the user is viewing
	ComparatorTimeStamp int          // Indexes Datasets. Selects what Viewed items are compared with.
	Simulations         Table        // Details of all simulations
	TableSets           []*TableSet  // Repository for the data objects generated during the simulation
}

// Constructor for a standard initial User.
func NewUser(username string) *User {
	newUser := User{
		UserName:            username,
		Password:            "",
		ApiKey:              "",
		CurrentSimulationID: 0,
		CurrentPage:         CurrentPager{"", 0},
		TimeStamp:           0,
		ViewedTimeStamp:     0,
		ComparatorTimeStamp: 0,
		TableSets:           []*TableSet{},
		Simulations: Table{
			ApiUrl: `/simulations`,
			Table:  new([]Simulation),
			Name:   "Simulations",
		},
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

func (u *User) Write() string {
	if u == nil {
		return "no such user"
	}
	s, _ := json.MarshalIndent(*u, " ", " ")
	return string(s)
}

// Find the commodity with a given id.
//
//	u: the user to whom the commodity belongs
//	Return: pointer to the commodity if it found
//	Return: pointer to NotFoundCommodity if not found.
func (u User) Commodity(id int) *Commodity {
	// commodityList := *LoggedInUsers[u.UserName].OldCommodities()
	commodityList := *ViewedObjects[Commodity](u, `commodities`)
	for i := 0; i < len(commodityList); i++ {
		c := commodityList[i]
		if id == c.Id {
			return &c
		}
	}
	return &NotFoundCommodity
}

// Find the simulation with a given id.
//
//	u: the user to whom the simulation belongs
//	Return: pointer to the simulation if it found
//	Return: nil if not found.
func (u *User) Simulation(id int) *Simulation {
	simulationList := LoggedInUsers[u.UserName].Simulations.Table.(*[]Simulation)
	for i := 0; i < len(*simulationList); i++ {
		s := (*simulationList)[i]
		if id == s.Id {
			return &s
		}
	}
	return nil
}

// Find the class with a given id.
//
//	u: the user to whom the class belongs
//	Return: pointer to the class if it found
//	Return: pointer to NotFoundClass if not found.
func (u User) Class(id int) *Class {
	// classList := *LoggedInUsers[u.UserName].OldClasses()
	classList := *ViewedObjects[Class](u, `classes`)
	for i := 0; i < len(classList); i++ {
		c := classList[i]
		if id == c.Id {
			return &c
		}
	}
	return &NotFoundClass
}

// Find the industry with a given id.
//
//	u: the user to whom the industry belongs
//	Return: pointer to the industry if it found
//	Return: pointer to NotFoundIndustry if not found.
func (u User) Industry(id int) *Industry {
	// industryList := *LoggedInUsers[u.UserName].OldIndustries()
	industryList := *ViewedObjects[Industry](u, `industries`)
	for i := 0; i < len(industryList); i++ {
		ind := industryList[i]
		if id == ind.Id {
			return &ind
		}
	}
	return &NotFoundIndustry
}

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

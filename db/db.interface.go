package db

import (
	"gorilla-client/models"
)

// Interface for database solutions for authorization purposes.
// We don't need extensive query facilities.
// We just need to create, delete and find users.
type DataHandler interface {
	FindRegisteredUser(Name string) (*models.RegisteredUser, error)
	CreateRegisteredUser(u *models.RegisteredUser) (err error)
	UpdateRegisteredUser(u *models.RegisteredUser) (*models.RegisteredUser, error)
	List() string
}

// global variable for the database created when the server starts
var DataBase DataHandler

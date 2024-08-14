// PATH: go-auth/models/inmemory.go
// legacy in-memory databas for authorization

package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorilla-client/models"
	"gorilla-client/utils"
)

// Barebones in memory database
// Implements DataHander interface
//  NewDB creates an instance of an in-memory store for users
//  CreateUser adds a user to the store
//  FindUser finds a user in the store

// An imdbStruct is a single database.
// it should be created using NewDB()
type imdbStruct struct {
	store map[string]*models.RegisteredUser
}

// Creates a new in-memory store
func NewImDB() imdbStruct {
	var imdb = imdbStruct{}
	imdb.store = make(map[string]*models.RegisteredUser)
	return imdb
}

// Implements DataHandler Create(*RegisteredUser)
//
//	u: the address of a RegisteredUser
//	returns: error (nil if create succeeds)
func (s imdbStruct) CreateRegisteredUser(u *models.RegisteredUser) (err error) {
	// Check for exists already
	if _, ok := s.store[u.UserName]; ok {
		return utils.TraceError(fmt.Sprintf("user %s already exists", u.UserName))
	}
	s.store[u.UserName] = u
	utils.TraceInfo(utils.BrightMagenta, fmt.Sprintf("User %s has been added to the local Database", u.UserName))
	return nil
}

// Implements DataHandler Find(*User)
//
//	name: the name of the user
func (s imdbStruct) FindRegisteredUser(name string) (*models.RegisteredUser, error) {
	result, ok := s.store[name]
	if !ok {
		return nil, errors.New("user does not exist")
	}
	return result, nil
}

// diagnostic functiom to dump the whole store
//
//	returns: formatted string containing the contents of the store
func (s imdbStruct) List() string {
	result, _ := json.MarshalIndent(s.store, " ", " ")
	return string(result)
}

func (s imdbStruct) UpdateRegisteredUser(u *models.RegisteredUser) (*models.RegisteredUser, error) {
	return nil, errors.New("working on it")
}

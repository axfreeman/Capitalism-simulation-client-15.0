package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorilla-client/models"
	"gorilla-client/utils"
)

// Database using the remote server API
// Implements DataHander interface
//
//  NewDB creates an instance of the db.
//  CreateUser adds a user to the store
//  FindUser finds a user in the store

// A RemoteDBStruct defines a single database.
// it should be created using NewDB()
type RemoteDbStruct struct {
	store map[string]*models.User
}

// Creates a new in-memory store
func NewRemoteDB() RemoteDbStruct {
	var rdb = RemoteDbStruct{}
	rdb.store = make(map[string]*models.User)
	return rdb
}

// Implements DataHandler Create(*User)
//
//	u: the address of a User
func (s RemoteDbStruct) CreateUser(u *models.User) (err error) {
	name := u.UserName
	utils.TraceInfo(utils.BrightMagenta, fmt.Sprintf("User %s has been added to the local Database", name))
	// Check for exists already
	if _, ok := s.store[name]; ok {
		return utils.TraceError(fmt.Sprintf("user %s already exists", u.UserName))
	}
	s.store[name] = u
	return nil
}

// Implements DataHandler Find(*User)
//
//	name: the name of the user
func (s RemoteDbStruct) FindUser(name string) (*models.User, error) {
	result, ok := s.store[name]
	if !ok {
		return nil, errors.New("user does not exist")
	}
	return result, nil
}

// diagnostic functiom to dump the whole store
//
//	returns: formatted string containing the contents of the store
func (s RemoteDbStruct) List() string {
	result, _ := json.MarshalIndent(s.store, " ", " ")
	return string(result)
}

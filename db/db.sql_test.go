package db

import (
	"gorilla-client/models"
	"gorilla-client/utils"
	"testing"
)

func TestSQLDB(t *testing.T) {
	var err error
	var TestUser *models.RegisteredUser
	utils.LogInit()
	db := NewSQLDB()
	db.CreateRegisteredUser(models.NewRegisteredUser("TestUser", "", ""))
	if TestUser, err = db.FindRegisteredUser("TestUser"); err != nil {
		t.Errorf("Database failed to find test user because: %s", err)
	}
	db.List()
	if TestUser.UserName != `TestUser` {
		t.Errorf("Database found the wrong user")
	}
	if _, err = db.FindRegisteredUser("NonExistentUser"); err == nil {
		t.Errorf("Database failed to report non existent user because: %s", err)
	}
}

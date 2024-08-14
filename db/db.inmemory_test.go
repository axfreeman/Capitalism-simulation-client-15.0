package db

import (
	"gorilla-client/models"
	"gorilla-client/utils"
	"testing"
)

func TestDB(t *testing.T) {
	var err error
	var TestUser *models.RegisteredUser
	utils.LogInit()
	db := NewImDB()
	db.CreateRegisteredUser(models.NewRegisteredUser("TestUser", "", ""))
	if TestUser, err = db.FindRegisteredUser("TestUser"); err != nil {
		t.Errorf("Database failed to find test user because: %s", err)
	}
	if TestUser.UserName != `TestUser` {
		t.Errorf("Database found the wrong user")
	}
	if _, err = db.FindRegisteredUser("NonExistentUser"); err == nil {
		t.Errorf("Database failed to report non existent user because: %s", err)
	}
}

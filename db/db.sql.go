package db

import (
	"database/sql"
	"errors"
	"gorilla-client/config"
	"gorilla-client/models"
	"gorilla-client/utils"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Database using SQlite
// Implements DataHander interface
//
//  NewDB creates an instance of the db.
//  CreateUser adds a user to the store
//  FindUser finds a user in the store

// A RemoteDBStruct defines a single database.
// it should be created using NewDB()
type SQLDbStruct struct {
	db *sql.DB
}

// Receiver for a query
// TODO should simply use the RegisteredUser struct
type SQLdbEntry struct {
	username string
	password string
	apikey   string
}

// Creates a new SQLite store
func NewSQLDB() SQLDbStruct {
	// This is a short-term database. We recreat it every time this app
	// starts, rather than investing in time-consuming migration software.
	// The api server is our permanent repository.
	os.Remove(config.Config.SQLiteFile)

	file, err := os.Create(config.Config.SQLiteFile) // Create new SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	utils.TraceInfo(utils.BrightWhite, "New sqlite file created")

	// Now open the database we just created
	sdb, err := sql.Open("sqlite3", config.Config.SQLiteFile)
	if err != nil {
		log.Fatalf("Could not open SQLite file because:%v. Cannot continue", err)
	}

	// defer sdb.Close() // Defer Closing the database NOTE this stops us inserting anything

	// Create the user table
	_, err = sdb.Exec("CREATE TABLE IF NOT EXISTS `users` (`username` VARCHAR(64) PRIMARY KEY, `password` VARCHAR(256) NOT NULL,`apikey` VARCHAR(256) NOT NULL);")
	if err != nil {
		log.Fatalf("Could not create an SQlite table because:%v. Cannot continue", err)
	}

	utils.TraceInfo(utils.BrightMagenta, "Local Database Created")
	return SQLDbStruct{sdb}
}

// Implements DataHandler Create(*User)
//
//	u: the address of a RegisteredUser
func (s SQLDbStruct) CreateRegisteredUser(u *models.RegisteredUser) error {
	var err error
	insertSQL := "INSERT INTO users (username,password,apikey) VALUES(?,?,?)"
	statement, err := s.db.Prepare(insertSQL)
	if err != nil {
		utils.TraceErrorf("Failed to add user because %v: ", err)
		return err
	}
	_, err = statement.Exec(u.UserName, u.Password, u.ApiKey)

	if err != nil {
		utils.TraceErrorf("The treatment worked but the statement died %v", err)
	}
	utils.TraceInfof(utils.BrightMagenta, "User %s has been added to the local Database", u.UserName)
	return nil
}

// Implements DataHandler Find(*User)
//
//	name: the name of the user
func (s SQLDbStruct) FindRegisteredUser(name string) (*models.RegisteredUser, error) {
	var entry SQLdbEntry
	var err error
	row := s.db.QueryRow("SELECT * FROM users WHERE username=?", name)
	if err = row.Scan(&entry.username, &entry.password, &entry.apikey); err == sql.ErrNoRows {
		return nil, errors.New("user does not exist")
	}

	utils.TraceInfof(utils.BrightMagenta, "Found user %s", entry.username)
	return models.NewRegisteredUser(entry.username, entry.password, entry.apikey), nil
}

// diagnostic functiom to dump the whole store
//
//	returns: formatted string containing the contents of the store
func (s SQLDbStruct) List() string {
	row, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var username string
		var password string
		row.Scan(&username, &password)
		log.Println("User: ", username, " ", password)
	}
	return "Completed Listing of the database"
}

func (s SQLDbStruct) UpdateRegisteredUser(u *models.RegisteredUser) (*models.RegisteredUser, error) {
	var entry SQLdbEntry
	var err error
	row := s.db.QueryRow("UPDATE users SET password=? WHERE username= ?", u.Password, u.UserName)
	if err = row.Scan(&entry.username, &entry.password, &entry.apikey); err == sql.ErrNoRows {
		return nil, errors.New("user does not exist")
	}
	utils.TraceInfof(utils.BrightMagenta, "Updated user %s", entry.username)
	return models.NewRegisteredUser(entry.username, entry.password, entry.password), nil
}

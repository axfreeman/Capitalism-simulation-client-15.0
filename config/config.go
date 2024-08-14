package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Cfg struct {
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	SSLMode    string
	ApiSource  string
	AdminUser  string
	AdminKey   string
	ClientHost string
	LogFile    string
	SQLiteFile string
}

var Config Cfg

func Init() (err error) {
	// Load .env file and Create a new connection to the database
	// NOTE not used but preserved for reference
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	Config = Cfg{
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		SSLMode:    os.Getenv("DB_SSLMODE"),
		ApiSource:  os.Getenv("APISOURCE"),
		AdminUser:  os.Getenv("ADMINUSER"),
		AdminKey:   os.Getenv("ADMINKEY"),
		ClientHost: os.Getenv("CLIENT_HOST"),
		LogFile:    os.Getenv("LOG_FILE"),
		SQLiteFile: os.Getenv("SQLITE_FILE"),
	}
	return err
}

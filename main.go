package main

import (
	"gorilla-client/api"
	"gorilla-client/config"
	"gorilla-client/controllers"
	"gorilla-client/db"
	"gorilla-client/models"
	"gorilla-client/routes"
	"gorilla-client/utils"
	"html/template"
	"log"
	"net/http"
)

func main() {
	utils.LogInit()

	config.Init()

	models.InitViews()

	utils.TraceInfo(utils.Yellow, "The Rosy Dawn of Capitalism has begun")

	db.DataBase = db.NewSQLDB()

	api.LoadRegisteredUsers()

	controllers.Tpl, _ = template.ParseGlob("./templates/*/*")

	routes.AuthRoutes()

	err := http.ListenAndServe("localhost:8080", routes.Router)
	if err != nil {
		log.Fatal(err)
	}
}

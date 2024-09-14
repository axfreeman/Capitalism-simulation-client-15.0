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

	var err error

	// controllers.Tpl, err = template.ParseGlob("./templates/*/*")

	funcMap := template.FuncMap{
		"Show":          models.Show,
		"Link":          models.Link,
		"OriginGraphic": models.OriginGraphic,
		"UsageGraphic":  models.UsageGraphic,
		"CommodityLink": models.CommodityLink,
	}

	controllers.Tpl, err = template.New("").Funcs(funcMap).ParseGlob("./templates/*/*")
	if err != nil {
		log.Fatal(err)
	}

	routes.AuthRoutes()

	err = http.ListenAndServe("localhost:8080", routes.Router)
	if err != nil {
		log.Fatal(err)
	}
}

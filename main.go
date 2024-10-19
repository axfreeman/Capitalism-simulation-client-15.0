package main

import (
	"html/template"
	"log"
	"net/http"
	"simulation-client/api"
	"simulation-client/config"
	"simulation-client/controllers"
	"simulation-client/db"
	"simulation-client/models"
	"simulation-client/routes"
	"simulation-client/utils"
	"simulation-client/views"
)

func main() {

	utils.LogInit()

	config.Init()

	utils.TraceInfo(utils.Yellow, "The Rosy Dawn of Capitalism has begun")

	db.DataBase = db.NewImDB()

	api.LoadRegisteredUsers()

	var err error

	// controllers.Tpl, err = template.ParseGlob("./templates/*/*")

	funcMap := template.FuncMap{
		"Show":                       views.Show,
		"ShowString":                 views.ShowString,
		"ShowDecimal":                views.ShowDecimal,
		"Link":                       views.Link,
		"OriginGraphic":              models.OriginGraphic,
		"UsageGraphic":               models.UsageGraphic,
		"IndustryCommodityLink":      models.IndustryCommodityLink,
		"ClassCommodityLink":         models.ClassCommodityLink,
		"StockIndustryLink":          models.StockIndustryLink,
		"IndustryStockCommodityLink": models.IndustryStockCommodityLink,
		"ClassStockCommodityLink":    models.ClassStockCommodityLink,
		"StockClassLink":             models.StockClassLink,
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

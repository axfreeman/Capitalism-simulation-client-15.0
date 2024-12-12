package main

import (
	"html/template"
	"log"
	"net/http"
	"simulation-client/api"
	"simulation-client/config"
	"simulation-client/db"
	"simulation-client/logging"
	"simulation-client/routes"
	"simulation-client/views"
)

func main() {

	logging.LogInit()

	config.Init()

	logging.TraceInfo(logging.Yellow, "The Rosy Dawn of Capitalism has begun")

	db.DataBase = db.NewImDB()

	api.LoadRegisteredUsers()

	var err error

	// logging.Tpl, err = template.ParseGlob("./templates/*/*")

	funcMap := template.FuncMap{
		"Show":                       views.Show,
		"ShowString":                 views.ShowString,
		"ShowPlainString":            views.ShowPlainString,
		"ShowDecimal":                views.ShowDecimal,
		"Link":                       views.Link,
		"OriginGraphic":              views.OriginGraphic,
		"UsageGraphic":               views.UsageGraphic,
		"IndustryCommodityLink":      views.IndustryCommodityLink,
		"ClassCommodityLink":         views.ClassCommodityLink,
		"StockIndustryLink":          views.StockIndustryLink,
		"IndustryStockCommodityLink": views.IndustryStockCommodityLink,
		"ClassStockCommodityLink":    views.ClassStockCommodityLink,
		"StockClassLink":             views.StockClassLink,
	}

	views.Tpl, err = template.New("").Funcs(funcMap).ParseGlob("./templates/*/*")
	if err != nil {
		log.Fatal(err)
	}

	routes.AuthRoutes()

	err = http.ListenAndServe("localhost:8080", routes.Router)
	if err != nil {
		log.Fatal(err)
	}
}

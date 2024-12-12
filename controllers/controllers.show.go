// display.objects.go
// handlers to display the objects of the simulation on the user's browser

package controllers

import (
	"fmt"
	"net/http"
	"simulation-client/api"
	"simulation-client/config"
	"simulation-client/logging"
	"simulation-client/models"
	"simulation-client/views"
)

// TODO temporary fix to reset the API database when we are offline
func Reset(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered Reset")
	var err error

	// user := *models.LoggedInUsers["Admin"]

	user := CurrentUser(r)

	status, err := api.AdminGetRequest(config.Config.ApiSource+"/action/reset/", &user)
	logging.TraceInfo(logging.BrightGreen, fmt.Sprintf("The server responded with status %d and error %v", status, err))
	if status != http.StatusOK {
		logging.TraceError("The server doesn't know this user, sorry")
		views.Tpl.ExecuteTemplate(w, "login.html", "Check username and password")
		return
	}
}

// display all commodities in the current simulation
func ShowCommodities(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "commodities.html", Id: 0}

	logging.TraceInfof(logging.BrightYellow, "Fetching commodities for user %s", user.UserName)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, ""))
}

// display all industries in the current simulation
func ShowIndustries(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "industries.html", Id: 0}

	logging.TraceInfof(logging.BrightYellow, "Fetching industries for user %s", user.UserName)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, ""))
}

// display all classes in the current simulation
func ShowClasses(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "classes.html", Id: 0}

	logging.TraceInfo(logging.BrightYellow, fmt.Sprintf("Fetching classes for user %s", user.UserName))
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, ""))
}

// display all industry stocks in the current simulation
func ShowIndustryStocks(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "industry_stocks.html", Id: 0}

	logging.TraceInfof(logging.BrightYellow, "Fetching industry stocks for user %s", user.UserName)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, ""))
}

// display all the class stocks in the current simulation
func ShowClassStocks(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "class_stocks.html", Id: 0}

	logging.TraceInfof(logging.BrightYellow, "Fetching class stocks for user %s", user.UserName)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, ""))
}

// display all Trace records in the current simulation
func ShowTrace(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "trace.html", Id: 0}
	logging.TraceInfof(logging.BrightYellow, "Fetching trace for user %s", user.UserName)
	t := models.Traces(user)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.TemplateData{Trace: t})
}

// Display one specific commodity
func ShowCommodity(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	user := CurrentUser(r)
	if id, err = FetchIDfromURL(r); err != nil {
		ReportError(user, w, err.Error())
	}
	user.CurrentPage = models.CurrentPageType{Url: "commodity.html", Id: id}

	logging.TraceInfof(logging.BrightYellow, "Fetching commodity %d for user %s", id, user.UserName)
	views.Tpl.ExecuteTemplate(w,
		user.CurrentPage.Url,
		views.CommodityDisplayData(user, "", id))
}

// Display one specific industry
func ShowIndustry(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	user := CurrentUser(r)
	if id, err = FetchIDfromURL(r); err != nil {
		ReportError(user, w, err.Error())
	}
	user.CurrentPage = models.CurrentPageType{Url: "industry.html", Id: id}

	logging.TraceInfof(logging.BrightYellow, "Fetching industry %d for user %s", id, user.UserName)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.IndustryDisplayData(user, "", id))
}

// Display one specific class
func ShowClass(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	user := CurrentUser(r)
	if id, err = FetchIDfromURL(r); err != nil {
		ReportError(user, w, err.Error())
	}
	user.CurrentPage = models.CurrentPageType{Url: "class.html", Id: id}

	logging.TraceInfof(logging.BrightYellow, "Fetching class %d for user %s", id, user.UserName)
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.ClassDisplayData(user, "", id))
}

// Display one specific industry stock
func ShowIndustryStock(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	user := CurrentUser(r)
	if id, err = FetchIDfromURL(r); err != nil {
		ReportError(user, w, err.Error())
	}
	user.CurrentPage = models.CurrentPageType{Url: "industry_stock.html", Id: id}

	logging.TraceInfof(logging.BrightYellow, "Fetching industry_stock %d for user %s", id, user.UserName)
	views.Tpl.ExecuteTemplate(w,
		user.CurrentPage.Url,
		views.IndustryStockDisplayData(user, "", id))
}

// Displays a snapshot of the economy
func ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "index.html", Id: 0}

	logging.TraceInfo(logging.BrightYellow, fmt.Sprintf("Showing Index Page for user %s", user.UserName))
	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, ""))
}

func UserDashboard(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "user-dashboard.html", Id: 0}

	views.Tpl.ExecuteTemplate(w, user.CurrentPage.Url, views.CreateTemplateData(user, ""))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	views.Tpl.ExecuteTemplate(w, "404.html", "")
}

// check session for logged in done with middleware Auth()
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	logging.TraceInfo(logging.BrightGreen, "Enter WelcomeHandler")
	user := CurrentUser(r)
	views.Tpl.ExecuteTemplate(w, "welcome.html", views.CreateTemplateData(user, ""))
}

// TODO remove. Just a basic test page
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter AboutHandler")
	views.Tpl.ExecuteTemplate(w, "about.html", "test")
}

// Diagnostic function mainly for the developer, to show all the DisplayData
// TODO this is a crude implementation. There is probably a better way
func AllDisplayData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter ShowDisplayData")
	user := CurrentUser(r)
	// views.Tpl.ExecuteTemplate(w, "displayData.html", views.CreateTemplateData(user,""))

	templateData := views.CreateTemplateData(user, "")
	logging.TraceLogf(logging.White, "Template data is %v\n", templateData)

	// Log all the commodities
	commodityData := templateData.CommodityViews
	fmt.Println("Commodities")
	logging.TraceLogf(logging.White, "CommodityViews (%v)\n", commodityData)
	for i := range *commodityData {
		v := (*commodityData)[i].(*views.CommodityView).Viewed().(*models.Commodity)
		fmt.Println(v.Write())
	}
	// Log all the Industries
	industryData := templateData.IndustryViews
	fmt.Println("Industries")
	logging.TraceLogf(logging.White, "IndustryViews (%v)\n", industryData)
	for i := range *industryData {
		v := (*industryData)[i].(*views.IndustryView).Viewed().(*models.Industry)
		fmt.Println(v.Write())
	}

	// TODO send this to a page to be viewed
}

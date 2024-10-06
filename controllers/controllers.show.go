// display.objects.go
// handlers to display the objects of the simulation on the user's browser

package controllers

import (
	"fmt"
	"gorilla-client/models"
	"gorilla-client/utils"
	"net/http"
)

// display all commodities in the current simulation
func ShowCommodities(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "commodities.html", Id: 0}

	utils.TraceInfof(utils.BrightYellow, "Fetching commodities for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all industries in the current simulation
func ShowIndustries(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "industries.html", Id: 0}

	utils.TraceInfof(utils.BrightYellow, "Fetching industries for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all classes in the current simulation
func ShowClasses(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "classes.html", Id: 0}

	utils.TraceInfo(utils.BrightYellow, fmt.Sprintf("Fetching classes for user %s", user.UserName))
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all industry stocks in the current simulation
func ShowIndustryStocks(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "industry_stocks.html", Id: 0}

	utils.TraceInfof(utils.BrightYellow, "Fetching industry stocks for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all the class stocks in the current simulation
func ShowClassStocks(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "class_stocks.html", Id: 0}

	utils.TraceInfof(utils.BrightYellow, "Fetching class stocks for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all Trace records in the current simulation
func ShowTrace(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	// user.CurrentPage = models.CurrentPageType{Url: "trace.html", Id: 0}
	user.CurrentPage = models.CurrentPageType{Url: "elodieb.html", Id: 0}
	utils.TraceInfof(utils.BrightYellow, "Fetching classes for user %s", user.UserName)
	t := models.Traces(user)
	utils.UNUSED(t)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, models.TemplateData{Trace: t})
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

	utils.TraceInfof(utils.BrightYellow, "Fetching commodity %d for user %s", id, user.UserName)
	Tpl.ExecuteTemplate(w,
		user.CurrentPage.Url,
		models.CommodityDisplayData(user, "", id))
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

	utils.TraceInfof(utils.BrightYellow, "Fetching industry %d for user %s", id, user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.IndustryDisplayData("", id))
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

	utils.TraceInfof(utils.BrightYellow, "Fetching class %d for user %s", id, user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.ClassDisplayData("", id))
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

	utils.TraceInfof(utils.BrightYellow, "Fetching industry_stock %d for user %s", id, user.UserName)
	Tpl.ExecuteTemplate(w,
		user.CurrentPage.Url,
		user.IndustryStockDisplayData("", id))
}

// Displays a snapshot of the economy
func ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "index.html", Id: 0}

	utils.TraceInfo(utils.BrightYellow, fmt.Sprintf("Showing Index Page for user %s", user.UserName))
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

func UserDashboard(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "user-dashboard.html", Id: 0}

	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Tpl.ExecuteTemplate(w, "404.html", "")
}

// check session for logged in done with middleware Auth()
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Enter WelcomeHandler")
	user := CurrentUser(r)
	Tpl.ExecuteTemplate(w, "welcome.html", user.CreateTemplateData(""))
}

// TODO remove. Just a basic test page
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter AboutHandler")
	Tpl.ExecuteTemplate(w, "about.html", "Logged In")
}

// Diagnostic function mainly for the developer, to show all the DisplayData
// TODO this is a crude implementation. There is probably a better way
func AllDisplayData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter ShowDisplayData")
	user := CurrentUser(r)
	// Tpl.ExecuteTemplate(w, "displayData.html", user.CreateTemplateData(""))

	templateData := user.CreateTemplateData("")
	utils.TraceLogf(utils.White, "Template data is %v\n", templateData)

	// Log all the commodities
	commodityData := templateData.CommodityViews
	fmt.Println("Commodities")
	utils.TraceLogf(utils.White, "CommodityViews (%v)\n", commodityData)
	for i := range *commodityData {
		v := (*commodityData)[i].(*models.CommodityView).Viewed().(*models.Commodity)
		fmt.Println(v.Write())
	}
	// Log all the Industries
	industryData := templateData.IndustryViews
	fmt.Println("Industries")
	utils.TraceLogf(utils.White, "IndustryViews (%v)\n", industryData)
	for i := range *industryData {
		v := (*industryData)[i].(*models.IndustryView).Viewed().(*models.Industry)
		fmt.Println(v.Write())
	}

	// TODO send this to a page to be viewed
}

// under development
func DisplayStage(*models.Stage) string {
	return "wtf"
}

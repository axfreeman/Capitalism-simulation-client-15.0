// PATH: go-auth/routes/auth.go

package routes

import (
	"net/http"
	"simulation-client/auth"
	"simulation-client/controllers"

	"github.com/gorilla/mux"
)

var Router *mux.Router

// TODO convert " to ` where appropriate

func AuthRoutes() {
	// Export router to globally accessible variable
	var ds http.Handler
	Router = mux.NewRouter()

	// create static file server for css, js, and image files
	fs := http.FileServer(http.Dir("./static"))
	ds = http.StripPrefix("/static/", fs)
	Router.PathPrefix("/static/").Handler(ds)

	Router.HandleFunc("/auth/login", auth.LoginFormDisplay)
	Router.HandleFunc("/auth/loginauth", auth.LoginAuthHandler)
	Router.HandleFunc("/auth/logout", auth.LogoutHandler)
	Router.HandleFunc("/auth/register", auth.RegisterHandler)
	Router.HandleFunc("/auth/registerauth", auth.RegisterAuthHandler)

	Router.HandleFunc("/about", auth.Auth(controllers.AboutHandler))
	Router.HandleFunc("/welcome", auth.Auth(controllers.WelcomeHandler))
	Router.HandleFunc("/user/data", controllers.AllData)
	Router.HandleFunc("/user/table-data", controllers.DisplayData)
	Router.HandleFunc("/user/dashboard", auth.Auth(controllers.UserDashboard))
	Router.HandleFunc(`/user/delete/{id}`, auth.Auth(controllers.DeleteSimulation))
	Router.HandleFunc(`/user/switch/{id}`, auth.Auth(controllers.SwitchSimulation))
	Router.HandleFunc(`/user/restart/{id}`, auth.Auth(controllers.RestartSimulation))

	// actions
	Router.HandleFunc("/action/{action}", controllers.ActionHandler)

	// Display controls
	Router.HandleFunc("/user/forward", controllers.Forward)
	Router.HandleFunc("/user/back", controllers.Back)
	Router.HandleFunc("/user/create/{id}", controllers.CreateSimulation)
	Router.HandleFunc("/user/display-size", controllers.DisplaySize)
	Router.HandleFunc("/user/display-value", controllers.DisplayValue)
	Router.HandleFunc("/user/display-price", controllers.DisplayPrice)

	// Table displays
	Router.HandleFunc("/commodities", auth.Auth(controllers.ShowCommodities))
	Router.HandleFunc("/industries", auth.Auth(controllers.ShowIndustries))
	Router.HandleFunc("/classes", auth.Auth(controllers.ShowClasses))
	Router.HandleFunc("/industry_stocks", auth.Auth(controllers.ShowIndustryStocks))
	Router.HandleFunc("/industry_stock/{id}", auth.Auth(controllers.ShowIndustryStock))
	Router.HandleFunc("/class_stocks", auth.Auth(controllers.ShowClassStocks))
	Router.HandleFunc("/class_stock/{id}", auth.Auth(controllers.ShowClassStock))
	Router.HandleFunc("/industry_stock/{id}", auth.Auth(controllers.ShowIndustryStock))
	Router.HandleFunc("/commodity/{id}", auth.Auth(controllers.ShowCommodity))
	Router.HandleFunc("/industry/{id}", auth.Auth(controllers.ShowIndustry))
	Router.HandleFunc("/class/{id}", auth.Auth(controllers.ShowClass))
	Router.HandleFunc("/trace", auth.Auth(controllers.ShowTrace))
	Router.HandleFunc("/index", auth.Auth(controllers.ShowIndexPage))
	Router.HandleFunc("/", auth.Auth(controllers.ShowIndexPage))
	Router.HandleFunc(`/download`, auth.Auth(controllers.Download))
	Router.HandleFunc(`/all-display-data`, auth.Auth(controllers.AllDisplayData))

	// TODO purely temporary fix to send reset instruction to API
	// for developmental use when we have no internet
	Router.HandleFunc(`/reset`, auth.Auth(controllers.Reset))

	// Not found page
	Router.NotFoundHandler = http.HandlerFunc(controllers.NotFound)

}

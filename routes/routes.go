// PATH: go-auth/routes/auth.go

package routes

import (
	"gorilla-client/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

var Router *mux.Router

func AuthRoutes() {
	// Export router to globally accessible variable
	Router = mux.NewRouter()
	Router.HandleFunc("/auth/login", controllers.LoginHandler)
	Router.HandleFunc("/auth/loginauth", controllers.LoginAuthHandler)
	Router.HandleFunc("/auth/logout", controllers.LogoutHandler)
	Router.HandleFunc("/auth/register", controllers.RegisterHandler)
	Router.HandleFunc("/auth/registerauth", controllers.RegisterAuthHandler)

	Router.HandleFunc("/about", controllers.Auth(controllers.AboutHandler))
	Router.HandleFunc("/welcome", controllers.Auth(controllers.WelcomeHandler))
	Router.HandleFunc("/user/data", controllers.AllData)
	Router.HandleFunc("/user/table-data", controllers.DisplayData)
	Router.HandleFunc("/user/dashboard", controllers.Auth(controllers.UserDashboard))
	Router.HandleFunc(`/user/delete/{id}`, controllers.Auth(controllers.DeleteSimulation))
	Router.HandleFunc(`/user/switch/{id}`, controllers.Auth(controllers.SwitchSimulation))
	Router.HandleFunc(`/user/restart/{id}`, controllers.Auth(controllers.RestartSimulation))

	// actions
	Router.HandleFunc("/action/{action}", controllers.ActionHandler)
	Router.HandleFunc("/user/forward", controllers.Forward)
	Router.HandleFunc("/user/back", controllers.Back)
	Router.HandleFunc("/user/create/{id}", controllers.CreateSimulation)

	// Table displays
	Router.HandleFunc("/commodities", controllers.Auth(controllers.ShowCommodities))
	Router.HandleFunc("/industries", controllers.Auth(controllers.ShowIndustries))
	Router.HandleFunc("/classes", controllers.Auth(controllers.ShowClasses))
	Router.HandleFunc("/industry_stocks", controllers.Auth(controllers.ShowIndustryStocks))
	Router.HandleFunc("/class_stocks", controllers.Auth(controllers.ShowClassStocks))
	Router.HandleFunc("/commodity/{id}", controllers.Auth(controllers.ShowCommodity))
	Router.HandleFunc("/industry/{id}", controllers.Auth(controllers.ShowIndustry))
	Router.HandleFunc("/class/{id}", controllers.Auth(controllers.ShowClass))
	Router.HandleFunc("/trace", controllers.Auth(controllers.ShowTrace))
	Router.HandleFunc("/index", controllers.Auth(controllers.ShowIndexPage))
	Router.HandleFunc("/", controllers.Auth(controllers.ShowIndexPage))
	Router.HandleFunc(`/download`, controllers.Auth(controllers.Download))

	Router.NotFoundHandler = http.HandlerFunc(controllers.NotFound)

}

package views

import (
	"html/template"
	"simulation-client/logging"
	"simulation-client/models"
)

func ViewedObjects[T models.Object](u models.User, objectType string) *[]T {
	return (*u.GetViewedStage())[objectType].Table.(*[]T)
}

func ComparedObjects[T models.Object](u models.User, objectType string) *[]T {
	return (*u.GetComparatorStage())[objectType].Table.(*[]T)
}

func ViewedObject[T models.Object](u models.User, objectType string, id int) *T {
	// fmt.Println("ViewedObject was asked to display an object of type ", objectType)
	objectList := (*u.GetViewedStage())[objectType].Table.(*[]T)
	for i := 0; i < len(*objectList); i++ {
		o := (*objectList)[i]
		if id == o.GetId() {
			return &o
		}
	}
	return nil
}

func ComparedObject[T models.Object](u models.User, objectType string, id int) *T {
	objectList := (*u.GetComparatorStage())[objectType].Table.(*[]T)
	for i := 0; i < len(*objectList); i++ {
		o := (*objectList)[i]
		if id == o.GetId() {
			return &o
		}
	}
	return nil
}

// Commonly-used Views to pass into templates
type TemplateData struct {
	Title              string
	Simulations        *[]models.Manager
	Templates          *[]models.Manager
	CommodityViews     *[]Viewer
	IndustryViews      *[]Viewer
	ClassViews         *[]Viewer
	IndustryStockViews *[]Viewer
	ClassStockViews    *[]Viewer
	Trace              *[]template.HTML // exported directly as safe HTML
	Count              int
	Username           string
	State              string
	SetPriceMode       string
	ViewedState        string
	ComparatorState    string
	DisplayDimension   string
	Message            string
}

// Supplies data to pass into Templates for display
//
//		u: a user
//
//		returns:
//	     if the user has no simulations, just the template list
//	     otherwise, the output data the users current simulation
func CreateTemplateData(u *models.User, message string) TemplateData {
	logging.TraceInfof(logging.BrightYellow, "TemplateData is retrieving data for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)
	if u.CurrentSimulationID == 0 {
		logging.TraceInfo(logging.BrightYellow, "User has no simulations")
		return TemplateData{
			Title:              "No simulations",
			Simulations:        nil,
			Templates:          &models.TemplateList,
			Count:              0,
			Username:           u.UserName,
			State:              "UNKNOWN",
			SetPriceMode:       "Auto",
			DisplayDimension:   "Size",
			ViewedState:        "UNKNOWN",
			ComparatorState:    "UNKNOWN",
			CommodityViews:     nil,
			IndustryViews:      nil,
			IndustryStockViews: nil,
			ClassStockViews:    nil,
			Trace:              nil,
			Message:            message,
		}
	}

	// retrieve comparator and viewed records for all data objects
	// to prepare for entry into Views in the DisplayData object

	// uncomment for diagnostics
	// manager := u.GetCurrentSimulation().Manager
	// fmt.Println("***Manager is ", manager.Write())

	cv := ViewedObjects[models.Commodity](*u, `commodities`)
	cc := ComparedObjects[models.Commodity](*u, `commodities`)
	iv := ViewedObjects[models.Industry](*u, `industries`)
	ic := ComparedObjects[models.Industry](*u, `industries`)
	clv := ViewedObjects[models.Class](*u, `classes`)
	clc := ComparedObjects[models.Class](*u, `classes`)
	isv := ViewedObjects[models.IndustryStock](*u, `industry_stocks`)
	isc := ComparedObjects[models.IndustryStock](*u, `industry_stocks`)
	csv := ViewedObjects[models.ClassStock](*u, `class_stocks`)
	csc := ComparedObjects[models.ClassStock](*u, `class_stocks`)

	// Create the DisplayData object
	templateData := TemplateData{
		Title:              "Hello",
		Templates:          &models.TemplateList,
		Username:           u.UserName,
		State:              u.CurrentState(),
		SetPriceMode:       u.SetPriceMode(),
		DisplayDimension:   u.DisplayDimension(),
		ViewedState:        u.ViewedState(),
		ComparatorState:    u.ComparatorState(),
		CommodityViews:     CommodityViews(cv, cc),
		IndustryViews:      IndustryViews(iv, ic),
		ClassViews:         ClassViews(clv, clc),
		IndustryStockViews: IndustryStockViews(isv, isc),
		ClassStockViews:    ClassStockViews(csv, csc),
		Message:            message,
	}

	return templateData
}

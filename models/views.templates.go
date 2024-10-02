package models

import (
	"fmt"
	"gorilla-client/utils"
	"gorilla-client/views"
)

type Object interface {
	Commodity | Industry | Class | IndustryStock | ClassStock | Manager | Trace
	GetId() int
}

func ViewedObjects[T Object](u User, objectType string) *[]T {
	return (*u.GetViewedStage())[objectType].Table.(*[]T)
}

func ComparedObjects[T Object](u User, objectType string) *[]T {
	return (*u.GetComparatorStage())[objectType].Table.(*[]T)
}

func ViewedObject[T Object](u User, objectType string, id int) *T {
	objectList := (*u.GetViewedStage())[objectType].Table.(*[]T)
	for i := 0; i < len(*objectList); i++ {
		o := (*objectList)[i]
		if id == o.GetId() {
			return &o
		}
	}
	return nil
}

func ComparedObject[T Object](u User, objectType string, id int) *T {
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
	Simulations        *[]Manager
	Templates          *[]Manager
	CommodityViews     *[]views.Viewer
	IndustryViews      *[]views.Viewer
	ClassViews         *[]views.Viewer
	IndustryStockViews *[]views.Viewer
	ClassStockViews    *[]views.Viewer
	Trace              *[]Trace
	Count              int
	Username           string
	State              string
	ComparatorState    string
	Message            string
}

// Supplies data to pass into Templates for display
//
//		u: a user
//
//		returns:
//	     if the user has no simulations, just the template list
//	     otherwise, the output data the users current simulation
func (u *User) CreateTemplateData(message string) TemplateData {
	utils.TraceInfof(utils.BrightYellow, "TemplateData is retrieving data for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)
	if u.CurrentSimulationID == 0 {
		utils.TraceInfo(utils.BrightYellow, "User has no simulations")
		return TemplateData{
			Title:              "No simulations",
			Simulations:        nil,
			Templates:          &TemplateList,
			Count:              0,
			Username:           u.UserName,
			State:              "UNKNOWN",
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
	manager := u.GetCurrentSimulation().Manager

	fmt.Println("***Manager is ", manager.Write())

	state := u.GetCurrentState()
	cv := ViewedObjects[Commodity](*u, `commodities`)
	cc := ComparedObjects[Commodity](*u, `commodities`)
	iv := ViewedObjects[Industry](*u, `industries`)
	ic := ComparedObjects[Industry](*u, `industries`)
	clv := ViewedObjects[Class](*u, `classes`)
	clc := ComparedObjects[Class](*u, `classes`)
	isv := ViewedObjects[IndustryStock](*u, `industry stocks`)
	isc := ComparedObjects[IndustryStock](*u, `industry stocks`)
	csv := ViewedObjects[ClassStock](*u, `class stocks`)
	csc := ComparedObjects[ClassStock](*u, `class stocks`)

	// Create the DisplayData object
	templateData := TemplateData{
		Title:              "Hello",
		Templates:          &TemplateList,
		Username:           u.UserName,
		State:              state,
		CommodityViews:     CommodityViews(cv, cc),
		IndustryViews:      IndustryViews(iv, ic),
		ClassViews:         ClassViews(clv, clc),
		IndustryStockViews: IndustryStockViews(isv, isc),
		ClassStockViews:    ClassStockViews(csv, csc),
		Message:            message,
	}

	return templateData
}

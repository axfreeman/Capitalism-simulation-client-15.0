package models

import (
	"fmt"
	"reflect"
)

// TODO fix up IndustryLink and CommodityLink for IndustryStock so it displays the name, not the Id

// Type for implementation of Viewer interface
type IndustryView struct {
	viewedRecord   *Industry
	comparedRecord *Industry
	MoneyView      *IndustryStockView
	VariableView   *IndustryStockView
	ConstantView   *IndustryStockView
	SalesView      *IndustryStockView
}

// Implements Viewer interface ViewedField method
func (i *IndustryView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements Viewer interface ComparedField method
func (i *IndustryView) ComparedField(f string) string {
	return reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f).String()
}

// Create a single IndustryView for display in a template
//
//	v: the currently viewed industry
//	c: the same industry at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateIndustryView(v *Industry, c *Industry) Viewer {
	return Viewer(&IndustryView{
		viewedRecord:   v,
		comparedRecord: c,
		MoneyView:      &IndustryStockView{v.Money, c.Money},
		SalesView:      &IndustryStockView{v.Sales, c.Sales},
		VariableView:   &IndustryStockView{v.Variable, c.Variable},
		ConstantView:   &IndustryStockView{v.Constant[0], c.Constant[0]},
	})
}

// Create a slice of IndustryView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func IndustryViews(v *[]Industry, c *[]Industry) *[]Viewer {
	var views = make([]Viewer, len(*v))
	for i := range *v {
		view := CreateIndustryView(&(*v)[i], &(*c)[i])
		// vs, _ := json.MarshalIndent(view.(*IndustryView).SalesView, " ", " ")
		// fmt.Printf("Sales View of industry %s is\n %v\n", (*v)[i].Name, string(vs))
		views[i] = view
	}
	return &views
}

// Returns a Viewer for the Money stock of i
func (i *IndustryView) Money() IndustryStockView {
	return IndustryStockView{viewedRecord: i.viewedRecord.Money, comparedRecord: i.comparedRecord.Money}
}

// Returns a Viewer for the Variable stock of i
func (i *IndustryView) Variable() IndustryStockView {
	return IndustryStockView{viewedRecord: i.viewedRecord.Variable, comparedRecord: i.comparedRecord.Variable}
}

// Returns a Viewer for the Sales stock of i
func (i *IndustryView) Sales() IndustryStockView {
	return IndustryStockView{viewedRecord: i.viewedRecord.Sales, comparedRecord: i.comparedRecord.Sales}
}

// Returns a Viewer for the Constant stock of i
// TODO extend to a slice of Constant
func (i *IndustryView) Constant() IndustryStockView {
	return IndustryStockView{viewedRecord: i.viewedRecord.Constant[0], comparedRecord: i.comparedRecord.Constant[0]}
}

// Type for implementation of Viewer interface
type IndustryStockView struct {
	viewedRecord   *IndustryStock
	comparedRecord *IndustryStock
}

// Implements Viewer interface ViewedField method
func (i *IndustryStockView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements Viewer interface ViewedField method
func (i *IndustryStockView) ComparedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Create a single IndustryStockView for display in a template
//
//	v: the currently viewed IndustryStock
//	c: the same IndustryStock at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateIndustryStockView(v *IndustryStock, c *IndustryStock) Viewer {
	return &IndustryStockView{
		viewedRecord:   v,
		comparedRecord: c,
	}
}

// Create a slice of IndustryStockViews for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func IndustryStockViews(v *[]IndustryStock, c *[]IndustryStock) *[]Viewer {
	var newViews = make([]Viewer, len(*v))
	for i := range *v {
		newView := CreateIndustryStockView(&(*v)[i], &(*c)[i])
		newViews[i] = newView
	}
	return &newViews
}

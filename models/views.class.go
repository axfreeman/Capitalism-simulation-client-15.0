package models

import (
	"fmt"
	"reflect"
)

// Type for implementation of Viewer interface
type ClassView struct {
	viewedRecord    *Class
	comparedRecord  *Class
	ConsumptionView *ClassStockView
	MoneyView       *ClassStockView
	SalesView       *ClassStockView
}

// Implements Viewer interface ViewedField method
func (i *ClassView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements Viewer interface ComparedField method
func (i *ClassView) ComparedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Create a single ClassView for display in a template
//
//	v: the currently viewed Class
//	c: the same Class at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateClassView(v *Class, c *Class) Viewer {
	return Viewer(&ClassView{
		viewedRecord:    v,
		comparedRecord:  c,
		MoneyView:       &ClassStockView{v.Money, c.Money},
		SalesView:       &ClassStockView{v.Sales, c.Sales},
		ConsumptionView: &ClassStockView{v.Consumption[0], c.Consumption[0]}, // TODO expand to slice
	})
}

// Create a slice of ClassView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func ClassViews(v *[]Class, c *[]Class) *[]Viewer {
	var views = make([]Viewer, len(*v))
	for i := range *v {
		view := CreateClassView(&(*v)[i], &(*c)[i])
		// vs, _ := json.MarshalIndent(view.(*ClassView).SalesView, " ", " ")
		// fmt.Printf("Sales View of class %s is\n %v\n", (*v)[i].Name, string(vs))
		views[i] = view
	}
	return &views
}

// Type for implementation of Viewer interface
type ClassStockView struct {
	viewedRecord   *ClassStock
	comparedRecord *ClassStock
}

// Implements Viewer interface ViewedField method
func (i *ClassStockView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements Viewer interface ComparedField method
func (i *ClassStockView) ComparedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Create a single ClassStockView for display in a template
//
//	v: the currently viewed ClassStock
//	c: the same ClassStock at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateClassStockView(v *ClassStock, c *ClassStock) Viewer {
	return &ClassStockView{viewedRecord: v, comparedRecord: c}
}

// Create a slice of ClassStockView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func ClassStockViews(v *[]ClassStock, c *[]ClassStock) *[]Viewer {
	var newViews = make([]Viewer, len(*v))
	var vc *ClassStock
	var cc *ClassStock
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateClassStockView(vc, cc)
		newViews[i] = newView
	}
	return &newViews
}

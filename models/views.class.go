package models

import (
	"fmt"
	"gorilla-client/views"
	"html/template"
	"reflect"
)

// Type for implementation of views.Viewer interface
type ClassView struct {
	viewedRecord    *Class
	comparedRecord  *Class
	ConsumptionView *ClassStockView
	MoneyView       *ClassStockView
	SalesView       *ClassStockView
}

// Embedded data for a single class, to pass into templates
type ClassData struct {
	TemplateData
	Class Class
}

// Implements views.Viewer interface ViewedField method
func (i *ClassView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements views.Viewer interface ComparedField method
func (i *ClassView) ComparedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Create a single ClassView for display in a template
//
//	v: the currently viewed Class
//	c: the same Class at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateClassView(v *Class, c *Class) views.Viewer {
	return views.Viewer(&ClassView{
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
func ClassViews(v *[]Class, c *[]Class) *[]views.Viewer {
	var views = make([]views.Viewer, len(*v))
	for i := range *v {
		view := CreateClassView(&(*v)[i], &(*c)[i])
		// vs, _ := json.MarshalIndent(view.(*ClassView).SalesView, " ", " ")
		// fmt.Printf("Sales View of class %s is\n %v\n", (*v)[i].Name, string(vs))
		views[i] = view
	}
	return &views
}

// Type for implementation of views.Viewer interface
type ClassStockView struct {
	viewedRecord   *ClassStock
	comparedRecord *ClassStock
}

// Implements views.Viewer interface ViewedField method
func (i *ClassStockView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements views.Viewer interface ComparedField method
func (i *ClassStockView) ComparedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Create a single ClassDataView to display in the class.html template.
// This is added dynamically to the DisplayData template when requested
//
//	u: the user
//	message: any message
//	id: the id of the social class to display
//
//	returns: classData which references this class, and embeds an OutputData
func (u User) ClassDisplayData(message string, id int) ClassData {
	return ClassData{
		u.CreateTemplateData(message),
		*ViewedObject[Class](u, `classes`, id),
	}
}

// Create a single ClassStockView for display in a template
//
//	v: the currently viewed ClassStock
//	c: the same ClassStock at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateClassStockView(v *ClassStock, c *ClassStock) views.Viewer {
	return &ClassStockView{viewedRecord: v, comparedRecord: c}
}

// Create a slice of ClassStockView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func ClassStockViews(v *[]ClassStock, c *[]ClassStock) *[]views.Viewer {
	var newViews = make([]views.Viewer, len(*v))
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

// Returns a safe HTML string with a link to the Commodity of a class
// Should be a method of IndustryView but haven't yet figured out how to fix this
//
//	v: Industry implementation of the views.Viewer interface
//	template.HTML: safe string using fields supplied by the Commodity implementation
func ClassCommodityLink(v ClassView) template.HTML {
	o := v.viewedRecord
	output := template.HTML(fmt.Sprintf(`<td><a href="/commodity/%d">%s</a></td>`, o.Commodity.Id, o.Output))
	// utils.TraceInfof(utils.Purple, "Industry Commodity Link says commodity Id is %s", string(output))
	return output
}

// Returns a safe HTML string with a link to stock's class
//
//	v: an implementation of the views.Viewer interface
//	urlBase: the root of the link url (eg `commodity`)
//	template.HTML: safe string using ID and Name fields supplied by the implementation
func StockClassLink(v ClassStockView) template.HTML {
	o := v.viewedRecord
	className := o.ClassName
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%d">%s</a>`, `class`, o.ClassId, className))
}

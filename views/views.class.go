package views

import (
	"fmt"
	"html/template"
	"reflect"
	"simulation-client/models"
)

// Type for implementation of Viewer interface
type ClassView struct {
	viewedRecord    *models.Class
	comparedRecord  *models.Class
	ConsumptionView *ClassStockView
	MoneyView       *ClassStockView
	SalesView       *ClassStockView
}

// Embedded data for a single class, to pass into templates
type ClassData struct {
	TemplateData
	Class models.Class
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
func CreateClassView(v *models.Class, c *models.Class) Viewer {
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
func ClassViews(v *[]models.Class, c *[]models.Class) *[]Viewer {
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
	viewedRecord   *models.ClassStock
	comparedRecord *models.ClassStock
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

// Diagnostic method exposes base viewed record to the Viewer interface
func (v *ClassView) Viewed() any {
	return v.viewedRecord
}

// Diagnostic method exposes base compared record to the Viewer interface
func (c *ClassView) Compared() any {
	return c.comparedRecord
}

// Create a single ClassDataView to display in the class.html template.
// This is added dynamically to the DisplayData template when requested
//
//	u: the user
//	message: any message
//	id: the id of the social class to display
//
//	returns: classData which references this class, and embeds an OutputData
func ClassDisplayData(u *models.User, message string, id int) ClassData {
	return ClassData{
		CreateTemplateData(u, message),
		*ViewedObject[models.Class](*u, `classes`, id),
	}
}

// Create a single ClassStockView for display in a template
//
//	v: the currently viewed ClassStock
//	c: the same ClassStock at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateClassStockView(v *models.ClassStock, c *models.ClassStock) Viewer {
	return &ClassStockView{viewedRecord: v, comparedRecord: c}
}

// Diagnostic method exposes base viewed record to the Viewer interface
func (v *ClassStockView) Viewed() any {
	return v.viewedRecord
}

// Diagnostic method exposes base compared record to the Viewer interface
func (c *ClassStockView) Compared() any {
	return c.comparedRecord
}

// Create a slice of ClassStockView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func ClassStockViews(v *[]models.ClassStock, c *[]models.ClassStock) *[]Viewer {
	var newViews = make([]Viewer, len(*v))
	var vc *models.ClassStock
	var cc *models.ClassStock
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateClassStockView(vc, cc)
		newViews[i] = newView
	}
	return &newViews
}

// Embedded data for a single ClassStock, to pass into templates
type ClassStockData struct {
	TemplateData
	ClassStock models.ClassStock
}

// Create a ClassStockData to display a single classStock in the
// class-stock.html template. This is added dynamically to the DisplayData
// template when the class-stock view is requested
//
//	u: the user
//	message: any message
//	id: the id of the class item to display
//
//	returns: industryStockData which references this industryStock, and embeds a TemplateData
func ClassStockDisplayData(u *models.User, message string, id int) ClassStockData {
	return ClassStockData{
		CreateTemplateData(u, message), *ViewedObject[models.ClassStock](*u, `class_stocks`, id),
	}
}

// Returns a safe HTML string with a link to the Commodity of a class
// Should be a method of IndustryView but haven't yet figured out how to fix this
//
//	v: Industry implementation of the Viewer interface
//	template.HTML: safe string using fields supplied by the Commodity implementation
func ClassCommodityLink(v ClassView) template.HTML {
	o := v.viewedRecord
	output := template.HTML(fmt.Sprintf(`<td><a href="/commodity/%d">%s</a></td>`, o.Commodity.Id, o.Output))
	// logging.TraceInfof(utils.Purple, "Industry Commodity Link says commodity Id is %s", string(output))
	return output
}

// Returns a safe HTML string with a link to stock's class
//
//	v: an implementation of the Viewer interface
//	urlBase: the root of the link url (eg `commodity`)
//	template.HTML: safe string using ID and Name fields supplied by the implementation
func StockClassLink(v ClassStockView) template.HTML {
	o := v.viewedRecord
	className := o.ClassName
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%d">%s</a></td>`, `class`, o.ClassId, className))
}

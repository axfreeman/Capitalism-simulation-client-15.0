package models

import (
	"fmt"
	"html/template"
	"reflect"
	"strconv"
)

// TODO Show should display decimals when required
// TODO figure out how to make Graphics a method of the implementation, not the interface

// Interface for all view types. Wrapped by the view struct to
// provide the 'Show' method, which compares a viewed field at
// the current stage of the simulation, with a compared field
// at a previous stage.
//
//	viewedField(f): returns the field f of a viewed record
//	comparedField(f): returns the field f of a compared record
type Viewer interface {
	ViewedField(f string) string
	ComparedField(f string) string
}

// Wrapper for the Viewer struct. Any view that implements Viewer
// can access the Show method of this type
type View struct {
	Viewer
}

// Provide a string, suitable for display in a template, that formats
// a viewed value and highlights values that have changed.
//
//	v: a View object
//	f: the name of the field to display
//	Returns: safe HTML string coloured red if the value has changed
func Show(v Viewer, f string) template.HTML {
	// fmt.Printf("   Entered Show with field %s\n", f)
	vv, _ := strconv.Atoi(v.ViewedField(f))
	vc, _ := strconv.Atoi(v.ComparedField(f))
	var htmlString string
	if vv == vc {
		htmlString = fmt.Sprintf("<td style=\"text-align:center\">%d</td>", vv)
	} else {
		htmlString = fmt.Sprintf("<td style=\"text-align:center; color:red\">%d</td>", vv)
	}
	return template.HTML(htmlString)
}

// Returns a safe HTML string with a link to the ViewedField object
// Assumes the implementation supplies Name and ID fields
//
// v: an implementation of the Viewer interface
// urlBase: the root of the link url (eg `commodity`)
// template.HTML: safe string using ID and Name fields supplied by the implementation
func Link(v Viewer, urlBase string) template.HTML {
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%s\">%s</a>`, urlBase, v.ViewedField(`Id`), v.ViewedField(`Name`)))
}

// Returns a safe HTML string with a link to the Commodity of an industry
// Should be a method of  IndustryView but haven't yet figured out how
//
// v: an implementation of the Viewer interface
// template.HTML: safe string using ID and Name fields supplied by the implementation
func CommodityLink(v Viewer) template.HTML {
	return template.HTML(fmt.Sprintf(`<td><a href="/commodity/%s">%s</a></td>`, v.ViewedField(`OutputCommodityID`), v.ViewedField("Output")))
}

// Returns a safe HTML string with a graphic illustrating the origin
//
//	v: a CommodityView
//	template.HTML: safe string with a graphic representing the origin
func OriginGraphic(v Viewer) template.HTML {
	var htmlString template.HTML
	switch v.ViewedField(`Origin`) {
	case `INDUSTRIAL`:
		htmlString = `<td style="text-align:center"><i style="font-weight: bolder; color:blue" class="fa fa-industry"></i></td>`
	case `SOCIAL`:
		if v.ViewedField(`Usage`) == `Useless` {
			htmlString = `<td style="text-align:center"><i style="font-weight: bolder; color:rgba(128, 0, 128, 0.696)" class="fas fa-user-tie"></i></td>`
		} else {
			htmlString = `<td style="text-align:center"><i style="font-weight: bolder; color:red" class="fa fa-user-friends"></i></td>`
		}
	case `MONEY`:
		htmlString = `<td style="text-align:center"><i style="font-weight: 900; color:goldenrod" class="fa fa-dollar"></i></td>`
	default:
		htmlString = `<td style="text-align:center">Unknown Origin</td>`
	}
	return template.HTML(htmlString)
}

// Returns a safe HTML string with a graphic illustrating the usage
//
//	v: a CommodityView
//	template.HTML: safe string with a graphic representing the usage
func UsageGraphic(v Viewer) template.HTML {
	var htmlString template.HTML
	switch v.ViewedField(`Usage`) {
	case `PRODUCTIVE`:
		htmlString = `<td style="text-align:center"><i style="font-weight: bolder; color:blue" class="fas fa-hammer"></i></td>`
	case `CONSUMPTION`:
		htmlString = `<td style="text-align:center"><i style="font-weight: bolder; color:green" class="fa fa-cutlery"></i></td>`
	case `MONEY`:
		htmlString = `<td style="text-align:center"><i class="fa fa-dollar" style="font-weight: 900; color:goldenrod"></i></td>`
	case `Useless`:
		htmlString = `<td style="text-align:center"><i class="fas fa-skull-crossbones" style="font-weight: bolder; color:black"></i></td>`
	default:
		htmlString = `<td style="text-align:center">Unknown Usage</td>`
	}
	return template.HTML(htmlString)
}

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
func CreateIndustryView(vi *Industry, ci *Industry) Viewer {
	moneyView := IndustryStockView{vi.Money, ci.Money}
	salesView := IndustryStockView{vi.Sales, ci.Sales}
	variableView := IndustryStockView{vi.Variable, ci.Variable}
	constantView := IndustryStockView{vi.Constant[0], ci.Constant[0]}

	return Viewer(&IndustryView{
		viewedRecord:   vi,
		comparedRecord: ci,
		MoneyView:      &moneyView,
		SalesView:      &salesView,
		VariableView:   &variableView,
		ConstantView:   &constantView,
	})
}

// Create a slice of IndustryView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func IndustryViews(v *[]Industry, c *[]Industry) *[]Viewer {
	var views = make([]Viewer, len(*v))
	var vi *Industry
	var ci *Industry
	for i := range *v {
		vi = &(*v)[i]
		ci = &(*c)[i]

		view := Viewer(CreateIndustryView(vi, ci))
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
	return reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f).String()
}

// Create a single IndustryStockView for display in a template
//
//	v: the currently viewed IndustryStock
//	c: the same IndustryStock at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateIndustryStockView(v *IndustryStock, c *IndustryStock) Viewer {
	return View{&IndustryStockView{
		viewedRecord:   v,
		comparedRecord: c,
	}}
}

// Create a slice of IndustryStockViews for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func IndustryStockViews(v *[]IndustryStock, c *[]IndustryStock) *[]Viewer {
	var newViews = make([]Viewer, len(*v))
	var vc *IndustryStock
	var cc *IndustryStock
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateIndustryStockView(vc, cc)
		newViews[i] = newView
	}
	return &newViews
}

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
	return reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f).String()
}

// Create a single ClassView for display in a template
//
//	v: the currently viewed Class
//	c: the same Class at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateClassView(v *Class, c *Class) Viewer {
	return View{&ClassView{
		viewedRecord:    v,
		comparedRecord:  c,
		MoneyView:       &ClassStockView{v.Money, c.Money},
		SalesView:       &ClassStockView{v.Sales, c.Sales},
		ConsumptionView: &ClassStockView{v.Consumption[0], c.Consumption[0]}, // TODO expand to slice
	}}
}

// Create a slice of ClassView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func ClassViews(v *[]Class, c *[]Class) *[]Viewer {
	var newViews = make([]Viewer, len(*v))
	var vc *Class
	var cc *Class
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateClassView(vc, cc)
		newViews[i] = newView
	}
	return &newViews
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
	return reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f).String()
}

// Create a single ClassStockView for display in a template
//
//	v: the currently viewed ClassStock
//	c: the same ClassStock at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateClassStockView(v *ClassStock, c *ClassStock) Viewer {
	return View{&ClassStockView{
		viewedRecord:   v,
		comparedRecord: c,
	}}
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

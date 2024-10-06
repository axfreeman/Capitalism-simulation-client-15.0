package models

import (
	"fmt"
	"gorilla-client/utils"
	"gorilla-client/views"
	"html/template"
	"reflect"
)

// Type for implementation of views.Viewer interface
type IndustryView struct {
	viewedRecord   *Industry
	comparedRecord *Industry
	MoneyView      *IndustryStockView
	VariableView   *IndustryStockView
	ConstantView   *IndustryStockView
	SalesView      *IndustryStockView
}

// Implements views.Viewer interface ViewedField method
func (i *IndustryView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements views.Viewer interface ComparedField method
func (i *IndustryView) ComparedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Create a single IndustryView for display in a template
//
//	v: the currently viewed industry
//	c: the same industry at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateIndustryView(v *Industry, c *Industry) views.Viewer {
	// fmt.Printf("Entering CreateIndustryView\n%v\n%v\n%v\n%v\n", v, c, v.Constant, c.Constant)
	return views.Viewer(&IndustryView{
		viewedRecord:   v,
		comparedRecord: c,
		MoneyView:      &IndustryStockView{v.Money, c.Money},
		SalesView:      &IndustryStockView{v.Sales, c.Sales},
		VariableView:   &IndustryStockView{v.Variable, c.Variable},
		ConstantView:   &IndustryStockView{v.Constant[0], c.Constant[0]},
	})
}

// Diagnostic method exposes base viewed record to the Viewer interface
func (v *IndustryView) Viewed() any {
	return v.viewedRecord
}

// Diagnostic method exposes base compared record to the Viewer interface
func (c *IndustryView) Compared() any {
	return c.comparedRecord
}

// Diagnostic method exposes base viewed record to the Viewer interface
func (v *IndustryStockView) Viewed() any {
	return v.viewedRecord
}

// Diagnostic method exposes base compared record to the Viewer interface
func (c *IndustryStockView) Compared() any {
	return c.comparedRecord
}

// Create a slice of IndustryView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func IndustryViews(v *[]Industry, c *[]Industry) *[]views.Viewer {
	var views = make([]views.Viewer, len(*v))
	for i := range *v {
		view := CreateIndustryView(&(*v)[i], &(*c)[i])
		// vs, _ := json.MarshalIndent(view.(*IndustryView).SalesView, " ", " ")
		// fmt.Printf("Sales View of industry %s is\n %v\n", (*v)[i].Name, string(vs))
		views[i] = view
	}
	return &views
}

// Returns a views.Viewer for the Money stock of i
func (i *IndustryView) Money() IndustryStockView {
	return IndustryStockView{viewedRecord: i.viewedRecord.Money, comparedRecord: i.comparedRecord.Money}
}

// Returns a views.Viewer for the Variable stock of i
func (i *IndustryView) Variable() IndustryStockView {
	return IndustryStockView{viewedRecord: i.viewedRecord.Variable, comparedRecord: i.comparedRecord.Variable}
}

// Returns a views.Viewer for the Sales stock of i
func (i *IndustryView) Sales() IndustryStockView {
	return IndustryStockView{viewedRecord: i.viewedRecord.Sales, comparedRecord: i.comparedRecord.Sales}
}

// Returns a views.Viewer for the Constant stock of i
// TODO extend to a slice of Constant
func (i *IndustryView) Constant() IndustryStockView {
	return IndustryStockView{viewedRecord: i.viewedRecord.Constant[0], comparedRecord: i.comparedRecord.Constant[0]}
}

// Type for implementation of views.Viewer interface
type IndustryStockView struct {
	viewedRecord   *IndustryStock
	comparedRecord *IndustryStock
}

// Implements views.Viewer interface ViewedField method
// ** DIAGNOSTICS SEE https://stackoverflow.com/questions/17262238/how-to-cast-reflect-value-to-its-type Last answer ***
func (i *IndustryStockView) ViewedField(f string) string {
	// utils.TraceInfof(utils.Yellow, "  Entered ViewedField for IndustryStockView with f=%s", f)
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	// Diagnostics - probably not needed now...
	if f == `Size` {
		utils.TraceInfof(utils.Yellow, "Displaying an IndustryStockView with f=%s", f)
		r := reflect.ValueOf(i.viewedRecord)
		in := reflect.Indirect(r)
		record := in.Interface().(IndustryStock)
		st := record.Write()
		fmt.Println("***The result is ", s)
		fmt.Printf("***The result formatted is %v\n", s)
		utils.TraceInfof(utils.Yellow, "Stock is:\n%v", st)
		sf := fmt.Sprint(s)
		fmt.Println("***The result sprinted is", sf)
	}
	// ...End of diagnostics
	return fmt.Sprint(s)
}

// Implements views.Viewer interface ViewedField method
func (i *IndustryStockView) ComparedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Create a single IndustryStockView for display in a template
//
//	v: the currently viewed IndustryStock
//	c: the same IndustryStock at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateIndustryStockView(v *IndustryStock, c *IndustryStock) views.Viewer {
	return &IndustryStockView{
		viewedRecord:   v,
		comparedRecord: c,
	}
}

// Embedded data for a single industry, to pass into templates
type IndustryData struct {
	TemplateData
	Industry Industry
}

// Create an IndustryData to display a single industry in the
// industry.html template. This is added dynamically to the DisplayData
// template when the Commodity view is requested
//
//	u: the user
//	message: any message
//	id: the id of the industry item to display
//
//	returns: industryData which references this industry, and embeds an OutputData
func (u User) IndustryDisplayData(message string, id int) IndustryData {
	return IndustryData{
		u.CreateTemplateData(message),
		*ViewedObject[Industry](u, `industries`, id),
	}
}

// Create a slice of IndustryStockViews for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func IndustryStockViews(v *[]IndustryStock, c *[]IndustryStock) *[]views.Viewer {
	var newViews = make([]views.Viewer, len(*v))
	for i := range *v {
		newView := CreateIndustryStockView(&(*v)[i], &(*c)[i])
		newViews[i] = newView
	}
	return &newViews
}

// Embedded data for a single IindustryStock, to pass into templates
type IndustryStockData struct {
	TemplateData
	IndustryStock IndustryStock
}

// Create an IndustryStockData to display a single industryStock in the
// industry-stock.html template. This is added dynamically to the DisplayData
// template when the industry-stock view is requested
//
//	u: the user
//	message: any message
//	id: the id of the industry item to display
//
//	returns: industryStockData which references this industryStock, and embeds a TemplateData
func (u User) IndustryStockDisplayData(message string, id int) IndustryStockData {
	return IndustryStockData{
		u.CreateTemplateData(message),
		*ViewedObject[IndustryStock](u, `industry_stocks`, id),
	}
}

// Implementation-specific template methods

// Returns a safe HTML string with a link to the Commodity of an industry
// Should be a method of IndustryView but haven't yet figured out how to fix this
//
//	v: Industry implementation of the views.Viewer interface
//	template.HTML: safe string using fields supplied by the Commodity implementation
func IndustryCommodityLink(v IndustryView) template.HTML {
	o := v.viewedRecord
	output := template.HTML(fmt.Sprintf(`<td><a href="/commodity/%d">%s</a></td>`, o.Commodity.Id, o.Output))
	// utils.TraceInfof(utils.Purple, "Industry Commodity Link says commodity Id is %s", string(output))
	return output
}

// Returns a safe HTML string with a link to industry stock's commodity
//
//	v: IndustryStock implementation of the views.Viewer interface
//	urlBase: the root of the link url (eg `commodity`)
//	template.HTML: safe string using ID and Name fields supplied by the implementation
func IndustryStockCommodityLink(v IndustryStockView) template.HTML {
	o := v.viewedRecord
	commodityName := o.CommodityName
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%s">%s</a>`, `commodity`, v.ViewedField(`CommodityId`), commodityName))
}

// Returns a safe HTML string with a link to industry stock's industry
//
//	v: an implementation of the views.Viewer interface
//	urlBase: the root of the link url (eg `commodity`)
//	template.HTML: safe string using ID and Name fields supplied by the implementation
func StockIndustryLink(v IndustryStockView) template.HTML {
	o := v.viewedRecord
	industryName := o.IndustryName
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%d">%s</a>`, `industry`, o.IndustryId, industryName))
}

// Returns a safe HTML string with a link to industry stock's commodity
//
//	v: IndustryStock implementation of the views.Viewer interface
//	urlBase: the root of the link url (eg `commodity`)
//	template.HTML: safe string using ID and Name fields supplied by the implementation
func ClassStockCommodityLink(v ClassStockView) template.HTML {
	o := v.viewedRecord
	commodityName := o.CommodityName
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%s">%s</a>`, `commodity`, v.ViewedField(`CommodityId`), commodityName))
}

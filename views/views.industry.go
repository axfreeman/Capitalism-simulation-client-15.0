package views

import (
	"fmt"
	"html/template"
	"reflect"
	"simulation-client/models"
)

// Type for implementation of Viewer interface
type IndustryView struct {
	viewedRecord   *models.Industry
	comparedRecord *models.Industry
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
	s := reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Create a single IndustryView for display in a template
//
//	v: the currently viewed industry
//	c: the same industry at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateIndustryView(v *models.Industry, c *models.Industry) Viewer {
	// fmt.Printf("Entering CreateIndustryView\n%v\n%v\n%v\n%v\n", v, c, v.Constant, c.Constant)
	return Viewer(&IndustryView{
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
func IndustryViews(v *[]models.Industry, c *[]models.Industry) *[]Viewer {
	var views = make([]Viewer, len(*v))
	for i := range *v {
		view := CreateIndustryView(&(*v)[i], &(*c)[i])
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
	viewedRecord   *models.IndustryStock
	comparedRecord *models.IndustryStock
}

// Implements Viewer interface ViewedField method
// ** DIAGNOSTICS SEE https://stackoverflow.com/questions/17262238/how-to-cast-reflect-value-to-its-type Last answer ***
func (i *IndustryStockView) ViewedField(f string) string {
	// logging.TraceInfof(utils.Yellow, "  Entered ViewedField for IndustryStockView with f=%s", f)
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	// Diagnostics - probably not needed now...
	// if f == `Size` {
	// 	logging.TraceInfof(utils.Yellow, "Displaying an IndustryStockView with f=%s", f)
	// 	r := reflect.ValueOf(i.viewedRecord)
	// 	in := reflect.Indirect(r)
	// 	record := in.Interface().(IndustryStock)
	// 	st := record.Write()
	// 	fmt.Println("***The result is ", s)
	// 	fmt.Printf("***The result formatted is %v\n", s)
	// 	logging.TraceInfof(utils.Yellow, "Stock is:\n%v", st)
	// 	sf := fmt.Sprint(s)
	// 	fmt.Println("***The result sprinted is", sf)
	// }
	// ...End of diagnostics
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
func CreateIndustryStockView(v *models.IndustryStock, c *models.IndustryStock) Viewer {
	return &IndustryStockView{
		viewedRecord:   v,
		comparedRecord: c,
	}
}

// Embedded data for a single industry, to pass into templates
type IndustryData struct {
	TemplateData
	Industry models.Industry
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
func IndustryDisplayData(u *models.User, message string, id int) IndustryData {
	return IndustryData{
		CreateTemplateData(u, message), *ViewedObject[models.Industry](*u, `industries`, id),
	}
}

// Create a slice of IndustryStockViews for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func IndustryStockViews(v *[]models.IndustryStock, c *[]models.IndustryStock) *[]Viewer {
	var newViews = make([]Viewer, len(*v))
	for i := range *v {
		newView := CreateIndustryStockView(&(*v)[i], &(*c)[i])
		newViews[i] = newView
	}
	return &newViews
}

// Embedded data for a single IindustryStock, to pass into templates
type IndustryStockData struct {
	TemplateData
	IndustryStock models.IndustryStock
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
func IndustryStockDisplayData(u *models.User, message string, id int) IndustryStockData {
	return IndustryStockData{
		CreateTemplateData(u, message), *ViewedObject[models.IndustryStock](*u, `industry_stocks`, id),
	}
}

// Implementation-specific template methods

// Returns a safe HTML string with a link to the Commodity of an industry
// Should be a method of IndustryView but haven't yet figured out how to fix this
//
//	v: Industry implementation of the Viewer interface
//	template.HTML: safe string using fields supplied by the Commodity implementation
func IndustryCommodityLink(v IndustryView) template.HTML {
	o := v.viewedRecord
	output := template.HTML(fmt.Sprintf(`<td><a href="/commodity/%d">%s</a></td>`, o.Commodity.Id, o.Output))
	// logging.TraceInfof(utils.Purple, "Industry Commodity Link says commodity Id is %s", string(output))
	return output
}

// Returns a safe HTML string with a link to industry stock's commodity
//
//	v: IndustryStock implementation of the Viewer interface
//	urlBase: the root of the link url (eg `commodity`)
//	template.HTML: safe string using ID and Name fields supplied by the implementation
func IndustryStockCommodityLink(v IndustryStockView) template.HTML {
	o := v.viewedRecord
	commodityName := o.CommodityName
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%s">%s</a>`, `commodity`, v.ViewedField(`CommodityId`), commodityName))
}

// Returns a safe HTML string with a link to industry stock's industry
//
//	v: an implementation of the Viewer interface
//	urlBase: the root of the link url (eg `commodity`)
//	template.HTML: safe string using ID and Name fields supplied by the implementation
func StockIndustryLink(v IndustryStockView) template.HTML {
	o := v.viewedRecord
	industryName := o.IndustryName
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%d">%s</a>`, `industry`, o.IndustryId, industryName))
}

// Returns a safe HTML string with a link to industry stock's commodity
//
//	v: IndustryStock implementation of the Viewer interface
//	urlBase: the root of the link url (eg `commodity`)
//	template.HTML: safe string using ID and Name fields supplied by the implementation
func ClassStockCommodityLink(v ClassStockView) template.HTML {
	o := v.viewedRecord
	commodityName := o.CommodityName
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%s">%s</a>`, `commodity`, v.ViewedField(`CommodityId`), commodityName))
}

package models

import (
	"fmt"
	"gorilla-client/views"
	"html/template"
	"reflect"
)

// implements View for the Commodity Object
type CommodityView struct {
	viewedRecord   *Commodity
	comparedRecord *Commodity
}

// Provides the value of the field f in the viewedRecord of a CommodityView
//
//	f: the name of a field
//	c: a CommodityView
//	returns: the stringified value of the field (easiest generic solution)
func (c *CommodityView) ViewedField(f string) string {
	return fmt.Sprint(reflect.Indirect(reflect.ValueOf(c.viewedRecord)).FieldByName(f))
}

// Provides the value of the field f in the comparedRecord of a CommodityView
//
//	f: the name of a field
//	c: a CommodityView
//	returns: the stringified value of the field (easiest generic solution)
func (c *CommodityView) ComparedField(f string) string {
	return fmt.Sprint(reflect.Indirect(reflect.ValueOf(c.comparedRecord)).FieldByName(f))
}

// Create a single CommodityView for display in a template
//
//	v: the currently viewed commodity
//	c: the same commodity at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateCommodityView(v *Commodity, c *Commodity) views.Viewer {
	return &CommodityView{
		viewedRecord:   v,
		comparedRecord: c,
	}
}

// Diagnostic method exposes base viewed record to the Viewer interface
func (v *CommodityView) Viewed() any {
	return v.viewedRecord
}

// Diagnostic method exposes base compared record to the Viewer interface
func (c *CommodityView) Compared() any {
	return c.comparedRecord
}

// Embedded data for a single commodity, to pass into templates
type CommodityData struct {
	TemplateData
	Commodity Commodity
}

// Create a CommodityData to display a single commodity in the
// commodity.html template. This is added dynamically to the DisplayData
// template when the Commodity view is requested
//
//	u: the user
//	message: any message
//	id: the id of the commodity to display
//
//	returns: CommodityData which references this commodity, and embeds an OutputData
func CommodityDisplayData(u *User, message string, id int) CommodityData {
	return CommodityData{
		u.CreateTemplateData(message),
		*ViewedObject[Commodity](*u, `commodities`, id),
	}
}

// Create a slice of CommodityView for display in a template
//
//	v: a slice of all commodities in the simulation at the current stage
//	c: a slice of the same commodities at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func CommodityViews(v *[]Commodity, c *[]Commodity) *[]views.Viewer {
	var view = make([]views.Viewer, len(*v))
	var vc *Commodity
	var cc *Commodity
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateCommodityView(vc, cc)
		view[i] = newView
	}
	return &view
}

// Returns a safe HTML string with a graphic illustrating the origin
//
//	v: a CommodityView
//	template.HTML: safe string with a graphic representing the origin
func OriginGraphic(v views.Viewer) template.HTML {
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
func UsageGraphic(v views.Viewer) template.HTML {
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

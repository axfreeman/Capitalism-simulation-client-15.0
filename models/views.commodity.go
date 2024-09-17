package models

import (
	"fmt"
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
func CreateCommodityView(v *Commodity, c *Commodity) Viewer {
	return &CommodityView{
		viewedRecord:   v,
		comparedRecord: c,
	}
}

// Create a slice of CommodityView for display in a template
//
//	v: a slice of all commodities in the simulation at the current stage
//	c: a slice of the same commodities at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func CommodityViews(v *[]Commodity, c *[]Commodity) *[]Viewer {
	var view = make([]Viewer, len(*v))
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

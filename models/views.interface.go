package models

import (
	"fmt"
	"html/template"
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

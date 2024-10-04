package views

import (
	"fmt"
	"gorilla-client/utils"
	"html/template"
	"strconv"
)

// TODO Show should display decimals when required

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
	Viewed() any
	Compared() any
}

// Provide a string, suitable for display in a template, that formats
// a viewed value and highlights values that have changed.
//
//	v: a View object
//	f: the name of the field to display
//	Returns: safe HTML string coloured red if the value has changed
func Show(v Viewer, f string) template.HTML {
	vv, _ := strconv.Atoi(v.ViewedField(f))
	vc, _ := strconv.Atoi(v.ComparedField(f))
	// Diagnostics - turn on if problems with display ...
	if f == "Size" {
		firstPart := utils.TraceInfoPart(utils.Yellow, " Viewed %d, Compared %d", vv, vc)
		secondPart := utils.TraceInfoPart(utils.Cyan, " Viewed Id %s, Compared Id %s", v.ViewedField(`Id`), v.ComparedField(`Id`))
		utils.TraceInfof(utils.BrightYellow, "Show %s: %s [%s]", f, firstPart, secondPart)
	}
	// ...End of Diagnostics

	var htmlString string
	if vv == vc {
		htmlString = fmt.Sprintf("<td style=\"text-align:center\">%d</td>", vv)
	} else {
		htmlString = fmt.Sprintf("<td style=\"text-align:center; color:red\">%d</td>", vv)
	}
	return template.HTML(htmlString)
}

// Provide a string, suitable for display in a template, that formats
// a decimal viewed value and highlights values that have changed.
//
//	v: a View object
//	f: the name of the field to display
//	Returns: safe HTML string coloured red if the value has changed
func ShowDecimal(v Viewer, f string) template.HTML {
	sv, _ := strconv.ParseFloat(v.ViewedField(f), 32)
	sc, _ := strconv.ParseFloat(v.ComparedField(f), 32)

	var htmlString string
	if sv == sc {
		htmlString = fmt.Sprintf("<td style=\"text-align:center\">%0.2f</td>", sv)
	} else {
		htmlString = fmt.Sprintf("<td style=\"text-align:center; color:red\">%0.2f</td>", sc)
	}
	return template.HTML(htmlString)
}

// Provide a string representing the named field
//
//	v: a View object
//	f: the name of the field to display
//	Returns: safe HTML string
func ShowString(v Viewer, f string) template.HTML {
	return template.HTML(fmt.Sprintf("<td style=\"text-align:center\">%s</td>", v.ViewedField(f)))
}

// Returns a safe HTML string with a link to the ViewedField object
// Assumes the implementation supplies Name and ID fields
//
// v: an implementation of the views.Viewer interface
// urlBase: the root of the link url (eg `commodity`)
// f: the field to display
// template.HTML: safe string using ID and Name fields supplied by the implementation
func Link(v Viewer, urlBase string, f string) template.HTML {
	return template.HTML(fmt.Sprintf(`<td style="text-align:left"><a href="/%s/%s">%s</a>`, urlBase, v.ViewedField(`Id`), v.ViewedField(f)))
}

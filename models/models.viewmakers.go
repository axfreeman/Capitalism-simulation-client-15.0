package models

import (
	"gorilla-client/utils"
	"reflect"
)

// Contains the generic view constructor which uses reflect
type Record interface{}

var PairType reflect.Type
var StringType reflect.Type

// Used in the View object
// Allows a template to distinguish magnitudes that have changed
// (eg by displaying them in a different colour)
//
//	Viewed: the current field in the simulation
//	Compared: the same field earlier in the simulation
type Pair struct {
	Viewed   float32
	Compared float32
}

// Create frequently-used constants
func InitViews() {
	PairType = reflect.TypeOf(Pair{})
	StringType = reflect.TypeOf("")
}

// Populate a View which is then passed into a template for display. The
// parameter is an unpopulated View object containing two Record objects,
// contained in the RecordBase embedded field, and called Viewed and Compared.
//
// Each 'Pair' field in the view is given a Pair object containing the
// corresponding viewed and compared magnitudes in Viewed and Compared
//
// NOTE: fields in RecordBase must have the same names as in View
//
//		view: the view to be populated
//		T: a struct whose type is specified by the Record interface
//	    This will be one of the main objects of the simulation, viz:
//	    Commodity | Simulation | Class
func PopulateView[T Record](View *T) {
	recordBase := reflect.ValueOf(*View).FieldByName("RecordBase")
	viewedRecord := reflect.Indirect(recordBase.FieldByName("Viewed"))
	comparedRecord := reflect.Indirect(recordBase.FieldByName("Compared"))
	vPtr := reflect.ValueOf(View)  // Pointer to the view
	vVal := reflect.ValueOf(*View) // Dereferenced Copy of the view
	vValTyp := vVal.Type()
	vElem := vPtr.Elem()

	for i := 0; i < vVal.NumField(); i++ {
		f := vPtr.Elem().Field(i)
		fTyp := f.Type()
		fieldFromValType := vValTyp.Field(i)
		name := fieldFromValType.Name

		switch fTyp {
		case StringType:
			x := viewedRecord.FieldByName(name)
			v := vElem.Field(i)
			v.Set(x)

		case PairType:
			f0 := f.Field(0)
			f1 := f.Field(1)
			vm := viewedRecord.FieldByName(name)
			cm := comparedRecord.FieldByName(name)
			if vm.IsValid() {
				f0.Set(vm)
				f1.Set(cm)
			} else {
				utils.TraceInfof(utils.Red, "Producing view item for field called %s", name)
				vmbn := viewedRecord.MethodByName(name)
				if vmbn.IsValid() {
					cmbn := comparedRecord.MethodByName(name) // if vmbn is valid, we can safely assume cmbn is too
					in := make([]reflect.Value, 0)
					vval := vmbn.Call(in)
					cval := cmbn.Call(in)
					f0.Set(vval[0])
					f1.Set(cval[0])
				}
			}
		default:
		}
	}
}

// Create a CommodityView for display in a template
// taking data from two Commodity objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func CommodityView(v *Commodity, c *Commodity) *CommodityViewer {
	recordBase := RecordBase[Commodity]{
		Viewed:   v,
		Compared: c,
	}
	view := CommodityViewer{RecordBase: recordBase}
	PopulateView(&view)
	return &view
}

// Create a slice of CommodityView for display in a template
// taking data from two Commodity objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func CommodityViews(v *[]Commodity, c *[]Commodity) *[]CommodityViewer {
	var newViews = make([]CommodityViewer, len(*v))
	for i := range *v {
		newView := CommodityView(&(*v)[i], &(*c)[i])
		newViews[i] = *newView
	}
	return &newViews
}

// Create an IndustryView for display in a template
// taking data from two Industry objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func IndustryView(v *Industry, c *Industry) *IndustryViewer {
	recordBase := RecordBase[Industry]{
		Viewed:   v,
		Compared: c,
	}
	view := IndustryViewer{RecordBase: recordBase}
	PopulateView(&view)
	return &view
}

// Create a slice of IndustryView for display in a template
// taking data from two Commodity objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func IndustryViews(v *[]Industry, c *[]Industry) *[]IndustryViewer {
	utils.TraceInfof(utils.BrightRed, "Creating an industry view of length %d", len(*v))
	var newViews = make([]IndustryViewer, len(*v))
	for i := range *v {
		newView := IndustryView(&(*v)[i], &(*c)[i])
		newViews[i] = *newView
	}
	return &newViews
}

// Create an ClassView for display in a template
// taking data from two Class objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func ClassView(v *Class, c *Class) *ClassViewer {
	recordBase := RecordBase[Class]{
		Viewed:   v,
		Compared: c,
	}
	view := ClassViewer{RecordBase: recordBase}
	PopulateView(&view)
	return &view
}

// Create a slice of ClassView for display in a template
// taking data from two Commodity objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func ClassViews(v *[]Class, c *[]Class) *[]ClassViewer {
	var newViews = make([]ClassViewer, len(*v))
	for i := range *v {
		newView := ClassView(&(*v)[i], &(*c)[i])
		newViews[i] = *newView
	}
	return &newViews
}

package models

import "reflect"

// Contains the generic view constructor which uses reflect
type Record interface{}

var PairType reflect.Type
var StringType reflect.Type

// Used in the View object
// Allows a template to distinguish magnitudes that have changed
// (eg by displaying them in a different colour)
//
//	 Viewed: the current field in the simulation
//	 Compared: the same field earlier in the simulation
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
//	view: the view to be populated
//	T: a struct whose type is specified by the Record interface
//     This will be one of the main objects of the simulation, viz:
//     Commodity | Simulation | Class
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

			f0.Set(vm)
			f1.Set(cm)
		default:
		}
	}
}

// Create a CommodityView for display in a template
// taking data from two Commodity objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func VeryNewCommodityView(v *Commodity, c *Commodity) *CommodityView {
	recordBase := RecordBase[Commodity]{
		Viewed:   v,
		Compared: c,
	}
	view := CommodityView{RecordBase: recordBase}
	PopulateView(&view)
	return &view
}

// Create a slice of CommodityView objects for display in a template
// taking data from two Commodity objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func VeryNewCommodityViews(v *[]Commodity, c *[]Commodity) *[]CommodityView {
	var newViews = make([]CommodityView, len(*v))
	for i := range *v {
		newView := VeryNewCommodityView(&(*v)[i], &(*c)[i])
		newViews[i] = *newView
	}
	return &newViews
}

package models

import (
	"encoding/json"
	"fmt"
	"gorilla-client/utils"
	"reflect"
)

// List of the user's Simulations.
//
//	u: the user
//	returns:
//	 Slice of SimulationsList
//	 If the user has no simulations, an empty slice
func (u User) SimulationsList() *[]Simulation {
	list := u.Simulations.Table.(*[]Simulation)
	if len(*list) == 0 {
		var fakeList []Simulation = *new([]Simulation)
		return &fakeList
	}
	return list
}

// supplies outputData to be passed into Templates for display
//
//		u: a user
//
//		returns:
//	     if the user has no simulations, just the template list
//	     otherwise, the output data the users current simulation
func (u *User) TemplateData(message string) OutputData {
	slist := u.SimulationsList()
	state := u.GetCurrentState()
	utils.TraceInfof(utils.BrightYellow, "Entering TemplateData for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)
	if u.CurrentSimulationID == 0 {
		utils.TraceInfo(utils.BrightYellow, "User has no simulations")
		return OutputData{
			Title:          "Hello",
			Simulations:    nil,
			Templates:      &TemplateList,
			Count:          0,
			Username:       u.UserName,
			State:          state,
			CommodityViews: nil,
			IndustryViews:  nil,
			ClassViews:     nil,
			IndustryStocks: nil,
			ClassStocks:    nil,
			Trace:          nil,
			Message:        message,
		}
	}
	utils.TraceInfof(utils.BrightYellow, "TemplateData is retrieving data for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)
	commodityView := u.CommodityViews()
	commodityViewAsString, _ := json.MarshalIndent(commodityView, " ", " ")
	utils.TraceLogf(utils.BrightYellow, "CommodityViews returned %s", string(commodityViewAsString))
	return OutputData{
		Title:          "Hello",
		Simulations:    slist,
		Templates:      &TemplateList,
		Count:          len(*slist),
		Username:       u.UserName,
		State:          state,
		CommodityViews: u.CommodityViews(),
		IndustryViews:  u.IndustryViews(),
		ClassViews:     u.ClassViews(),
		IndustryStocks: u.IndustryStocks(*u.GetViewedTimeStamp()),
		ClassStocks:    u.ClassStocks(*u.GetViewedTimeStamp()),
		Trace:          u.Traces(*u.GetViewedTimeStamp()),
		Message:        message,
	}
}
func (u User) OutputCommodityData(message string, id int) CommodityData {
	return CommodityData{
		u.TemplateData(message),
		*u.Commodity(id),
	}
}

// Get a ClassData to display a single social class in the class.html template
//
//	u: the user
//	message: any message
//	id: the id of the social class item to be displayed
//
//	returns: classData which references this class, and embeds an OutputData
func (u User) OutputClassData(message string, id int) ClassData {
	return ClassData{
		u.TemplateData(message),
		*u.Class(id),
	}
}

// Get an IndustryData to display a single industry in the industry.html template
//
//	u: the user
//	message: any message
//	id: the id of the industry item to be displayed
//
//	returns: industryData which references this industry, and embeds an OutputData
func (u User) OutputIndustryData(message string, id int) IndustryData {
	return IndustryData{
		u.TemplateData(message),
		*u.Industry(id),
	}
}

func (u *User) LogTemplateData() string {
	output := u.TemplateData("hello")
	outputAsString, _ := json.MarshalIndent(output, " ", " ")
	return string(outputAsString)
}

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
				vmbn := viewedRecord.MethodByName(name)
				if vmbn.IsValid() {
					fmt.Printf(utils.Green+"The Field called %s is a function\n", name)
					cmbn := viewedRecord.MethodByName(name) // if vmbn is valid, we can safely assume cmbn is too
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

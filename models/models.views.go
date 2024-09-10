package models

import (
	"fmt"
	"reflect"
)

type Recorder interface {
	Commodity | Industry | Class | IndustryStock | ClassStock
}

type RecordBase[T Recorder] struct {
	Viewed   *T
	Compared *T
}

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

func (v View) ShowPlaceHolder() string {
	fmt.Println("Enter Show PlaceHolder")
	return "Placeholder"
}

type NewCommodityView struct {
	viewedRecord   *Commodity
	comparedRecord *Commodity
}

// Provides the value of the field f in the viewedRecord of a CommodityView
//
//	 f: the name of a field
//	 c: a CommodityView
//	 returns: the float32 value of the field
//	TODO return an int for an int field
func (c *NewCommodityView) ViewedField(f string) string {
	return reflect.Indirect(reflect.ValueOf(c.viewedRecord)).FieldByName(f).String()
}

// Provides the value of the field f in the comparedRecord of a CommodityView
//
//	 f: the name of a field
//	 c: a CommodityView
//	 returns: the float32 value of the field
//	TODO return an int for an int field
func (c *NewCommodityView) ComparedField(f string) string {
	return reflect.Indirect(reflect.ValueOf(c.comparedRecord)).FieldByName(f).String()
}

func CreateCommodityView(v *Commodity, c *Commodity) View {
	return View{&NewCommodityView{
		viewedRecord:   v,
		comparedRecord: c,
	}}
}

// Create a slice of CommodityView for display in a template
// taking data from two Commodity objects; one being viewed now,
// the other showing the state of the simulation at some time in the 'past'
func NewCommodityViews(v *[]Commodity, c *[]Commodity) *[]View {
	var newViews = make([]View, len(*v))
	var vc *Commodity
	var cc *Commodity
	fmt.Println("Entered New CommodityViews")
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateCommodityView(vc, cc)
		newViews[i] = newView
		fmt.Println(newView.Show("Size"))
		fmt.Println(newView.ShowPlaceHolder())
	}
	return &newViews
}

// Provide a string, suitable for display in a template, that formats
// a viewed value and highlights values that have changed.
//
//	v: a View object
//	f: the name of the field to display
//	Returns: a safe HTML formatted string
func (v *View) Show(f string) string {
	vv := v.Viewer.ViewedField(f)
	vc := v.Viewer.ComparedField(f)
	if vv == vc {
		return "same"
	}
	return "different"
}

type NewIndustryView struct {
	viewedRecord   *Industry
	comparedRecord *Industry
}

func (i *NewIndustryView) ViewedField(f string) float32 {
	return float32(reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f).Float())
}

func (i *NewIndustryView) ComparedField(f string) float32 {
	return float32(reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f).Float())
}

type CommodityViewer struct {
	RecordBase[Commodity]
	Id                        int
	Name                      string
	Origin                    string
	Usage                     string
	Size                      Pair
	TotalValue                Pair
	TotalPrice                Pair
	UnitValue                 Pair
	UnitPrice                 Pair
	TurnoverTime              Pair
	Demand                    Pair
	Supply                    Pair
	AllocationRatio           Pair
	MonetarilyEffectiveDemand float32
	InvestmentProportion      float32
}

type IndustryViewer struct {
	RecordBase[Industry]
	Id                   int
	Name                 string
	OutputCommodityId    int
	Output               string
	OutputScale          Pair
	OutputGrowthRate     Pair
	InitialCapital       Pair
	WorkInProgress       Pair
	CurrentCapital       Pair
	ConstantCapitalSize  Pair
	ConstantCapitalValue Pair
	ConstantCapitalPrice Pair
	VariableCapitalSize  Pair
	VariableCapitalValue Pair
	VariableCapitalPrice Pair
	MoneyStockSize       Pair
	MoneyStockValue      Pair
	MoneyStockPrice      Pair
	SalesStockSize       Pair
	SalesStockValue      Pair
	SalesStockPrice      Pair
	Profit               Pair
	ProfitRate           Pair
}

type ClassViewer struct {
	RecordBase[Class]
	Id                    int
	Name                  string
	SimulationId          int32
	TimeStamp             int
	UserName              string
	Population            Pair
	ParticipationRatio    float32
	ConsumptionRatio      float32
	Revenue               Pair
	Assets                Pair
	ConsumptionStockSize  Pair
	ConsumptionStockValue Pair
	ConsumptionStockPrice Pair
	MoneyStockSize        Pair
	MoneyStockValue       Pair
	MoneyStockPrice       Pair
	SalesStockSize        Pair
	SalesStockValue       Pair
	SalesStockPrice       Pair
}

type IndustryStockViewer struct {
	RecordBase[IndustryStock]
	Id           int
	SimulationId int
	IndustryId   int
	CommodityId  int
	UserName     string
	Name         string
	UsageType    string
	Size         Pair
	Value        Pair
	Price        Pair
	Requirement  Pair
	Demand       Pair
}

type ClassStockViewer struct {
	RecordBase[ClassStock]
	Id           int
	SimulationId int
	ClassId      int
	CommodityId  int
	UserName     string
	Name         string
	UsageType    string
	Size         Pair
	Value        Pair
	Price        Pair
	Requirement  Pair
	Demand       Pair
}

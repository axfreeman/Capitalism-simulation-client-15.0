package models

import "reflect"

type Recorder interface {
	Commodity | Industry | Class | IndustryStock | ClassStock
}

type RecordBase[T Recorder] struct {
	Viewed   *T
	Compared *T
}

type Viewer interface {
	ViewedField(f string) float32
	ComparedField(f string) float32
}

type View struct {
	Viewer
}

func (v *View) Show(f string) string {
	vv := v.Viewer.ViewedField(f)
	vc := v.Viewer.ComparedField(f)
	if vv == vc {
		return "same"
	}
	return "different"
}

type NewCommodityView struct {
	viewedRecord   *Commodity
	comparedRecord *Commodity
}

func (c *NewCommodityView) ViewedField(f string) float32 {
	return float32(reflect.Indirect(reflect.ValueOf(c.viewedRecord)).FieldByName(f).Float())
}

func (c *NewCommodityView) ComparedField(f string) float32 {
	return float32(reflect.Indirect(reflect.ValueOf(c.comparedRecord)).FieldByName(f).Float())
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

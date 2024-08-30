package models

type Recorder interface {
	Commodity | Industry | Class
}

type RecordBase[T Recorder] struct {
	Viewed   *T
	Compared *T
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

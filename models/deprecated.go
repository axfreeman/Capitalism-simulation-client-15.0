package models

func (u User) OldCommodities() *[]Commodity {
	return (*u.TableSets[*u.GetViewedTimeStamp()])["commodities"].Table.(*[]Commodity)
}

func (u User) OldIndustries() *[]Industry {
	return (*u.TableSets[*u.GetViewedTimeStamp()])["industries"].Table.(*[]Industry)
}

func (u User) OldClasses() *[]Class {
	return (*u.TableSets[*u.GetViewedTimeStamp()])["classes"].Table.(*[]Class)
}

// Wrapper for the IndustryStockList
func (u User) OldIndustryStocks(timeStamp int) *[]IndustryStock {
	return (*u.TableSets[timeStamp])["industry stocks"].Table.(*[]IndustryStock)
}

// Wrapper for the ClassStockList
func (u User) OldClassStocks(timeStamp int) *[]ClassStock {
	return (*u.TableSets[timeStamp])["class stocks"].Table.(*[]ClassStock)
}

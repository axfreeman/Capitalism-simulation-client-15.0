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

// Find the commodity with a given id.
//
//	u: the user to whom the commodity belongs
//	Return: pointer to the commodity if it found
//	Return: pointer to NotFoundCommodity if not found.
func (u User) OldCommodity(id int) *Commodity {
	// commodityList := *LoggedInUsers[u.UserName].OldCommodities()
	commodityList := *ViewedObjects[Commodity](u, `commodities`)
	for i := 0; i < len(commodityList); i++ {
		c := commodityList[i]
		if id == c.Id {
			return &c
		}
	}
	return &NotFoundCommodity
}

// Find the class with a given id.
//
//	u: the user to whom the class belongs
//	Return: pointer to the class if it found
//	Return: pointer to NotFoundClass if not found.
func (u User) OldClass(id int) *Class {
	// classList := *LoggedInUsers[u.UserName].OldClasses()
	classList := *ViewedObjects[Class](u, `classes`)
	for i := 0; i < len(classList); i++ {
		c := classList[i]
		if id == c.Id {
			return &c
		}
	}
	return &NotFoundClass
}

// Find the industry with a given id.
//
//	u: the user to whom the industry belongs
//	Return: pointer to the industry if it found
//	Return: pointer to NotFoundIndustry if not found.
func (u User) OldIndustry(id int) *Industry {
	// industryList := *LoggedInUsers[u.UserName].OldIndustries()
	industryList := *ViewedObjects[Industry](u, `industries`)
	for i := 0; i < len(industryList); i++ {
		ind := industryList[i]
		if id == ind.Id {
			return &ind
		}
	}
	return &NotFoundIndustry
}

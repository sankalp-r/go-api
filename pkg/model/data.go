package model

// DataContainer for stroring collection of Data
type DataContainer struct {
	Data []Data `json:"data"`
}

// Data for storing the data entity
type Data struct {
	URL            string  `json:"url"`
	Views          int     `json:"views"`
	RelevanceScore float32 `json:"relevanceScore"`
}

// DataByView to sort by views
type DataByView []Data

func (dv DataByView) Len() int {
	return len(dv)
}

func (dv DataByView) Swap(i, j int) {
	dv[i], dv[j] = dv[j], dv[i]
}

func (dv DataByView) Less(i, j int) bool {
	return dv[i].Views < dv[j].Views
}

// DataByRelevanceScore to sort by relevance score
type DataByRelevanceScore []Data

func (dr DataByRelevanceScore) Len() int {
	return len(dr)
}

func (dr DataByRelevanceScore) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

func (dr DataByRelevanceScore) Less(i, j int) bool {
	return dr[i].RelevanceScore < dr[j].RelevanceScore
}

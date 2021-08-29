package model

type CensusVariablesGroups struct {
	Name        string `db:"groupid"`
	Description string `db:"description"`
	Variables   string `db:"variableslink"`
}

type CensusVariables struct {
	VariableID string `json:"variableid"`
	Label      string `json:"label"`
	Concept    string `json:"concept"`
	GroupID    string `json:"groupid"`
}

type CBSAInfo struct {
	CbsaCode  string       `json:"cbsafipscode"`
	CbsaTitle string       `json:"cbsaname"`
	Counties  []CountyInfo `json:"counties"`
}

type StateInfo struct {
	FipsStateCode     string `json:"statefipscode"`
	StateName         string `json:"statename"`
	StateAbbreviation string `json:"stateabbreviation"`
}

type CountyInfo struct {
	CountyFullCode string `json:"countyfullcode"`
	FipsCountyCode string `json:"countyfipscode"`
	CountyName     string `json:"countyname"`
	StateInfo      StateInfo
}

type TractInfo struct {
	TractCode      string `json:"tractcode"`
	CensusYear     int    `json:"censusyear"`
	CountyFullCode string `json:"countyfullcode"`
	FipsStateCode  string `json:"statefipscode"`
	FipsCountyCode string `json:"countyfipscode"`
}

type EsriTractsInfo struct {
	TractCode               string                  `json:"tractcode"`
	CountyFullCode          string                  `json:"countyfullcode"`
	FipsStateCode           string                  `json:"statefipscode"`
	EsriStandardGeoFeatures EsriStandardGeoFeatures `json:"standardgeofeatures"`
}

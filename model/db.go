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

type StateInfo struct {
	FipsStateCode string `json:"statefipscode"`
	StateName     string `json:"statename"`
}

type CountyInfo struct {
	FipsCountyCode         string `json:"countyfipscode"`
	CountyCountyEquivalent string `json:"countyname"`
	StateInfo              StateInfo
}

type CBSAInfo struct {
	CbsaCode  string       `json:"cbsafipscode"`
	CbsaTitle string       `json:"cbsaname"`
	Counties  []CountyInfo `json:"counties"`
}

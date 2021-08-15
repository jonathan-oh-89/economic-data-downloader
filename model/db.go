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
	fipsstatecode string `json:"statefipscode"`
	statename     string `json:"statename"`
}

type CountyInfo struct {
	fipscountycode         string `json:"countyfipscode"`
	countycountyequivalent string `json:"countyname"`
	StateInfo              StateInfo
}

type CBSAInfo struct {
	cbsacode  string       `json:"cbsacode"`
	cbsatitle string       `json:"cbsaname"`
	counties  []CountyInfo `json:"counties"`
}

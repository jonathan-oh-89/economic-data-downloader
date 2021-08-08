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

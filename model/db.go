package model

type CensusVariablesGroups struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	Variables   string `db:"variableslink"`
}

// type CensusVariables struct {
// 	Name        string `db:"name"`
// 	Label       string `db:"name"`
// 	Concept string `db:"concept"`
// }

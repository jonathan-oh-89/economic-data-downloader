package model

type CensusVariablesResponse struct {
	Label         string `json:"label"`
	Concept       string `json:"concept"`
	PredicateType string `json:"predicateType"`
	Group         string `json:"group"`
	Limit         int    `json:"limit"`
	PredicateOnly bool   `json:"predicateOnly"`
}

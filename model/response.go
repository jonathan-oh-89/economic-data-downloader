package model

type CensusVariablesResponse struct {
	Label         string `json:"label"`
	Concept       string `json:"concept"`
	PredicateType string `json:"predicateType"`
	Group         string `json:"group"`
	Limit         int    `json:"limit"`
	PredicateOnly bool   `json:"predicateOnly"`
}

type EsriStandardGeoResponse struct {
	Results []EsriStandardGeo `json:"results"`
}

type EsriStandardGeo struct {
	ParamName interface{}          `json:"paramname"`
	DataType  interface{}          `json:"dataType"`
	Value     EsriStandardGeoValue `json:"value"`
}

type EsriStandardGeoValue struct {
	Features []EsriStandardGeoFeatures
}

type EsriStandardGeoFeatures struct {
	Attributes EsriStandardGeoFeatureAttributes `json:"attributes"`
	Geometry   EsriStandardGeoFeatureGeometry   `json:"geometry"`
}

type EsriStandardGeoFeatureAttributes struct {
	DatasetID            string `json:"datasetid"`
	DataLayerID          string `json:"datalayerid"`
	AreaID               string `json:"areaid"`
	AreaName             string `json:"areaname"`
	MajorSubdivisionName string `json:"majorsubdivisionname"`
	MajorSubdivisionAbbr string `json:"majorsubdivisionabbr"`
	MajorSubdivisionType string `json:"majorsubdivisiontype"`
	CountryAbbr          string `json:"countryabbr"`
	ObjectId             int    `json:"objectid"`
	Score                int    `json:"score"`
}

type EsriStandardGeoFeatureGeometry struct {
	Rings [][][]float64
}

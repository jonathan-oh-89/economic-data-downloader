package model

type CensusVariablesResponse struct {
	Label         string `json:"label"`
	Concept       string `json:"concept"`
	PredicateType string `json:"predicateType"`
	Group         string `json:"group"`
	Limit         int    `json:"limit"`
	PredicateOnly bool   `json:"predicateOnly"`
}

// ESRI STANDARD GEO

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

// ESRI ENRICH

type EsriEnrichResponse struct {
	Results  []EsriEnrichResults `json:"results"`
	Messages []string            `json:"messages"`
}

type EsriEnrichResults struct {
	ParamName interface{}     `json:"paramname"`
	DataType  interface{}     `json:"dataType"`
	Value     EsriEnrichValue `json:"value"`
}

type EsriEnrichValue struct {
	FeatureSet []EsriEnrichFeatureSet
}

type EsriEnrichFeatureSet struct {
	Features []EsriEnrichFeatures
}

type EsriEnrichFeatures struct {
	Attributes EsriEnrichFeaturesAttributes `json:"attributes"`
}

type EsriEnrichFeaturesAttributes struct {
	StdGeographyLevel string `json:"stdgeographylevel"`
	StdGeographyName  string `json:"stdgeographyname"`
	StdGeographyID    string `json:"stdgeographyid"`
	CRMCYPERC         int    `json:"crmcyperc"`
	CRMCYMURD         int    `json:"crmcymurd"`
	CRMCYRAPE         int    `json:"crmcyrape"`
	CRMCYROBB         int    `json:"crmcyrobb"`
	CRMCYASST         int    `json:"crmcyasst"`
	CRMCYPROC         int    `json:"crmcyproc"`
	CRMCYBURG         int    `json:"crmcyburg"`
	CRMCYLARC         int    `json:"crmcylarc"`
	CRMCYMVEH         int    `json:"crmcymveh"`
}

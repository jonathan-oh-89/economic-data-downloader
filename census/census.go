package census

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/jonathan-oh-89/economic-data-downloader/utils"
)

type variablesDescGroup struct {
	variableDesc string
	category     string
}

var censusGroups = []string{"B01001", "B03002", "B19013", "B08301", "B08303", "B11012"}

func Test(censusGroupID string, geoLevel string) {

	variablesDescGroup := getVariablesToInclude(censusGroupID, geoLevel, true)

	apiParams := ""
	for key, _ := range variablesDescGroup {
		apiParams = apiParams + key + ","
	}

	testAPI(apiParams[:len(apiParams)-1], censusGroupID, geoLevel)
}

func testAPI(apiparams string, groupid string, geoLevel string) {

	url := fmt.Sprintf("https://api.census.gov/data/2019/acs/acs5?get=NAME,%s&for=county:059,037&in=state:06", apiparams)

	response, err := http.Get(url)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var arr [][]string
	err = json.Unmarshal(responseData, &arr)
	if err != nil {
		log.Fatal(err)
	}

	censusVariablesLookup2 := getVariablesToInclude(groupid, geoLevel, false)

	for i, dataRow := range arr[1:] {
		testAggregate(arr[0], dataRow, censusVariablesLookup2, true)

		if i > 1 {
			log.Print("*** TEMP STOP ITERATE")
			break
		}
	}
}

// Pass in geoFips and geoLevel. Iterate through census groups and gather all the variables/data needed.
func Do(geoFips string, geoLevel string) {

	for _, censusGroupID := range censusGroups {
		variablesDescGroup := getVariablesToInclude(censusGroupID, geoLevel, true)

		apiParams := ""
		for key, _ := range variablesDescGroup {
			apiParams = apiParams + key + ","
		}

		callAPI(apiParams[:len(apiParams)-1], censusGroupID, geoLevel, geoFips)
	}
}

func callAPI(apiparams string, groupid string, geoLevel string, geoFips string) {
	url := ""
	if geoLevel == "county" {
		url = fmt.Sprintf("https://api.census.gov/data/2019/acs/acs5?get=NAME,%s&for=county:059,037&in=state:%s", apiparams, geoFips)
	}

	response, err := http.Get(url)
	if err != nil {
		fmt.Print("ERROR: ", err)
		panic(err)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var arr [][]string
	err = json.Unmarshal(responseData, &arr)
	if err != nil {
		fmt.Print(err)
	}

	censusVariablesLookup2 := getVariablesToInclude(groupid, geoLevel, false)

	for _, dataRow := range arr[1:] {
		_ = dataRow
		aggregate(arr[0], dataRow, censusVariablesLookup2, true)
	}
}

// Returns a map of census variable groups. VariableId is used as key. Define all the data to be stored here.
func getVariablesToInclude(groupID string, geoLevel string, paramsOnly bool) map[string]variablesDescGroup {

	includeVariables := map[string]variablesDescGroup{}

	switch groupID {
	case "B01001":
		includeVariables = map[string]variablesDescGroup{
			"B01001_001E": {variableDesc: "Total", category: "Total"},
			"B01001_003E": {variableDesc: "Male!!Under 5 years", category: "under18"},
			"B01001_027E": {variableDesc: "Female!!Under 5 years", category: "under18"},
			"B01001_004E": {variableDesc: "Male!!5 to 9 years", category: "under18"},
			"B01001_028E": {variableDesc: "Female!!5 to 9 years", category: "under18"},
			"B01001_029E": {variableDesc: "Female!!10 to 14 years", category: "under18"},
			"B01001_005E": {variableDesc: "Male!!10 to 14 years", category: "under18"},
			"B01001_030E": {variableDesc: "Female!!15 to 17 years", category: "under18"},
			"B01001_006E": {variableDesc: "Male!!15 to 17 years", category: "under18"},
			"B01001_007E": {variableDesc: "Male!!18 and 19 years", category: "age18to24"},
			"B01001_031E": {variableDesc: "Female!!18 and 19 years", category: "age18to24"},
			"B01001_032E": {variableDesc: "Female!!20 years", category: "age18to24"},
			"B01001_008E": {variableDesc: "Male!!20 years", category: "age18to24"},
			"B01001_033E": {variableDesc: "Female!!21 years", category: "age18to24"},
			"B01001_009E": {variableDesc: "Male!!21 years", category: "age18to24"},
			"B01001_010E": {variableDesc: "Male!!22 to 24 years", category: "age18to24"},
			"B01001_034E": {variableDesc: "Female!!22 to 24 years", category: "age18to24"},
			"B01001_035E": {variableDesc: "Female!!25 to 29 years", category: "age25to39"},
			"B01001_011E": {variableDesc: "Male!!25 to 29 years", category: "age25to39"},
			"B01001_012E": {variableDesc: "Male!!30 to 34 years", category: "age25to39"},
			"B01001_036E": {variableDesc: "Female!!30 to 34 years", category: "age25to39"},
			"B01001_037E": {variableDesc: "Female!!35 to 39 years", category: "age25to39"},
			"B01001_013E": {variableDesc: "Male!!35 to 39 years", category: "age25to39"},
			"B01001_014E": {variableDesc: "Male!!40 to 44 years", category: "age40to64"},
			"B01001_038E": {variableDesc: "Female!!40 to 44 years", category: "age40to64"},
			"B01001_039E": {variableDesc: "Female!!45 to 49 years", category: "age40to64"},
			"B01001_015E": {variableDesc: "Male!!45 to 49 years", category: "age40to64"},
			"B01001_040E": {variableDesc: "Female!!50 to 54 years", category: "age40to64"},
			"B01001_016E": {variableDesc: "Male!!50 to 54 years", category: "age40to64"},
			"B01001_041E": {variableDesc: "Female!!55 to 59 years", category: "age40to64"},
			"B01001_017E": {variableDesc: "Male!!55 to 59 years", category: "age40to64"},
			"B01001_042E": {variableDesc: "Female!!60 and 61 years", category: "age40to64"},
			"B01001_018E": {variableDesc: "Male!!60 and 61 years", category: "age40to64"},
			"B01001_043E": {variableDesc: "Female!!62 to 64 years", category: "age40to64"},
			"B01001_019E": {variableDesc: "Male!!62 to 64 years", category: "age40to64"},
			"B01001_044E": {variableDesc: "Female!!65 and 66 years", category: "age65up"},
			"B01001_020E": {variableDesc: "Male!!65 and 66 years", category: "age65up"},
			"B01001_045E": {variableDesc: "Female!!67 to 69 years", category: "age65up"},
			"B01001_021E": {variableDesc: "Male!!67 to 69 years", category: "age65up"},
			"B01001_046E": {variableDesc: "Female!!70 to 74 years", category: "age65up"},
			"B01001_022E": {variableDesc: "Male!!70 to 74 years", category: "age65up"},
			"B01001_023E": {variableDesc: "Male!!75 to 79 years", category: "age65up"},
			"B01001_047E": {variableDesc: "Female!!75 to 79 years", category: "age65up"},
			"B01001_048E": {variableDesc: "Female!!80 to 84 years", category: "age65up"},
			"B01001_024E": {variableDesc: "Male!!80 to 84 years", category: "age65up"},
			"B01001_025E": {variableDesc: "Male!!85 years and over", category: "age65up"},
			"B01001_049E": {variableDesc: "Female!!85 years and over", category: "age65up"},
		}
	case "B03002":
		includeVariables = map[string]variablesDescGroup{
			"B03002_001E": {variableDesc: "Total", category: "Total"},
			"B03002_012E": {variableDesc: "Hispanic or Latino", category: "hispanicLatino"},
			"B03002_005E": {variableDesc: "Not Hispanic or Latino!!American Indian and Alaska Native alone", category: "americanIndianAlaskaNative"},
			"B03002_006E": {variableDesc: "Not Hispanic or Latino!!Asian alone", category: "asian"},
			"B03002_004E": {variableDesc: "Not Hispanic or Latino!!Black or African American alone", category: "black"},
			"B03002_007E": {variableDesc: "Not Hispanic or Latino!!Native Hawaiian and Other Pacific Islander alone", category: "nativeHawaiianPacificIslander"},
			"B03002_008E": {variableDesc: "Not Hispanic or Latino!!Some other race alone", category: "biracialandOther"},
			"B03002_009E": {variableDesc: "Not Hispanic or Latino!!Two or more races", category: "biracialandOther"},
			"B03002_003E": {variableDesc: "Not Hispanic or Latino!!White alone", category: "white"},
		}
	case "B19013":
		includeVariables = map[string]variablesDescGroup{
			"B19013_001E": {variableDesc: "Median household income in the past 12 months (in 2019 inflation-adjusted dollars)", category: "medianHouseholdIncome"},
		}
	case "B08301":
		includeVariables = map[string]variablesDescGroup{
			"B08301_018E": {variableDesc: "Bicycle", category: "bicycle"},
			"B08301_004E": {variableDesc: "Car, truck, or van!!Carpooled", category: "other"},
			"B08301_003E": {variableDesc: "Car, truck, or van!!Drove alone", category: "droveAlone"},
			"B08301_017E": {variableDesc: "Motorcycle", category: "droveAlone"},
			"B08301_020E": {variableDesc: "Other means", category: "other"},
			"B08301_010E": {variableDesc: "Public transportation (excluding taxicab)", category: "publicTransportation"},
			"B08301_019E": {variableDesc: "Walked", category: "walked"},
			"B08301_021E": {variableDesc: "Worked from home", category: "workFromHome"},
			"B08301_016E": {variableDesc: "Taxicab", category: "taxi"},
			"B08301_001E": {variableDesc: "Total", category: ""},
		}
	case "B08303":
		includeVariables = map[string]variablesDescGroup{
			"B08303_001E": {variableDesc: "Total", category: ""},
			"B08303_004E": {variableDesc: "10 to 14 minutes", category: "lessThan15"},
			"B08303_005E": {variableDesc: "15 to 19 minutes", category: "travelTime15_30"},
			"B08303_006E": {variableDesc: "20 to 24 minutes", category: "travelTime15_30"},
			"B08303_007E": {variableDesc: "25 to 29 minutes", category: "travelTime15_30"},
			"B08303_008E": {variableDesc: "30 to 34 minutes", category: "timeTravel30_45"},
			"B08303_009E": {variableDesc: "35 to 39 minutes", category: "timeTravel30_45"},
			"B08303_010E": {variableDesc: "40 to 44 minutes", category: "timeTravel30_45"},
			"B08303_011E": {variableDesc: "45 to 59 minutes", category: "timeTravel45_1hr"},
			"B08303_003E": {variableDesc: "5 to 9 minutes", category: "lessThan15"},
			"B08303_012E": {variableDesc: "60 to 89 minutes", category: "timeTravel1hr_more"},
			"B08303_013E": {variableDesc: "90 or more minutes", category: "timeTravel1hr_more"},
			"B08303_002E": {variableDesc: "Less than 5 minutes", category: "lessThan15"},
		}
	case "B11012":
		includeVariables = map[string]variablesDescGroup{
			"B11012_007E": {variableDesc: "Cohabiting couple household!!With no own children of the householder under 18 years", category: "couplesWithNoChildren"},
			"B11012_006E": {variableDesc: "Cohabiting couple household!!With own children of the householder under 18 years", category: "nonmarriedCoupleWithChildren"},
			"B11012_009E": {variableDesc: "Female householder, no spouse or partner present!!Living alone", category: "onePerson"},
			"B11012_012E": {variableDesc: "Female householder, no spouse or partner present!!With only nonrelatives present", category: "other"},
			"B11012_010E": {variableDesc: "Female householder, no spouse or partner present!!With own children under 18 years", category: "singleParentWithChildren"},
			"B11012_011E": {variableDesc: "Female householder, no spouse or partner present!!With relatives, no own children under 18 years", category: "other"},
			"B11012_014E": {variableDesc: "Male householder, no spouse or partner present!!Living alone", category: "onePerson"},
			"B11012_017E": {variableDesc: "Male householder, no spouse or partner present!!With only nonrelatives present", category: "other"},
			"B11012_015E": {variableDesc: "Male householder, no spouse or partner present!!With own children under 18 years", category: "singleParentWithChildren"},
			"B11012_016E": {variableDesc: "Male householder, no spouse or partner present!!With relatives, no own children under 18 years", category: "other"},
			"B11012_004E": {variableDesc: "Married-couple household!!With no own children under 18 years", category: "couplesWithNoChildren"},
			"B11012_003E": {variableDesc: "Married-couple household!!With own children under 18 years", category: "marriedWithChildren"},
			"B11012_001E": {variableDesc: "Total", category: ""},
		}
	}

	if paramsOnly {
		return includeVariables
	}

	includeVariables["NAME"] = variablesDescGroup{variableDesc: "GeoName", category: "GeoName"}

	switch geoLevel {
	case "state":
		includeVariables["state"] = variablesDescGroup{variableDesc: "state", category: "state"}
	case "county":
		includeVariables["state"] = variablesDescGroup{variableDesc: "state", category: "state"}
		includeVariables["county"] = variablesDescGroup{variableDesc: "county", category: "county"}
	}
	return includeVariables
}

// Gets a sample set for a census variable groupid and stores into csv.
func DownloadToCSV(testgroup string) {
	censusVariablesLookup := map[string]string{}

	lines := utils.ReadCSV("/files/censusvariables.csv")

	groupid := ""
	apiParams := ""
	for i, line := range lines {
		if i < 2 {
			if i < 1 {
				continue
			} else {
				groupid = line[3]
			}
		}

		apiParams += line[0] + ","
		censusVariablesLookup[line[0]] = line[1]

		if len(lines) == i+1 {
			break
		}

		if lines[i+1][3] != groupid {
			if groupid == testgroup {

				url := fmt.Sprintf("https://api.census.gov/data/2019/acs/acs5?get=NAME,%s&for=county:059&in=state:06", apiParams[:len(apiParams)-1])

				response, err := http.Get(url)
				if err != nil {
					fmt.Print("ERROR: ", err)
					panic(err)
				}

				defer response.Body.Close()

				responseData, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatal(err)
				}

				var arr [][]string
				err = json.Unmarshal(responseData, &arr)
				if err != nil {
					fmt.Print(err)
				}

				variableNameArr := []string{}

				for _, variableID := range arr[0] {
					colname := censusVariablesLookup[variableID]

					if colname == "" {
						switch variableID {
						case "NAME":
							colname = "NAME"
						case "state":
							colname = "state"
						case "county":
							colname = "county"
						default:
							panic("FOUND MISSING VARIABLE NAME")
						}
					}
					variableNameArr = append(variableNameArr, colname)
				}

				df := dataframe.New(
					series.New(variableNameArr, series.String, "variableNameArr"),
					series.New(arr[0], series.String, "COL.1"),
					series.New(arr[1], series.String, "COL.2"),
				)

				utils.SaveToCSV(df)
			}
			groupid = lines[i+1][3]
			apiParams = ""
			censusVariablesLookup = map[string]string{}
		}
	}

}

func getPopulation() {
	//https://api.census.gov/data/2019/pep/population?get=COUNTY,DATE_CODE,DATE_DESC,DENSITY,POP,NAME,STATE&for=state:01

}

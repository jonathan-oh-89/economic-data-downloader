package census

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jonathan-oh-89/economic-data-downloader/utils"
)

func AggregateRace(arr [][]string, censusvariableslookup map[string]string, groupid string) {
	variableDescLookup := variableInclusion(censusvariableslookup, groupid)

	for k, _ := range variableDescLookup {
		variableDescLookup[k] = strings.Replace(variableDescLookup[k], "Not Hispanic or Latino!!", "", -1)
		variableDescLookup[k] = strings.Replace(variableDescLookup[k], " alone", "", -1)
	}

	for _, dataRow := range arr[1:] {
		raceAggregate := mapVariableDescToValue(arr[0], dataRow, variableDescLookup)

		for k, v := range raceAggregate {

			if k == "Two or more races" {
				raceAggregate["Other/Multiracial"] += v
				delete(raceAggregate, "Two or more races")
			} else if k == "Some other race" {
				raceAggregate["Other/Multiracial"] += v
				delete(raceAggregate, "Some other race")
			}
		}

		newAggregate := calculatePercentage(raceAggregate)

		raceAggregate = map[string]int{}

		fmt.Print(newAggregate)
		// Do something with mapFinalAggregate
	}
}

func AggregateAge(arr [][]string, censusvariableslookup map[string]string, groupid string) {
	variableDescLookup := variableInclusion(censusvariableslookup, groupid)

	for variableID, variableName := range variableDescLookup {
		variableName = strings.Replace(variableName, "Male!!", "", -1)
		variableName = strings.Replace(variableName, "Female!!", "", -1)
		variableDescLookup[variableID] = variableName
	}

	finalAgeAggregate := map[string]int{}

	for _, dataRow := range arr[1:] {
		ageAggregate := mapVariableDescToValue(arr[0], dataRow, variableDescLookup)

		for k, v := range ageAggregate {

			under18 := []string{"Under 5 years", "5 to 9 years", "10 to 14 years", "15 to 17 years"}
			age18to24 := []string{"18 and 19 years", "20 years", "21 years", "xxxxxx", "22 to 24 years"}
			age25to39 := []string{"25 to 29 years", "30 to 34 years", "35 to 39 years"}
			age40to64 := []string{"40 to 44 years", "45 to 49 years", "50 to 54 years", "55 to 59 years", "60 and 61 years", "62 to 64 years"}
			age65up := []string{"65 and 66 years", "67 to 69 years", "70 to 74 years", "75 to 79 years", "80 to 84 years", "85 years and over"}
			total := []string{"Total"}

			if utils.CheckStringInList(under18, k) {
				finalAgeAggregate["under18"] += v
			} else if utils.CheckStringInList(age18to24, k) {
				finalAgeAggregate["age18to24"] += v
			} else if utils.CheckStringInList(age25to39, k) {
				finalAgeAggregate["age25to39"] += v
			} else if utils.CheckStringInList(age40to64, k) {
				finalAgeAggregate["age40to64"] += v
			} else if utils.CheckStringInList(age65up, k) {
				finalAgeAggregate["age65up"] += v
			} else if utils.CheckStringInList(total, k) {
				finalAgeAggregate["Total"] += v
			}
		}

		newAggregate := calculatePercentage(finalAgeAggregate)
		fmt.Print(newAggregate)
		/*
			*** CONTINUE ****
			SAVE OFF TO MONGO BY GEO
		*/
		finalAgeAggregate = map[string]int{}
		ageAggregate = map[string]int{}
	}
}

func AggregateMedianIncome(arr [][]string, censusvariableslookup map[string]string, groupid string) {
	variableDescLookup := variableInclusion(censusvariableslookup, groupid)

	for variableID, variableName := range variableDescLookup {
		variableName = strings.Replace(variableName, " in the past 12 months", "", -1)
		variableName = strings.Replace(variableName, " inflation-adjusted dollars", "", -1)
		variableName = strings.Replace(variableName, "in ", "", -1)
		variableDescLookup[variableID] = variableName
	}

	for _, dataRow := range arr[1:] {
		medianIncomeAggregate := mapVariableDescToValue(arr[0], dataRow, variableDescLookup)

		fmt.Print(medianIncomeAggregate)
		/*
			*** CONTINUE ****
			SAVE OFF TO MONGO BY GEO
		*/
	}
}

func AggregateTransportationToWork(arr [][]string, censusvariableslookup map[string]string, groupid string) []map[string]float64 {
	variableDescLookup := variableInclusion(censusvariableslookup, groupid)

	for variableID, variableName := range variableDescLookup {
		variableName = strings.Replace(variableName, "Car, truck, or van!!", "", -1)
		variableDescLookup[variableID] = variableName
	}

	postAggregate := []map[string]float64{}

	for _, dataRow := range arr[1:] {
		variableAggregate := mapVariableDescToValue(arr[0], dataRow, variableDescLookup)

		for k, v := range variableAggregate {

			bicycle := []string{"Bicycle"}
			other := []string{"Carpooled", "Other means"}
			droveAlone := []string{"Drove alone", "Motorcycle"}
			taxi := []string{"Taxicab"}
			walked := []string{"Walked"}
			workFromHome := []string{"Worked from home"}
			total := []string{"Total"}

			if utils.CheckStringInList(bicycle, k) {
				variableAggregate["bike"] += v
			} else if utils.CheckStringInList(other, k) {
				variableAggregate["other"] += v
			} else if utils.CheckStringInList(droveAlone, k) {
				variableAggregate["droveAlone"] += v
			} else if utils.CheckStringInList(taxi, k) {
				variableAggregate["taxi"] += v
			} else if utils.CheckStringInList(walked, k) {
				variableAggregate["walked"] += v
			} else if utils.CheckStringInList(workFromHome, k) {
				variableAggregate["workFromHome"] += v
			} else if strings.Contains(k, "Public transportation") {
				variableAggregate["publicTransportation"] += v
			} else if utils.CheckStringInList(total, k) {
				variableAggregate["Total"] += v
			}
		}

		percentageAggregate := calculatePercentage(variableAggregate)
		postAggregate = append(postAggregate, percentageAggregate)
		variableAggregate = map[string]int{}
	}

	return postAggregate
}

func AggregateTimeTravelToWork(arr [][]string, censusvariableslookup map[string]string, groupid string) {
	variableDescLookup := variableInclusion(censusvariableslookup, groupid)

	for variableID, variableName := range variableDescLookup {
		variableDescLookup[variableID] = variableName
	}

	finalAgeAggregate := map[string]int{}

	for _, dataRow := range arr[1:] {
		ageAggregate := mapVariableDescToValue(arr[0], dataRow, variableDescLookup)

		for k, v := range ageAggregate {

			lessThan15 := []string{"Less than 5 minutes", "5 to 9 minutes", "10 to 14 minutes"}
			travelTime15_30 := []string{"15 to 29 minutes", "20 to 24 minutes", "25 to 29 minutes"}
			timeTravel30_45 := []string{"30 to 34 minutes", "35 to 39 minutes", "40 to 44 minutes"}
			timeTravel45_1hr := []string{"45 to 59 minute"}
			timeTravel1hr_more := []string{"60 to 89 minutes", "90 or more minutes"}
			total := []string{"Total"}

			if utils.CheckStringInList(lessThan15, k) {
				finalAgeAggregate["lessThan15"] += v
			} else if utils.CheckStringInList(travelTime15_30, k) {
				finalAgeAggregate["travelTime15_30"] += v
			} else if utils.CheckStringInList(timeTravel30_45, k) {
				finalAgeAggregate["timeTravel30_45"] += v
			} else if utils.CheckStringInList(timeTravel45_1hr, k) {
				finalAgeAggregate["timeTravel45_1hr"] += v
			} else if utils.CheckStringInList(timeTravel1hr_more, k) {
				finalAgeAggregate["timeTravel1hr_more"] += v
			} else if utils.CheckStringInList(total, k) {
				finalAgeAggregate["Total"] += v
			}
		}

		newAggregate := calculatePercentage(finalAgeAggregate)
		fmt.Print(newAggregate)
		/*
			*** CONTINUE ****
			SAVE OFF TO MONGO BY GEO
		*/
		finalAgeAggregate = map[string]int{}
		ageAggregate = map[string]int{}
	}
}

func variableInclusion(censusvariableslookup map[string]string, groupID string) map[string]string {

	includeVariables := []string{}

	if groupID == "B01001" {
		includeVariables = []string{
			"Male!!Under 5 years",
			"Male!!85 years and over",
			"Male!!80 to 84 years",
			"Male!!75 to 79 years",
			"Male!!70 to 74 years",
			"Male!!67 to 69 years",
			"Male!!65 and 66 years",
			"Male!!62 to 64 years",
			"Male!!60 and 61 years",
			"Male!!55 to 59 years",
			"Male!!50 to 54 years",
			"Male!!5 to 9 years",
			"Male!!45 to 49 years",
			"Male!!40 to 44 years",
			"Male!!35 to 39 years",
			"Male!!30 to 34 years",
			"Male!!25 to 29 years",
			"Male!!22 to 24 years",
			"Male!!21 years",
			"Male!!20 years",
			"Male!!18 and 19 years",
			"Male!!15 to 17 years",
			"Male!!10 to 14 years",
			"Female!!Under 5 years",
			"Female!!85 years and over",
			"Female!!80 to 84 years",
			"Female!!75 to 79 years",
			"Female!!70 to 74 years",
			"Female!!67 to 69 years",
			"Female!!65 and 66 years",
			"Female!!62 to 64 years",
			"Female!!60 and 61 years",
			"Female!!55 to 59 years",
			"Female!!50 to 54 years",
			"Female!!5 to 9 years",
			"Female!!45 to 49 years",
			"Female!!40 to 44 years",
			"Female!!35 to 39 years",
			"Female!!30 to 34 years",
			"Female!!25 to 29 years",
			"Female!!22 to 24 years",
			"Female!!21 years",
			"Female!!20 years",
			"Female!!18 and 19 years",
			"Female!!15 to 17 years",
			"Female!!10 to 14 years",
		}
	} else if groupID == "B03002" {
		includeVariables = []string{
			"Hispanic or Latino",
			"Not Hispanic or Latino!!American Indian and Alaska Native alone",
			"Not Hispanic or Latino!!Asian alone",
			"Not Hispanic or Latino!!Black or African American alone",
			"Not Hispanic or Latino!!Native Hawaiian and Other Pacific Islander alone",
			"Not Hispanic or Latino!!Some other race alone",
			"Not Hispanic or Latino!!Two or more races",
			"Not Hispanic or Latino!!White alone",
		}
	} else if groupID == "B19013" {
		includeVariables = []string{
			"Median household income in the past 12 months (in 2019 inflation-adjusted dollars)",
		}
	} else if groupID == "B08301" {
		includeVariables = []string{
			"Bicycle",
			"Car, truck, or van!!Carpooled",
			"Car, truck, or van!!Drove alone",
			"Motorcycle",
			"Other means",
			"Public transportation (excluding taxicab)",
			"Taxicab",
			"Walked",
			"Worked from home",
		}
	} else if groupID == "B08303" {
		includeVariables = []string{
			"10 to 14 minutes",
			"15 to 19 minutes",
			"20 to 24 minutes",
			"25 to 29 minutes",
			"30 to 34 minutes",
			"35 to 39 minutes",
			"40 to 44 minutes",
			"45 to 59 minutes",
			"5 to 9 minutes",
			"60 to 89 minutes",
			"90 or more minutes",
			"Less than 5 minutes",
		}
	}

	includeVariables = append(includeVariables, "Total")
	variableDescLookup := map[string]string{}

	for variableID, variableName := range censusvariableslookup {
		if utils.CheckStringInList(includeVariables, variableName) {
			variableDescLookup[variableID] = variableName
		}
	}

	return variableDescLookup
}

func mapVariableDescToValue(variableCols []string, arr []string, variableDescLookup map[string]string) map[string]int {

	mapVariableDesc := map[string]int{}

	for i, row := range arr {
		if variableDesc, ok := variableDescLookup[variableCols[i]]; ok {

			val, err := strconv.Atoi(row)
			if err != nil {
				fmt.Print("Error converting int: ", row)
			}

			mapVariableDesc[variableDesc] += int(val)
		}
	}

	return mapVariableDesc
}

func calculatePercentage(dataMap map[string]int) map[string]float64 {
	mapFinalAggregate := map[string]float64{}
	totalPop := float64(0)
	if val, ok := dataMap["Total"]; ok {
		totalPop = float64(val)
	}

	for k, v := range dataMap {
		value := float64(v)
		num := (value / totalPop)
		mapFinalAggregate[k] = float64(math.Round(num*100) / 100)
	}

	return mapFinalAggregate
}

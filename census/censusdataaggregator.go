package census

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jonathan-oh-89/economic-data-downloader/utils"
)

var geoLevels = []string{"state", "county"}

func testAggregate(header []string, dataRow []string, censusvariableslookup map[string]variablesDescGroup, calculatePercentage bool) {
	aggregate := map[string]interface{}{}
	totalSizeForGroup := 0.0

	headerDescriptions := mapHeaderToDescription(header, censusvariableslookup)

	if len(dataRow) != len(headerDescriptions) {
		panic("WE RECEIVED MORE/LESS FIELDS FROM CENSUS API")
	}

	for i, v := range dataRow {
		variabledescgroup := headerDescriptions[i]

		if utils.CheckStringInList(geoLevels, variabledescgroup.variableDesc) {
			aggregate[variabledescgroup.variableDesc] = utils.FormatLeadingZeroes(variabledescgroup.variableDesc, v)
			continue
		} else if variabledescgroup.variableDesc == "GeoName" {
			aggregate[variabledescgroup.variableDesc] = v
			continue
		}

		val, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("ERROR CONVERTING INT: %v", v))
		}

		if variabledescgroup.variableDesc == "Total" {
			totalSizeForGroup = float64(val)
			continue
		}

		if existingValue, ok := aggregate[variabledescgroup.category]; ok {
			aggregate[variabledescgroup.category] = existingValue.(int) + val
		} else {
			aggregate[variabledescgroup.category] = val
		}
	}

	if calculatePercentage {
		aggregate = performPercentageCalculation(aggregate, totalSizeForGroup)
	}

	_ = aggregate
	fmt.Print(aggregate)
}

// Summarizes census variables into, more generalized categories. Categories are defined in the getVariablesToInclude function in census.go
func aggregate(header []string, dataRow []string, censusvariableslookup map[string]variablesDescGroup, calculatePercentage bool) map[string]interface{} {
	aggregate := map[string]interface{}{}
	totalSizeForGroup := 0.0

	headerDescriptions := mapHeaderToDescription(header, censusvariableslookup)

	if len(dataRow) != len(headerDescriptions) {
		panic("WE RECEIVED MORE/LESS FIELDS FROM CENSUS API")
	}

	for i, v := range dataRow {
		variabledescgroup := headerDescriptions[i]

		if utils.CheckStringInList(geoLevels, variabledescgroup.variableDesc) {
			aggregate[variabledescgroup.variableDesc] = utils.FormatLeadingZeroes(variabledescgroup.variableDesc, v)
			continue
		} else if variabledescgroup.variableDesc == "GeoName" {
			aggregate[variabledescgroup.variableDesc] = v
			continue
		}

		val, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("ERROR CONVERTING INT: %v", v))
		}

		if variabledescgroup.variableDesc == "Total" {
			totalSizeForGroup = float64(val)
			continue
		}

		if existingValue, ok := aggregate[variabledescgroup.category]; ok {
			aggregate[variabledescgroup.category] = existingValue.(int) + val
		} else {
			aggregate[variabledescgroup.category] = val
		}
	}

	if calculatePercentage {
		aggregate = performPercentageCalculation(aggregate, totalSizeForGroup)
	}

	return aggregate
}

func mapHeaderToDescription(header []string, censusvariableslookup map[string]variablesDescGroup) []variablesDescGroup {
	headerDescriptions := []variablesDescGroup{}

	for _, row := range header {
		if variabledescgroup, ok := censusvariableslookup[row]; ok {

			headerDescriptions = append(headerDescriptions, variabledescgroup)

		} else {
			fmt.Printf("\nVARIABLE NAME NOT FOUND FOR: %v", row)
		}
	}

	return headerDescriptions
}

func performPercentageCalculation(aggregate map[string]interface{}, totalSizeForGroup float64) map[string]interface{} {
	for key, val := range aggregate {
		if utils.CheckStringInList(append(geoLevels, "GeoName"), key) {
			continue
		}

		categoryTotal := float64(val.(int))

		aggregate[key] = float64(math.Round((categoryTotal/totalSizeForGroup)*1000) / 1000)
	}
	return aggregate
}

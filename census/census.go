package census

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/jonathan-oh-89/economic-data-downloader/db"
	"github.com/jonathan-oh-89/economic-data-downloader/model"
	"github.com/jonathan-oh-89/economic-data-downloader/utils"
)

type myInterface interface {
	get() int
	set(i int)
}

func GetCensusVariableGroups() map[string][]model.CensusVariablesGroups {
	variablegroupssurl := "https://api.census.gov/data/2019/acs/acs5/groups"

	response, err := http.Get(variablegroupssurl)
	if err != nil {
		fmt.Print("", err)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responsedata map[string][]model.CensusVariablesGroups

	err = json.Unmarshal(responseData, &responsedata)
	if err != nil {
		fmt.Print(err)
	}

	return responsedata
}

func DumpSelectedCensusVariables() {
	storeInDBDone := make(chan bool)

	lines := utils.ReadCSV("/files/censusgroups.csv")

	count := 0
	for _, line := range lines {
		if line[0] == "y" {
			processCensusVariables(line[1], line[2], line[3], storeInDBDone)
			count++
		}
	}

	for {
		select {
		case <-storeInDBDone:
			count--
		default:
			time.Sleep(1 * time.Second)
		}

		if count == 0 {
			fmt.Print("Done!")
			return
		}
	}
}

//processCensusVariables - gets variables and stores into db
func processCensusVariables(groupname string, groupdesc string, variableslink string, storeInDBDone chan bool) {
	response, err := http.Get(variableslink)
	if err != nil {
		fmt.Print("", err)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responsedata map[string]map[string]model.CensusVariablesResponse

	err = json.Unmarshal(responseData, &responsedata)
	if err != nil {
		fmt.Print(err)
	}

	var censusVariablesForDB []model.CensusVariables

	for key, variable := range responsedata["variables"] {
		if strings.Contains(variable.Label, "Margin of Error") || strings.Contains(variable.Label, "Annotation") {
			continue
		}

		variable.Label = strings.Replace(variable.Label, "Estimate!!", "", -1)
		variable.Label = strings.Replace(variable.Label, "Total:!!", "", -1)
		variable.Label = strings.Replace(variable.Label, ":", "", -1)
		variable.Label = strings.Replace(variable.Label, "\"", "", -1)
		variable.Concept = strings.Replace(variable.Concept, "\"", "", -1)

		censusVariablesForDB = append(censusVariablesForDB, model.CensusVariables{VariableID: key, Label: variable.Label, Concept: variable.Concept, GroupID: variable.Group})
	}

	go db.DumpCensusVariables(censusVariablesForDB, storeInDBDone)

}

func Test(testgroup string, skip bool) {
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

		if lines[i+1][3] != groupid {
			if groupid == testgroup {
				testAPI(apiParams[0:len(apiParams)-1], censusVariablesLookup, groupid, skip)
			}
			censusVariablesLookup = map[string]string{}
			groupid = lines[i+1][3]
			apiParams = ""
		}
	}

}

func testAPI(apiparams string, censusVariablesLookup map[string]string, groupid string, skip bool) {

	url := fmt.Sprintf("https://api.census.gov/data/2019/acs/acs5?get=NAME,%s&for=county:059,037&in=state:06", apiparams)

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

	if !skip {
		if groupid == "B01001" {
			AggregateAge(arr, censusVariablesLookup, groupid)
		} else if groupid == "B03002" {
			AggregateRace(arr, censusVariablesLookup, groupid)
		} else if groupid == "B19013" {
			AggregateMedianIncome(arr, censusVariablesLookup, groupid)
		} else if groupid == "B08301" {
			AggregateTransportationToWork(arr, censusVariablesLookup, groupid)
		} else if groupid == "B08303" {
			AggregateTimeTravelToWork(arr, censusVariablesLookup, groupid)
		}

	}

	for i, variableID := range arr[0] {
		colname := censusVariablesLookup[variableID]

		if colname == "" {
			continue
		}

		arr[0][i] = colname
	}

	df := dataframe.LoadRecords(arr)

	utils.GotaToCSV(df)
}

func CheckAPI() {
	//https://api.census.gov/data/2019/pep/population?get=COUNTY,DATE_CODE,DATE_DESC,DENSITY,POP,NAME,STATE&for=state:01

}

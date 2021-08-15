package census

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jonathan-oh-89/economic-data-downloader/db"
	"github.com/jonathan-oh-89/economic-data-downloader/model"
	"github.com/jonathan-oh-89/economic-data-downloader/utils"
)

func DumpCensusGeoFips(geoLevel string) {
	lines := utils.ReadCSV("/files/cbsa2fipsxw.csv")

	db.MongoStoreGeo(lines, geoLevel)
}

func DumpCensusVariableGroups() {
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

	var censusVariableGroups map[string][]model.CensusVariablesGroups

	err = json.Unmarshal(responseData, &censusVariableGroups)
	if err != nil {
		fmt.Print(err)
	}

	db.InitializeCensusGroups(censusVariableGroups["groups"])
}

func DumpSelectedCensusVariables() {
	storeInDBDone := make(chan bool)

	lines := utils.ReadCSV("/files/censusgroups.csv")

	count := 0
	for _, line := range lines {
		if line[0] == "y" {
			cleanAndStoreCensusVariables(line[1], line[2], line[3], storeInDBDone)
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

//cleanAndStoreCensusVariables - gets variables and stores into db
func cleanAndStoreCensusVariables(groupname string, groupdesc string, variableslink string, storeInDBDone chan bool) {
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

	go db.InitializeCensusVariables(censusVariablesForDB, storeInDBDone)
}

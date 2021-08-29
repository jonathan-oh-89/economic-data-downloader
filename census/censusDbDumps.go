package census

import (
	"bufio"
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
	"go.mongodb.org/mongo-driver/mongo"
)

type CensusGeoIds struct {
	CountyFullCode string `json:"countyfullcode"`
	FipsCountyCode string `json:"countyfipscode"`
	FipsStateCode  string
	CensusYear     int
}

func DumpCensusTracts(censusyear int) {
	/*
		2020 Census Tract URLs: https://www.census.gov/geographies/reference-maps/2020/geo/2020pl-maps/2020-census-tract.html
	*/
	stateMap := db.MongoGetStatesMap()
	xlFile := utils.ReadXLSX("/files/countytractfiles.xlsx")
	maxMongoConnections := make(chan int, 200)
	var success chan bool
	runNumber := 0
	totalRuns := 0
	mongodbclient := db.ConnectToMongo()

	for _, sheet := range xlFile.Sheets {
		//exclude column row
		totalRuns = len(sheet.Rows) - 1
		success = make(chan bool, totalRuns)
		for i, row := range sheet.Rows {
			if i < 1 {
				continue
			}

			if row.Cells[3].Value != fmt.Sprintf("%v", censusyear) {
				continue
			}

			fileName := row.Cells[0].Value
			stateinfo := stateMap[row.Cells[1].Value]
			statefipscode := stateinfo.FipsStateCode
			stateabbreviation := strings.ToLower(stateinfo.StateAbbreviation)
			countyfipscode := row.Cells[2].Value
			countyfullcode := statefipscode + countyfipscode

			geoInfo := CensusGeoIds{
				CountyFullCode: countyfullcode,
				FipsCountyCode: countyfipscode,
				FipsStateCode:  statefipscode,
				CensusYear:     censusyear,
			}

			// fix the url in use below
			url := ""
			if geoInfo.CensusYear == 2020 {
				url = fmt.Sprintf("https://www2.census.gov/geo/maps/DC2020/PL20/st%s_%s/censustract_maps/%sDC20CT_C%s_CT2MS.txt", geoInfo.FipsStateCode, stateabbreviation, fileName, geoInfo.CountyFullCode)
			} else {
				url = fmt.Sprintf("https://www2.census.gov/geo/maps/dc10map/tract/st%s_%s/%sDC10CT_C%s_CT2MS.txt", geoInfo.FipsStateCode, stateabbreviation, fileName, geoInfo.CountyFullCode)
			}

			runNumber++
			maxMongoConnections <- runNumber
			log.Printf("Running run#: %d", runNumber)

			go parseCensusTract(url, maxMongoConnections, success, geoInfo, mongodbclient)
		}
	}

	for {
		select {
		case <-success:
			totalRuns--
		default:
			fmt.Print("\nWaiting for all tracts to finish")
			time.Sleep(1 * time.Second)
		}

		if totalRuns == 0 {
			fmt.Print("Done storing tracts!")
			return
		}
	}
}

func parseCensusTract(url string, maxMongoConnections chan int, success chan bool, geoInfo CensusGeoIds, client *mongo.Client) {
	tracts := make([]model.TractInfo, 0)

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("ERROR getting response from: %s", url)
		return
	}

	sc := bufio.NewScanner(response.Body)
	lines := make([]string, 0)

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	for i, line := range lines {
		if i < 1 {
			continue
		}

		lineSplit := strings.Split(line, ";")

		tractcode := lineSplit[1]

		tracts = append(tracts, model.TractInfo{
			TractCode:      tractcode,
			CensusYear:     geoInfo.CensusYear,
			CountyFullCode: geoInfo.CountyFullCode,
			FipsStateCode:  geoInfo.FipsStateCode,
			FipsCountyCode: geoInfo.FipsCountyCode,
		})
	}
	db.MongoStoreTracts(tracts, maxMongoConnections, success, client)
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

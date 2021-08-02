package census

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jonathan-oh-89/economic-data-downloader/model"
)

func Test() {
	// "https://api.census.gov/data/2019/pep/population?get=COUNTY,DATE_CODE,DATE_DESC,DENSITY,POP,NAME,STATE&for=state:01",

	urls := []string{
		// "https://api.census.gov/data/2018/acs/acs5?get=NAME,B25034_001E,B25034_002E&for=tract:*&in=state:06&in=county:*",
		"https://api.census.gov/data/2019/acs/acs5?get=NAME,B25034_001E,B25034_002E&for=tract:*&in=state:06&in=county:*",
	}

	for _, url := range urls {
		response, err := http.Get(url)
		if err != nil {
			fmt.Print("", err)
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

		for i, row := range arr {
			//Skip first 4 rows - just description
			if i < 1 {
				continue
			}

			// groupid := strings.Split(row[0], "_")[0]

			// //Skip Puerto Rico
			// if groupid[len(groupid)-2:] == "PR" {
			// 	continue
			// }

			fmt.Print(row, "\n")

			// censusVariableList = append(censusVariableList, CensusVariables{
			// 	VariableID:  row[0],
			// 	Name:        row[1],
			// 	Description: row[2],
			// 	GroupID:     groupid,
			// })

			i++

			if i > 10 {
				break
			}
		}

	}
	// fmt.Print(censusVariableList, "\n")

	// users := []User{
	// 	{"Aram", 17, 0.2, true},
	// 	{"Juan", 18, 0.8, true},
	// 	{"Ana", 22, 0.5, true},
	// }
	// df := dataframe.LoadStructs(censusVariableList)

	// fmt.Print("df:\n", df)

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

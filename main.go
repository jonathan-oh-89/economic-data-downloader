package main

import (
	"fmt"

	"github.com/jonathan-oh-89/economic-data-downloader/census"
	// "github.com/go-gota/gota/dataframe"
)

func main() {

	// db.Test()

	fmt.Print("Startimg")

	if true {
		census.DumpCensusGeoFips("state")
	}

	if false {
		// db.InitializeDB()
		census.DumpCensusVariableGroups()
		census.DumpSelectedCensusVariables()
	}

	// // Select groups to store in csv and run

	// census.DownloadToCSV("B15003")

	census.Test("B11012", "county")
	census.Do("06", "county")

	// census.CheckAPI()
}

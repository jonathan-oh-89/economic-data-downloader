package main

import (
	"fmt"

	"github.com/jonathan-oh-89/economic-data-downloader/census"
	"github.com/jonathan-oh-89/economic-data-downloader/db"
	// "github.com/go-gota/gota/dataframe"
)

func main() {

	// db.Test()

	fmt.Print("Starting")
	db.MongoGetCbsaCodes()

	if false {
		census.DumpCensusGeoFips("state")
		census.DumpCensusGeoFips("county")
		census.DumpCensusGeoFips("cbsa")
	}

	if false {
		//setup mysql database
		db.InitializeDB()
		census.DumpCensusVariableGroups()
		census.DumpSelectedCensusVariables()
	}

	// CENSUS SECTION
	// census.DownloadToCSV("B15003")
	// census.Test("B11012", "county")
	// census.Do("06", "county")

	census.GetBuildingPermits()
}

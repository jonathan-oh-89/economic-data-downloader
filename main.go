package main

import (
	"github.com/jonathan-oh-89/economic-data-downloader/census"
	"github.com/jonathan-oh-89/economic-data-downloader/mongoclient"
	// "github.com/go-gota/gota/dataframe"
)

func main() {

	// db.Test()

	if true {
		mongoclient.MongoStoreGeo()
	}

	if false {
		// db.InitializeDB()
		census.DumpCensusVariableGroups()
		census.DumpSelectedCensusVariables()
	}

	// // Select groups to store in csv and run

	// census.DownloadToCSV("B15003")

	census.DumpCensusGeoFips()

	census.Test("B11012", "county")
	census.Do("06", "county")

	// census.CheckAPI()
}

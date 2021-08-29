package main

import (
	"log"

	"github.com/jonathan-oh-89/economic-data-downloader/census"
	"github.com/jonathan-oh-89/economic-data-downloader/db"
	"github.com/jonathan-oh-89/economic-data-downloader/esri"
	// "github.com/go-gota/gota/dataframe"
)

func main() {

	// db.Test()

	log.Print("Starting")
	esri.DumpEsriTractData(50)

	if false {
		//setup mysql database
		db.InitializeDB()
		db.MongoStoreGeo()
		census.DumpCensusTracts(2010)
		census.DumpCensusTracts(2020)
		census.DumpCensusVariableGroups()
		census.DumpSelectedCensusVariables()
	}

	// CENSUS SECTION
	census.DownloadToCSV("B25056")
	// census.Test("B11012", "county")
	// census.Do("06", "county")

	//BUILDING PERMITS - map msa - dates &total housing units permitted
	// census.GetBuildingPermits()

	log.Print("Finished running")
}

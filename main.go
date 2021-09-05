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

	if false {
		//setup mysql database
		db.InitializeDB()
		db.MongoStoreGeo()
		esri.DumpEsriTractData(200)
		census.DumpCensusTracts(2010)
		census.DumpCensusTracts(2020)
		census.DumpCensusVariableGroups()
		census.DumpSelectedCensusVariables()
	}

	// CENSUS SECTION
	// census.DownloadToCSV("B25056")
	// census.Test("B11012", "county")
	// census.Do("06", "county")

	//BUILDING PERMITS - map msa - dates &total housing units permitted
	// census.GetBuildingPermits()

	//CRIME SECTION

	if false {
		//carefull with running this
		countiesToGetCrime := []string{"29189", "45045"}
		esri.DumpEsriCrimeData(len(countiesToGetCrime), 2021, countiesToGetCrime)
	}

	log.Print("Finished running")
}

/*
Misc scripts


*****	Get count of tracts per county	******
	// tracts := db.MongoGetEsriTractsList()

	// countytractcounty := make(map[string]int, 0)

	// for _, v := range tracts {
	// 	if _, ok := countytractcounty[v.CountyFullCode]; ok {
	// 		countytractcounty[v.CountyFullCode] = countytractcounty[v.CountyFullCode] + 1
	// 		continue
	// 	} else {
	// 		countytractcounty[v.CountyFullCode] = 1
	// 	}
	// }

	// mc := db.ConnectToMongo()
	// db.MongoStoreTempMap(countytractcounty, mc)

*/

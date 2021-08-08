package main

import (
	"github.com/jonathan-oh-89/economic-data-downloader/census"
	// "github.com/jonathan-oh-89/economic-data-downloader/config"
	"github.com/jonathan-oh-89/economic-data-downloader/db"
	// "github.com/go-gota/gota/dataframe"
)

func main() {

	// db.Test()

	if false {
		db.InitializeDB()
		grouplist := census.GetCensusVariableGroups()["groups"]
		db.DumpCensusVariableGroups(grouplist)
	}

	// // Select groups to store in csv and run
	// census.DumpSelectedCensusVariables()

	census.Test("B08301", false)
	// census.CheckAPI()
}

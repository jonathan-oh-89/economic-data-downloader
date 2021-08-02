package main

import (
	"github.com/jonathan-oh-89/economic-data-downloader/census"
	"github.com/jonathan-oh-89/economic-data-downloader/db"
	// "github.com/go-gota/gota/dataframe"
)

func main() {
	if false {
		db.InitializeDB()
		grouplist := census.GetCensusVariableGroups()["groups"]
		db.DumpCensusVariableGroups(grouplist)
	}

	census.Test()
}

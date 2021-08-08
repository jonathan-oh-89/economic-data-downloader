package utils

import (
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func GotaToCSV(df dataframe.DataFrame) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(wd + "/files/dataframeresults.csv")
	if err != nil {
		log.Fatal(err)
	}

	df.WriteCSV(f)
}

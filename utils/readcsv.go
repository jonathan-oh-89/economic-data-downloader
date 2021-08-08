package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func GetCSVFile(path string) *os.File {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(wd + path)
	if err != nil {
		fmt.Print("Error: ", err)
	}
	defer f.Close()

	return f
}

func ReadCSV(path string) [][]string {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(wd + path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // this needs to be after the err check

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return lines
}

package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/jonathan-oh-89/economic-data-downloader/model"
)

func Test() {
	db := connectToDb()
	db.Close()

	fmt.Print("connected to db: ", db)

}

func connectToDb() *sql.DB {
	// dbConfig := c.Config
	err := godotenv.Load(".env")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	un := os.Getenv("DB_USERNAME")
	pw := os.Getenv("DB_PASSWORD")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", un, pw, host, port))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + "Census")
	if err != nil {
		panic(err)
	}

	return db
}

func InitializeDB() {
	fmt.Print("db initialize disabled")
	// return

	db := connectToDb()
	defer db.Close()

	// _, err := db.Exec("CREATE DATABASE " + "Census")

	// if err != nil {
	// 	panic(err)
	// }

	_, err := db.Exec("CREATE TABLE census_variable_groups ( groupid varchar(10), description varchar(255), variableslink varchar(100) )")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE census_variables ( variableid varchar(12), label varchar(255), concept varchar(100), groupid varchar(10) )")
	if err != nil {
		panic(err)
	}
}

func DumpCensusVariableGroups(cvglist []model.CensusVariablesGroups) {
	db := connectToDb()
	defer db.Close()

	sqlCall := "INSERT INTO census_variable_groups(groupid, description, variableslink) VALUES "
	for _, row := range cvglist {
		sqlCall += fmt.Sprintf(`("%s", "%s", "%s"),`, row.Name, row.Description, row.Variables)
	}

	// Get rid of last comma
	sqlCall = sqlCall[0 : len(sqlCall)-1]

	_, err := db.Exec(sqlCall)
	if err != nil {
		fmt.Print("SQL ERROR: ", err)
		panic(err)
	}
}

func DumpCensusVariables(censusVariablesForDB []model.CensusVariables, storeInDBDone chan bool) {

	db := connectToDb()
	defer db.Close()

	sqlCall := "INSERT INTO census_variables(variableid, label, concept, groupid) VALUES "

	for _, cv := range censusVariablesForDB {
		sqlCall += fmt.Sprintf(`("%s", "%s", "%s", "%s"),`, cv.VariableID, cv.Label, cv.Concept, cv.GroupID)
	}

	// Get rid of last comma
	sqlCall = sqlCall[0 : len(sqlCall)-1]

	_, err := db.Exec(sqlCall)
	if err != nil {

		fmt.Printf("SQL ERROR: %s", err)
		panic(err)
	}

	storeInDBDone <- true

}

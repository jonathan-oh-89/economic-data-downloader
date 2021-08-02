package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jonathan-oh-89/economic-data-downloader/model"
)

func connectToDb() *sql.DB {
	db, err := sql.Open("mysql",
		"un:pw@tcp(db:3306)/")
	if err != nil {
		fmt.Print(err)
	}

	_, err = db.Exec("USE " + "Census")
	if err != nil {
		panic(err)
	}

	return db
}

func InitializeDB() {
	fmt.Print("temporarily disabled")
	return

	db := connectToDb()
	defer db.Close()

	_, err := db.Exec("CREATE DATABASE " + "Census")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE census_variable_groups ( name varchar(10), description varchar(255), variableslink varchar(100) )")
	if err != nil {
		panic(err)
	}
}

func DumpCensusVariableGroups(cvglist []model.CensusVariablesGroups) {

	db := connectToDb()
	defer db.Close()

	vals := []interface{}{}

	sqlCall := "INSERT INTO census_variable_groups(name, description, variableslink) VALUES "
	for _, row := range cvglist {
		sqlCall += fmt.Sprintf(`("%s", "%s", "%s"),`, row.Name, row.Description, row.Variables)

		vals = append(vals, row)
	}

	stmt, err := db.Prepare(sqlCall)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}

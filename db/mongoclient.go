package db

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/jonathan-oh-89/economic-data-downloader/model"
	"github.com/jonathan-oh-89/economic-data-downloader/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToMongo() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	host := os.Getenv("MONGO_HOST")
	database := os.Getenv("MONGO_DATABASE")
	un := os.Getenv("MONGO_USERNAME")
	pw := os.Getenv("MONGO_PASSWORD")

	// clientOptions := options.Client().ApplyURI("mongodb+srv://admin:<password>@scopeout.hdtom.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", un, pw, host, database))
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

func MongoStoreGeo(lines [][]string, geoLevel string) {
	// fmt.Print("Mongo initialize disabled")
	// return

	statesArray := []interface{}{}

	for i, line := range lines {

		if i < 1 {
			continue
		}

		stateInfo := model.StateInfo{FipsStateCode: line[9], StateName: line[8]}

		data, err := utils.MarshallStructtoBson(stateInfo)

		statesArray = append(statesArray, data)

		if err != nil {
			panic(err)
		}
	}

	client := connectToMongo()

	collection := client.Database("scopeout").Collection("State")

	_, err := collection.InsertMany(context.TODO(), statesArray)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done storing geo: ", geoLevel)
}

func MongoCensusGroups(cvglist []model.CensusVariablesGroups) {

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

func MongoCensusVariables(censusVariablesForDB []model.CensusVariables, storeInDBDone chan bool) {

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

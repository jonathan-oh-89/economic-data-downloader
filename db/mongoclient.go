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

	collectionName := ""
	stateInfo := model.StateInfo{}
	countyInfo := model.CountyInfo{}
	statesArray := []interface{}{}
	data := []byte{}

	for i, line := range lines {

		if i < 2 {
			continue
		}

		switch geoLevel {
		case "state":
			stateInfo = model.StateInfo{FipsStateCode: line[9], StateName: line[8]}
			data, _ = utils.MarshallStructtoBson(stateInfo)
			collectionName = "State"
		case "county":
			stateInfo = model.StateInfo{FipsStateCode: line[9], StateName: line[8]}
			countyInfo = model.CountyInfo{FipsCountyCode: line[10], CountyCountyEquivalent: line[7], StateInfo: stateInfo}
			data, _ = utils.MarshallStructtoBson(countyInfo)
			collectionName = "County"
		default:

		}

		statesArray = append(statesArray, data)
	}

	client := connectToMongo()

	collection := client.Database("scopeout").Collection(collectionName)

	_, err := collection.InsertMany(context.TODO(), statesArray)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done storing geo: ", geoLevel)
}

func getGeoMongoCollection(geoLevel string) {

}

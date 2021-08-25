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
	"go.mongodb.org/mongo-driver/bson"
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

func MongoGetCbsaCodes() []string {
	cbsa := make([]string, 0)

	client := connectToMongo()
	collection := client.Database("scopeout").Collection("Cbsa")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []map[string]interface{}
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, record := range results {
		str := fmt.Sprintf("%v", record["cbsacode"])

		cbsa = append(cbsa, str)
		_ = cbsa
	}

	return cbsa
}

func MongoStoreGeo(lines [][]string, geoLevel string) {

	cbsaCounties := make(map[string][]model.CountyInfo)

	collectionName := ""
	stateInfo := model.StateInfo{}
	countyInfo := model.CountyInfo{}
	geoArray := []interface{}{}
	data := []byte{}

	//gather list of counties/state per cbsa prior assigning them
	if geoLevel == "cbsa" {
		for i, line := range lines {
			if i < 2 {
				continue
			}

			stateInfo = model.StateInfo{FipsStateCode: line[9], StateName: line[8]}
			countyInfo = model.CountyInfo{FipsCountyCode: line[10], CountyCountyEquivalent: line[7], StateInfo: stateInfo}

			if _, ok := cbsaCounties[line[0]]; ok {
				cbsaCounties[line[0]] = append(cbsaCounties[line[0]], countyInfo)
			} else {
				cbsaCounties[line[0]] = []model.CountyInfo{countyInfo}
			}
		}
	}

	dupcheck := make(map[string]bool)

	for i, line := range lines {

		if i < 2 {
			continue
		}

		//skip puerto rico
		if line[9] == "72" {
			continue
		}

		switch geoLevel {
		case "state":
			if _, hasKey := dupcheck[line[9]]; hasKey {
				continue
			} else {
				dupcheck[line[9]] = true
				stateInfo = model.StateInfo{FipsStateCode: line[9], StateName: line[8]}
				data, _ = utils.MarshallStructtoBson(stateInfo)
				collectionName = "State"
			}
		case "county":
			if _, hasKey := dupcheck[line[9]+line[10]]; hasKey {
				continue
			} else {
				dupcheck[line[9]+line[10]] = true
				stateInfo = model.StateInfo{FipsStateCode: line[9], StateName: line[8]}
				countyInfo = model.CountyInfo{FipsCountyCode: line[10], CountyCountyEquivalent: line[7], StateInfo: stateInfo}
				data, _ = utils.MarshallStructtoBson(countyInfo)
				collectionName = "County"
			}
		case "cbsa":
			if _, hasKey := dupcheck[line[0]]; hasKey {
				continue
			} else {
				dupcheck[line[0]] = true
				countiesToGet := cbsaCounties[line[0]]
				cbsaInfo := model.CBSAInfo{CbsaCode: line[0], CbsaTitle: line[3], Counties: countiesToGet}
				data, _ = utils.MarshallStructtoBson(cbsaInfo)
				collectionName = "Cbsa"
			}

		default:

		}

		geoArray = append(geoArray, data)
	}

	client := connectToMongo()

	collection := client.Database("scopeout").Collection(collectionName)

	_, err := collection.InsertMany(context.TODO(), geoArray)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done storing geo: ", geoLevel)
}

func getGeoMongoCollection(geoLevel string) {

}

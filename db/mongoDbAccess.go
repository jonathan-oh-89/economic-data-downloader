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

func ConnectToMongo() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	host := os.Getenv("MONGO_HOST")
	database := os.Getenv("MONGO_DATABASE")
	un := os.Getenv("MONGO_USERNAME")
	pw := os.Getenv("MONGO_PASSWORD")

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", un, pw, host, database))
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Connected to MongoDB")

	return client
}

func MongoGetCbsaCodes() map[string]model.CBSAInfo {

	client := ConnectToMongo()
	collection := client.Database("scopeout").Collection("Cbsa")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []model.CBSAInfo
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	cbsaMap := make(map[string]model.CBSAInfo)

	for _, record := range results {
		cbsaMap[record.CbsaCode] = record
	}

	return cbsaMap
}

func MongoGetStatesMap() map[string]model.StateInfo {

	client := ConnectToMongo()

	collection := client.Database("scopeout").Collection("State")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []model.StateInfo
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	stateMap := make(map[string]model.StateInfo)

	for _, record := range results {
		stateMap[record.FipsStateCode] = record
	}

	return stateMap
}

func MongoGetCountiesMap() map[string]model.CountyInfo {

	client := ConnectToMongo()
	collection := client.Database("scopeout").Collection("County")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []model.CountyInfo
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	countyMap := make(map[string]model.CountyInfo)

	for _, record := range results {
		countyMap[record.StateInfo.FipsStateCode+record.FipsCountyCode] = record
	}

	return countyMap
}

func MongoGetEsriTractsMap() map[string]model.EsriTractsInfo {

	client := ConnectToMongo()
	collection := client.Database("scopeout").Collection("EsriTracts")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []model.EsriTractsInfo
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	esriTractsMap := make(map[string]model.EsriTractsInfo)

	for _, record := range results {
		esriTractsMap[record.CountyFullCode] = record
	}

	return esriTractsMap
}

func MongoStoreGeo() {
	stateInfo := model.StateInfo{}
	stateLookup := make(map[string]model.StateInfo)
	countyInfo := model.CountyInfo{}
	countyLookup := make(map[string]model.CountyInfo)
	cbsaCounties := make(map[string][]model.CountyInfo)
	cbsadupcheck := make(map[string]bool)
	data := []byte{}

	collectionName := ""

	for _, geoLevel := range []string{"state", "county", "cbsa"} {
		geoArray := []interface{}{}
		lines := [][]string{}

		if geoLevel == "state" || geoLevel == "county" {
			lines = utils.ReadCSV("/files/statecountyfips.csv")
		} else if geoLevel == "cbsa" {
			lines = utils.ReadCSV("/files/cbsafips.csv")

			for i, line := range lines {
				if i < 2 {
					continue
				}

				stateInfo = stateLookup[line[9]]
				countyInfo = countyLookup[line[9]+line[10]]

				if _, ok := cbsaCounties[line[0]]; ok {
					cbsaCounties[line[0]] = append(cbsaCounties[line[0]], countyInfo)
				} else {
					cbsaCounties[line[0]] = []model.CountyInfo{countyInfo}
				}
			}

		}

		for i, line := range lines {
			if i < 2 {
				continue
			}

			switch geoLevel {
			case "state":
				geoLevelID := utils.FormatLeadingZeroes("3", line[0])
				statecode := utils.FormatLeadingZeroes("2", line[1])
				stateabbreviation := utils.FormatLeadingZeroes("2", line[7])

				//Skip puerto rico and if not state
				if statecode == "72" || geoLevelID != "040" {
					continue
				}

				if _, hasKey := stateLookup[statecode]; hasKey {
					continue
				} else {
					stateLookup[statecode] = model.StateInfo{FipsStateCode: statecode, StateName: line[6], StateAbbreviation: stateabbreviation}
					stateInfo = model.StateInfo{FipsStateCode: statecode, StateName: line[6], StateAbbreviation: stateabbreviation}
					data, _ = utils.MarshallStructtoBson(stateInfo)
					collectionName = "State"
				}
			case "county":
				geoLevelID := utils.FormatLeadingZeroes("3", line[0])
				statecode := utils.FormatLeadingZeroes("2", line[1])
				countycode := utils.FormatLeadingZeroes("3", line[2])

				//Skip puerto rico and if not county
				if statecode == "72" || geoLevelID != "050" {
					continue
				}

				if _, hasKey := countyLookup[statecode+countycode]; hasKey {
					continue
				} else {
					stateInfo = stateLookup[statecode]
					stateCountyCode := statecode + countycode
					countyLookup[stateCountyCode] = model.CountyInfo{CountyFullCode: stateCountyCode, FipsCountyCode: countycode, CountyName: line[6], StateInfo: stateInfo}
					countyInfo = model.CountyInfo{CountyFullCode: stateCountyCode, FipsCountyCode: countycode, CountyName: line[6], StateInfo: stateInfo}
					data, _ = utils.MarshallStructtoBson(countyInfo)
					collectionName = "County"
				}
			case "cbsa":
				//skip puerto rico
				if line[9] == "72" {
					continue
				}

				if _, hasKey := cbsadupcheck[line[0]]; hasKey {
					continue
				} else {
					cbsadupcheck[line[0]] = true
					cbsaCounties := cbsaCounties[line[0]]
					cbsaInfo := model.CBSAInfo{CbsaCode: line[0], CbsaTitle: line[3], Counties: cbsaCounties}
					data, _ = utils.MarshallStructtoBson(cbsaInfo)
					collectionName = "Cbsa"
				}
			}

			geoArray = append(geoArray, data)
		}

		_ = collectionName
		client := ConnectToMongo()
		collection := client.Database("scopeout").Collection(collectionName)

		_, err := collection.InsertMany(context.TODO(), geoArray)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Finished storing geo")
}

func MongoStoreTracts(tracts []model.TractInfo, maxMongoConnections <-chan int, success chan bool, client *mongo.Client) {
	geoArray := []interface{}{}

	collection := client.Database("scopeout").Collection("Tracts")

	for _, tract := range tracts {
		data, _ := utils.MarshallStructtoBson(tract)
		geoArray = append(geoArray, data)
	}

	_, err := collection.InsertMany(context.TODO(), geoArray)
	if err != nil {
		log.Fatal(err)
	}

	runNumber := <-maxMongoConnections
	fmt.Printf("\nProcessed run#: %d", runNumber)

	success <- true
}

func MongoStoreEsriTracts(esriTracts []model.EsriTractsInfo, client *mongo.Client) {
	geoArray := []interface{}{}

	collection := client.Database("scopeout").Collection("EsriTracts")

	for _, tract := range esriTracts {
		data, _ := utils.MarshallStructtoBson(tract)
		geoArray = append(geoArray, data)
	}

	_, err := collection.InsertMany(context.TODO(), geoArray)
	if err != nil {
		log.Fatal(err)
	}
}

package esri

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jonathan-oh-89/economic-data-downloader/db"
	"github.com/jonathan-oh-89/economic-data-downloader/model"
	"github.com/jonathan-oh-89/economic-data-downloader/utils"
)

func esriApi(url string, params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	for k, v := range params {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func getEsriToken() string {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	client_id := os.Getenv("ESRI_OAUTH_CLIENT_ID_JAYLEEONG0913")
	client_secret := os.Getenv("ESRI_OAUTH_CLIENT_SECRET_JAYLEEONG0913")

	params := map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"grant_type":    "client_credentials",
	}

	resp, err := esriApi("https://www.arcgis.com/sharing/rest/oauth2/token", params)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var target struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		panic(err)
	}
	return target.AccessToken
}

func DumpEsriTractData(howManyApiHits int) {
	token := getEsriToken()

	countiesLeftToRun := checkIfCountyForTractsExists()

	if len(countiesLeftToRun) < 1 {
		return
	}

	for i, county := range countiesLeftToRun {
		if i > howManyApiHits {
			break
		}

		params := map[string]string{
			"sourceCountry":           "US",
			"geographylayers":         "US.Counties",
			"geographyids":            fmt.Sprintf("%s", county.CountyFullCode),
			"returnGeometry":          "true",
			"returnSubGeographyLayer": "true",
			"subGeographyLayer":       "US.Tracts",
			"generalizationLevel":     "6",
			"f":                       "pjson",
			"featureLimit":            "1000",
			"token":                   token,
		}

		log.Printf("Calling Standard Geography for: %s", county.CountyFullCode)

		utils.RandomTimeOut()

		resp, err := esriApi("https://geoenrich.arcgis.com/arcgis/rest/services/World/geoenrichmentserver/StandardGeographyQuery", params)
		if err != nil {
			log.Fatalf("FATAL ERROR: %s", err.Error())
		}
		defer resp.Body.Close()

		var structtarget model.EsriStandardGeoResponse

		err = json.NewDecoder(resp.Body).Decode(&structtarget)
		if err != nil {
			log.Fatalf("FATAL ERROR: %s", err.Error())
		}

		esriTracts := make([]model.EsriTractsInfo, 0)

		if len(structtarget.Results[0].Value.Features) < 1 {
			log.Printf("Warning: Receive no features for county: %s", county.CountyFullCode)
			continue
		}

		for i, data := range structtarget.Results[0].Value.Features {
			esriTracts = append(esriTracts, model.EsriTractsInfo{
				TractCode:               data.Attributes.AreaID,
				CountyFullCode:          county.CountyFullCode,
				FipsStateCode:           county.StateInfo.FipsStateCode,
				EsriStandardGeoFeatures: data,
			})

			if i >= 1000 {
				log.Fatal("ERROR: Standard geography returned more than 1000 records")
			}
		}

		mongodbclient := db.ConnectToMongo()
		db.MongoStoreEsriTracts(esriTracts, mongodbclient)
	}

	log.Print("Finished storing Esri Tracts")

}

func checkIfCountyForTractsExists() []model.CountyInfo {
	allCounties := db.MongoGetCountiesMap()

	existingCounties := db.MongoGetEsriTractsMap()

	countiesLeftToRun := make([]model.CountyInfo, 0)

	countiesSkipped := 0

	for k, v := range allCounties {
		if _, ok := existingCounties[k]; ok {
			countiesSkipped++
			continue
		}
		countiesLeftToRun = append(countiesLeftToRun, v)
	}

	log.Printf("Skipped %d/%d counties", countiesSkipped, len(allCounties))

	return countiesLeftToRun
}

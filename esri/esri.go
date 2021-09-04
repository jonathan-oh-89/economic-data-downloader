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

func DumpEsriCrimeData(howManyApiHits int, crimeyear int) {
	token := getEsriToken()

	countiesToRun := checkIfCountyForCrimeExists()

	if len(countiesToRun) < 1 {
		return
	}

	for i, v := range countiesToRun {
		if i >= howManyApiHits {
			break
		}

		countyId := v.CountyFullCode

		studyAreas := fmt.Sprintf("[{\"sourceCountry\":\"US\", \"layer\":\"US.Counties\",\"ids\":[\"%s\"],\"comparisonLevels\": [{\"layer\":\"US.Tracts\"}] }]", countyId)

		params := map[string]string{
			"studyAreas":        studyAreas,
			"analysisVariables": "[\"CRMCYPERC\",\"CRMCYMURD\",\"CRMCYRAPE\",\"CRMCYROBB\",\"CRMCYASST\",\"CRMCYPROC\",\"CRMCYBURG\",\"CRMCYLARC\",\"CRMCYMVEH\"]",
			"f":                 "pjson",
			"returnGeometry":    "false",
			"token":             token,
		}

		resp, err := esriApi("https://geoenrich.arcgis.com/arcgis/rest/services/World/geoenrichmentserver/GeoEnrichment/enrich", params)
		if err != nil {
			log.Fatalf("FATAL ERROR: %s", err.Error())
		}
		defer resp.Body.Close()

		var structtarget model.EsriEnrichResponse

		err = json.NewDecoder(resp.Body).Decode(&structtarget)
		if err != nil {
			log.Fatalf("FATAL ERROR: %s", err.Error())
		}

		if len(structtarget.Messages) > 0 {
			for _, message := range structtarget.Messages {
				log.Printf("Got a message from ESRI API: %s", message)
			}
		}

		allCrimeFeatures := structtarget.Results[0].Value.FeatureSet[0].Features

		esriCrimeCounty := model.EsriCrimeCountyInfo{}
		esriCrimeTracts := make([]model.EsriCrimeTractInfo, 0)

		if len(allCrimeFeatures) < 1 {
			log.Printf("!!!!! Warning: Receive no features for county: %s !!!!!", countyId)
			continue
		}

		for i, data := range allCrimeFeatures {
			if data.Attributes.StdGeographyLevel == "US.Counties" {
				esriCrimeCounty = model.EsriCrimeCountyInfo{
					CountyFullCode:    countyId,
					CrimeYear:         crimeyear,
					StdGeographyLevel: data.Attributes.StdGeographyLevel,
					StdGeographyName:  data.Attributes.StdGeographyName,
					StdGeographyID:    data.Attributes.StdGeographyID,
					CRMCYPERC:         data.Attributes.CRMCYPERC,
					CRMCYMURD:         data.Attributes.CRMCYMURD,
					CRMCYRAPE:         data.Attributes.CRMCYRAPE,
					CRMCYROBB:         data.Attributes.CRMCYROBB,
					CRMCYASST:         data.Attributes.CRMCYASST,
					CRMCYPROC:         data.Attributes.CRMCYPROC,
					CRMCYBURG:         data.Attributes.CRMCYBURG,
					CRMCYLARC:         data.Attributes.CRMCYLARC,
					CRMCYMVEH:         data.Attributes.CRMCYMVEH,
				}
			} else if data.Attributes.StdGeographyLevel == "US.Tracts" {
				esriCrimeTracts = append(esriCrimeTracts, model.EsriCrimeTractInfo{
					StdGeographyLevel: data.Attributes.StdGeographyLevel,
					StdGeographyName:  data.Attributes.StdGeographyName,
					StdGeographyID:    data.Attributes.StdGeographyID,
					CRMCYPERC:         data.Attributes.CRMCYPERC,
					CRMCYMURD:         data.Attributes.CRMCYMURD,
					CRMCYRAPE:         data.Attributes.CRMCYRAPE,
					CRMCYROBB:         data.Attributes.CRMCYROBB,
					CRMCYASST:         data.Attributes.CRMCYASST,
					CRMCYPROC:         data.Attributes.CRMCYPROC,
					CRMCYBURG:         data.Attributes.CRMCYBURG,
					CRMCYLARC:         data.Attributes.CRMCYLARC,
					CRMCYMVEH:         data.Attributes.CRMCYMVEH,
				})
			}

			if i >= 1000 {
				log.Fatal("ERROR: Standard geography returned more than 1000 records")
			}
		}

		esriCrimeCounty.TractsCrime = esriCrimeTracts

		mongodbclient := db.ConnectToMongo()
		db.MongoStoreEsriCrime(esriCrimeCounty, mongodbclient)
	}

	log.Print("Finished storing Esri Crime")
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
			"featureLimit":            "5000",
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

		standardGeoFeatures := structtarget.Results[0].Value.Features

		if len(standardGeoFeatures) < 1 {
			log.Printf("Warning: Receive no features for county: %s", county.CountyFullCode)
			continue
		}

		for _, data := range standardGeoFeatures {
			esriTracts = append(esriTracts, model.EsriTractsInfo{
				TractCode:               data.Attributes.AreaID,
				CountyFullCode:          county.CountyFullCode,
				FipsStateCode:           county.StateInfo.FipsStateCode,
				EsriStandardGeoFeatures: data,
			})
		}

		mongodbclient := db.ConnectToMongo()
		db.MongoStoreEsriTracts(esriTracts, mongodbclient)
	}

	log.Print("Finished storing Esri Tracts")
}

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

func checkIfCountyForCrimeExists() []model.CountyInfo {
	//Get select cbsa for now
	cbsas := db.MongoGetCbsaMap()
	cbasFilterList := []string{"31080"}
	countiesFilter := make([]string, 0)
	for k, cbsa := range cbsas {
		if utils.CheckStringInList(cbasFilterList, k) {
			for _, county := range cbsa.Counties {
				countiesFilter = append(countiesFilter, county.CountyFullCode)
			}
		}
	}

	allCounties := db.MongoGetCountiesMap()

	existingCounties := db.MongoGetEsriCrimeCounties()

	countiesLeftToRun := make([]model.CountyInfo, 0)

	countiesSkipped := 0

	for k, v := range allCounties {
		if _, ok := existingCounties[k]; ok {
			countiesSkipped++
			continue
		}

		if !utils.CheckStringInList(countiesFilter, k) {
			log.Print("Skipping county: ", k)
			continue
		}
		countiesLeftToRun = append(countiesLeftToRun, v)
	}

	log.Printf("Skipped %d/%d counties", countiesSkipped, len(allCounties))

	return countiesLeftToRun
}

func checkIfCountyForTractsExists() []model.CountyInfo {
	allCounties := db.MongoGetCountiesMap()

	existingEsriTracts := db.MongoGetEsriTractsList()

	existingCounties := make(map[string]bool)

	for _, record := range existingEsriTracts {
		existingCounties[record.CountyFullCode] = true
	}

	countiesLeftToRun := make([]model.CountyInfo, 0)

	countiesSkipped := 0

	// Counties with no data
	excludeCounties := []string{"02063", "02066"}

	for k, v := range allCounties {
		if _, ok := existingCounties[k]; ok {
			countiesSkipped++
			continue
		}

		if utils.CheckStringInList(excludeCounties, k) {
			continue
		}

		countiesLeftToRun = append(countiesLeftToRun, v)
	}

	log.Printf("Skipped %d/%d counties", countiesSkipped, len(allCounties))

	return countiesLeftToRun
}

package census

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jonathan-oh-89/economic-data-downloader/db"
	"github.com/jonathan-oh-89/economic-data-downloader/utils"
)

type BuildingPermits struct {
	dates      []string
	totalunits []string
}

func GetBuildingPermits() map[string]BuildingPermits {
	censusCbsaMap := utils.CensusCbsaVersionMap()
	cbsacodes := db.MongoGetCbsaCodes()

	buildingPermits := make(map[string]BuildingPermits, 0)

	for _, cbsa := range cbsacodes {
		buildingPermits[cbsa] = BuildingPermits{}
	}

	// years := []string{"05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21"}
	years := []string{"21"}
	months := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}

	for _, year := range years {
		fmt.Print("Starting year: ", year)

		for _, month := range months {
			url := fmt.Sprintf(`https://www2.census.gov/econ/bps/Metro/ma%s%sc.txt`, year, month)

			response, err := http.Get(url)
			if err != nil {
				fmt.Print("ERROR: ", err)
				panic(err)
			}

			if response.StatusCode != 200 {
				fmt.Print("Done finding building permits")
				break
			}

			defer response.Body.Close()

			sc := bufio.NewScanner(response.Body)
			lines := make([]string, 0)

			for sc.Scan() {
				lines = append(lines, sc.Text())
			}

			if err := sc.Err(); err != nil {
				log.Fatal(err)
			}

			//	Columns for reference:
			//	(0)Survey Date, (1)CSA Code, (2)CBSA Code, (3)MONCOV, (4)CBSA Name, (5)Bldgs, (6)1-unit Units, (7) Value, (8) Bldgs, (9)2-units Units, (10) Value,
			//	(11) Bldgs, (12)3-4 units Units, (13) Value, (14) Bldgs, (15)5+ units Units, (16) Value, (17) Bldgs, (18)1-unit rep Units, (19) Value, (20) Bldgs,
			//	(21)2-units rep Units, (22) Value, (23) Bldgs,  (24)3-4 units rep Units, (25) Value, (26) Bldgs, (27)5+ units rep Units, (28) Value

			for i, line := range lines {
				if i < 3 {
					continue
				}

				lineSplit := strings.Split(line, ",")

				date := lineSplit[0][4:] + "/01/" + lineSplit[0][:4]
				cbsacode := lineSplit[2]
				totalunits := ""

				if len(lineSplit) > 29 {
					panic("More than 29 fields found")
				}

				units_1, err := strconv.Atoi(lineSplit[6])
				if err != nil {
					panic(err)
				}
				units_2, err := strconv.Atoi(lineSplit[9])
				if err != nil {
					panic(err)
				}
				units_3to4, err := strconv.Atoi(lineSplit[12])
				if err != nil {
					panic(err)
				}
				units_5ormore, err := strconv.Atoi(lineSplit[15])
				if err != nil {
					panic(err)
				}

				totalUnitsSum := units_1 + units_2 + units_3to4 + units_5ormore
				totalunits = strconv.Itoa(totalUnitsSum)

				if oldcbsacode, ok := censusCbsaMap[cbsacode]; ok {
					cbsacode = censusCbsaMap[oldcbsacode]
				}

				alldates := buildingPermits[cbsacode].dates
				alldates = append(alldates, date)
				alltotalunits := buildingPermits[cbsacode].totalunits
				alltotalunits = append(alltotalunits, totalunits)

				buildingPermits[cbsacode] = BuildingPermits{dates: alldates, totalunits: alltotalunits}
			}

		}
	}

	return buildingPermits

}

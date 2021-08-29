package crime

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func ImportCrime() {
	url := "https://geoenrich.arcgis.com/arcgis/rest/services/World/geoenrichmentserver/GeoEnrichment/enrich"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("api_key", "key_from_environment_or_flag")
	q.Add("another_thing", "foo & bar")
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())
}

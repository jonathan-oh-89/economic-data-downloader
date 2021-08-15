package utils

import "fmt"

func FormatGeoFipsCode(geoLevel string, fips string) string {
	formattedFips := ""

	switch geoLevel {
	case "state":
		formattedFips = fmt.Sprintf("%02s", fips)
	case "county":
		formattedFips = fmt.Sprintf("%03s", fips)
	}

	return formattedFips
}

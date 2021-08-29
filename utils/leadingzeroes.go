package utils

import "fmt"

func FormatLeadingZeroes(formatstring string, fips string) string {
	formattedFips := ""

	switch formatstring {
	case "2":
		fallthrough
	case "state":
		formattedFips = fmt.Sprintf("%02s", fips)
	case "3":
		fallthrough
	case "county":
		formattedFips = fmt.Sprintf("%03s", fips)
	}

	return formattedFips
}

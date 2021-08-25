package utils

func CensusCbsaVersionMap() map[string]string {
	// file from: https://kenaninstitute.unc.edu/entrepreneurshipdata/wp-content/uploads/2019/04/lookup_table_for_2009_and_2013_codes.xls
	lines := ReadCSV("/files/cbsa_2009_and_2013_codes.csv")
	cbsaVersionMap := make(map[string]string)
	for i, line := range lines {
		if i < 1 {
			continue
		}

		if line[3] == "Removed" && line[4] != "" {
			cbsaVersionMap[line[1]] = line[4]
		}
	}

	return cbsaVersionMap

}

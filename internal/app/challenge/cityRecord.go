package challenge

import "strconv"

const (
	cityNameIndex            = 1
	cityASCIINameIndex       = 2
	cityAlernateNamesIndex   = 3
	cityLatitudeIndex        = 4
	cityLongitudeIndex       = 5
	cityCountryCodeIndex     = 8
	cityAdminLevel1CodeIndex = 10
)

type cityRecord struct {
	rawRecords []string
}

func (record *cityRecord) fetchCityName() string {
	if len(record.rawRecords) > cityNameIndex {
		return record.rawRecords[cityNameIndex]
	}
	return "-"
}

func (record *cityRecord) fetchCountryName() string {
	if len(record.rawRecords) > cityCountryCodeIndex {
		return record.rawRecords[cityCountryCodeIndex]
	}
	return "-"
}

func (record *cityRecord) fetchFirstAdministrationLevel() string {
	if len(record.rawRecords) > cityAdminLevel1CodeIndex {
		return record.rawRecords[cityAdminLevel1CodeIndex]
	}
	return "-"
}

func (record *cityRecord) fetchLatitude() float64 {
	if len(record.rawRecords) > cityLatitudeIndex {
		value, _ := strconv.ParseFloat(record.rawRecords[cityLatitudeIndex], 64)
		return value
	}
	return 0.0
}

func (record *cityRecord) fetchLongitude() float64 {
	if len(record.rawRecords) > cityLongitudeIndex {
		value, _ := strconv.ParseFloat(record.rawRecords[cityLongitudeIndex], 64)
		return value
	}
	return 0.0
}

func (record *cityRecord) fetchASCIIName() string {
	return record.rawRecords[cityASCIINameIndex]
}

func (record *cityRecord) fetchAlternateNames() string {
	return record.rawRecords[cityAlernateNamesIndex]
}

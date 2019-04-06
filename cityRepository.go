package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

const (
	name            = 1
	asciiname       = 2
	alternatenames  = 3
	countrycode     = 8
	adminLevel1Code = 10
)

type match struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Score     float32 `json:"score"`
}

type suggestions struct {
	Suggestions []match `json:"suggestions"`
}

type cityRepository struct {
	records [][]string
}

type cityRepositoryInterface interface {
	findSuggestionsFor(string) suggestions
}

// Creates a CityRepository using TSV file as the data source.
func createCityRepositoryFor(sourceTsvFilePath string) (cityRepository, error) {
	repository := cityRepository{}

	tsvFile, err := os.Open(sourceTsvFilePath)
	if err != nil {
		return repository, err
	}
	defer tsvFile.Close()

	reader := createReaderForTsvFileAndQuoteInValues(tsvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return repository, err
	}

	repository.records = records
	return repository, nil
}

func createReaderForTsvFileAndQuoteInValues(tsvFile *os.File) *csv.Reader {
	reader := csv.NewReader(tsvFile)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	return reader
}

func (repository *cityRepository) findSuggestionsFor(query string) suggestions {
	result := suggestions{Suggestions: []match{}}

	if query == "" {
		return result
	}

	queryName := strings.ToLower(query)

	for _, record := range repository.records {
		cityName := fetchCityNameOf(record)

		if strings.Contains(strings.ToLower(cityName), queryName) ||
			strings.Contains(strings.ToLower(record[asciiname]), queryName) ||
			strings.Contains(strings.ToLower(record[alternatenames]), queryName) {

			match := match{Name: fmt.Sprintf("%s, %s, %s", cityName, fetchFirstAdministrationLevelOf(record), fetchCountryNameOf(record))}
			result.Suggestions = append(result.Suggestions, match)
		}
	}

	return result
}

func fetchCityNameOf(record []string) string {
	if len(record) > name {
		return record[name]
	}
	return "-"
}

func fetchCountryNameOf(record []string) string {
	if len(record) > countrycode {
		return record[countrycode]
	}
	return "-"
}

func fetchFirstAdministrationLevelOf(record []string) string {
	if len(record) > adminLevel1Code {
		return record[adminLevel1Code]
	}
	return "-"
}

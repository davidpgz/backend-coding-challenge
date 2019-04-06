package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
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
	FindRankedSuggestionsFor(string) suggestions
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

func (repository *cityRepository) FindRankedSuggestionsFor(query string) suggestions {
	suggestions := repository.findSuggestionsFor(query)

	return suggestions
}

func (repository *cityRepository) findSuggestionsFor(query string) suggestions {
	result := suggestions{Suggestions: []match{}}

	if query == "" {
		return result
	}

	queryName := strings.ToLower(query)

	for _, record := range repository.records {
		matched, score := matchQueryName(record, queryName)
		if matched {
			cityName := fetchCityNameOf(record)
			match := match{
				Name:  fmt.Sprintf("%s, %s, %s", cityName, fetchFirstAdministrationLevelOf(record), fetchCountryNameOf(record)),
				Score: score,
			}
			result.Suggestions = append(result.Suggestions, match)
		}
	}

	return result
}

func matchQueryName(record []string, queryName string) (bool, float32) {
	matched := false
	score := 0.0

	if matched = strings.Contains(strings.ToLower(fetchCityNameOf(record)), queryName); matched {
		score = computeScoreFor(queryName, fetchCityNameOf(record))
	} else if matched = strings.Contains(strings.ToLower(record[asciiname]), queryName); matched {
		score = computeScoreFor(queryName, record[asciiname])
	} else if matched = strings.Contains(strings.ToLower(record[alternatenames]), queryName); matched {
		matchedWholeWord := findMatchingAlternateNameWholeWord(record[alternatenames], queryName)
		score = computeScoreFor(queryName, matchedWholeWord)
	}

	return matched, float32(score)
}

func computeScoreFor(queryName string, matchedWord string) float64 {
	return float64(utf8.RuneCountInString(queryName)) / float64(utf8.RuneCountInString(matchedWord))
}

func findMatchingAlternateNameWholeWord(recordAlternateNames string, queryName string) string {
	indexMatch := strings.Index(strings.ToLower(recordAlternateNames), queryName)
	alternateNames := []rune(recordAlternateNames)

	indexWordStart := findAlternateNameWordStartIndex(alternateNames, indexMatch)
	indexWordEnd := findAlternateNameWordEndIndex(alternateNames, indexMatch+utf8.RuneCountInString(queryName))
	return string(alternateNames[indexWordStart : indexWordEnd+1])
}

func findAlternateNameWordStartIndex(alternateNames []rune, searchStartIndex int) int {
	wordStartIndex := searchStartIndex
	for wordStartIndex > 0 {
		if alternateNames[wordStartIndex] == ',' {
			wordStartIndex++
			break
		}
		wordStartIndex--
	}
	return wordStartIndex
}

func findAlternateNameWordEndIndex(alternateNames []rune, searchStartIndex int) int {
	wordEndIndex := searchStartIndex
	alternateNamesLength := len(alternateNames)
	for wordEndIndex < alternateNamesLength {
		if alternateNames[wordEndIndex] == ',' {
			wordEndIndex--
			break
		}
		wordEndIndex++
	}
	return wordEndIndex
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

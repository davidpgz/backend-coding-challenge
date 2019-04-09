package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	name            = 1
	asciiname       = 2
	alternatenames  = 3
	latitude        = 4
	longitude       = 5
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
	FindRankedSuggestionsFor(cityQuery) suggestions
}

type cityQuery struct {
	name      string
	latitude  string
	longitude string
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

func (repository *cityRepository) FindRankedSuggestionsFor(query cityQuery) suggestions {
	suggestions := repository.findSuggestionsFor(query)
	// Sort suggestions by descending order
	sort.SliceStable(suggestions.Suggestions, func(i, j int) bool {
		return suggestions.Suggestions[i].Score > suggestions.Suggestions[j].Score
	})
	return suggestions
}

func (repository *cityRepository) findSuggestionsFor(query cityQuery) suggestions {
	result := suggestions{Suggestions: []match{}}

	if query.name == "" {
		return result
	}

	query.name = strings.ToLower(query.name)

	for _, record := range repository.records {
		matched, score := matchQueryName(record, query)
		if matched {
			cityName := fetchCityNameOf(record)
			match := match{
				Name:      fmt.Sprintf("%s, %s, %s", cityName, fetchFirstAdministrationLevelOf(record), fetchCountryNameOf(record)),
				Latitude:  fetchLatitude(record),
				Longitude: fetchLongitude(record),
				Score:     score,
			}
			result.Suggestions = append(result.Suggestions, match)
		}
	}

	return result
}

func matchQueryName(record []string, query cityQuery) (bool, float32) {
	matched := false
	score := 0.0

	if matched = strings.Contains(strings.ToLower(fetchCityNameOf(record)), query.name); matched {
		score = computeScoreFor(query, fetchCityNameOf(record), record)
	} else if matched = strings.Contains(strings.ToLower(record[asciiname]), query.name); matched {
		score = computeScoreFor(query, record[asciiname], record)
	} else if matched = strings.Contains(strings.ToLower(record[alternatenames]), query.name); matched {
		matchedWholeWord := findMatchingAlternateNameWholeWord(record[alternatenames], query.name)
		score = computeScoreFor(query, matchedWholeWord, record)
	}

	return matched, float32(score)
}

func computeScoreFor(query cityQuery, matchedWord string, record []string) float64 {
	matchingCharWeight := computeMatchingCharWeight(query.name, matchedWord)
	latitudeWeight := computeLatitudeScoreWeight(query, record)
	longitudeWeight := computeLongitudeScoreWeight(query, record)
	return matchingCharWeight * latitudeWeight * longitudeWeight
}

func computeMatchingCharWeight(queryName string, matchedWord string) float64 {
	return float64(utf8.RuneCountInString(queryName)) / float64(utf8.RuneCountInString(matchedWord))
}

func computeLatitudeScoreWeight(query cityQuery, record []string) float64 {
	const latitudeMaximumRange float64 = 180.0

	queryLatitude, err := strconv.ParseFloat(query.latitude, 64)
	if err == nil {
		recordLatitude := fetchLatitude(record)
		distanceRatio := math.Abs(queryLatitude-recordLatitude) / latitudeMaximumRange
		return 1 - distanceRatio
	}
	return 1
}

func computeLongitudeScoreWeight(query cityQuery, record []string) float64 {
	const longitudeMaximumRange float64 = 360.0

	queryLongitude, err := strconv.ParseFloat(query.longitude, 64)
	if err == nil {
		recordLongitude := fetchLongitude(record)
		distanceRatio := math.Abs(queryLongitude-recordLongitude) / longitudeMaximumRange
		return 1 - distanceRatio
	}
	return 1
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

func fetchLatitude(record []string) float64 {
	if len(record) > latitude {
		value, _ := strconv.ParseFloat(record[latitude], 64)
		return value
	}
	return 0.0
}

func fetchLongitude(record []string) float64 {
	if len(record) > longitude {
		value, _ := strconv.ParseFloat(record[longitude], 64)
		return value
	}
	return 0.0
}

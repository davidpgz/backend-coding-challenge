package challenge

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
	"github.com/ahmetb/go-linq"
)

type cityRepository struct {
	records []cityRecord
}

type cityRepositoryInterface interface {
	FindRankedSuggestionsFor(cityQuery) suggestions
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

	linq.From(records).SelectT(func(record []string) cityRecord {
		return cityRecord{record}
	}).ToSlice(&repository.records)

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
			cityName := record.fetchCityName()
			match := match{
				Name:      fmt.Sprintf("%s, %s, %s", cityName, record.fetchFirstAdministrationLevel(), record.fetchCountryName()),
				Latitude:  record.fetchLatitude(),
				Longitude: record.fetchLongitude(),
				Score:     score,
			}
			result.Suggestions = append(result.Suggestions, match)
		}
	}

	return result
}

func matchQueryName(record cityRecord, query cityQuery) (bool, float32) {
	matched := false
	score := 0.0

	if matched = strings.Contains(strings.ToLower(record.fetchCityName()), query.name); matched {
		score = computeScoreFor(query, record.fetchCityName(), record)
	} else if matched = strings.Contains(strings.ToLower(record.fetchASCIIName()), query.name); matched {
		score = computeScoreFor(query, record.fetchASCIIName(), record)
	} else if matched = strings.Contains(strings.ToLower(record.fetchAlternateNames()), query.name); matched {
		matchedWholeWord := findMatchingAlternateNameWholeWord(record.fetchAlternateNames(), query.name)
		score = computeScoreFor(query, matchedWholeWord, record)
	}

	return matched, float32(score)
}

func computeScoreFor(query cityQuery, matchedWord string, record cityRecord) float64 {
	matchingCharWeight := computeMatchingCharWeight(query.name, matchedWord)
	latitudeWeight := computeLatitudeScoreWeight(query, record)
	longitudeWeight := computeLongitudeScoreWeight(query, record)
	return matchingCharWeight * latitudeWeight * longitudeWeight
}

func computeMatchingCharWeight(queryName string, matchedWord string) float64 {
	return float64(utf8.RuneCountInString(queryName)) / float64(utf8.RuneCountInString(matchedWord))
}

func computeLatitudeScoreWeight(query cityQuery, record cityRecord) float64 {
	const latitudeMaximumRange float64 = 180.0

	queryLatitude, err := strconv.ParseFloat(query.latitude, 64)
	if err == nil {
		recordLatitude := record.fetchLatitude()
		distanceRatio := math.Abs(queryLatitude-recordLatitude) / latitudeMaximumRange
		return 1 - distanceRatio
	}
	return 1
}

func computeLongitudeScoreWeight(query cityQuery, record cityRecord) float64 {
	const longitudeMaximumRange float64 = 360.0

	queryLongitude, err := strconv.ParseFloat(query.longitude, 64)
	if err == nil {
		recordLongitude := record.fetchLongitude()
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

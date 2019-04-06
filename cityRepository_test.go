package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createCityRepository() cityRepository {
	cityRepository, _ := createCityRepositoryFor("./data/cities_canada-usa.tsv")
	return cityRepository
}

func TestFindSuggestionsForExactName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("Québec")
	assert.New(t).Contains(result.Suggestions[0].Name, "Québec")
}

func TestFindSuggestionsForInexistantName(t *testing.T) {
	cityRepository := createCityRepository()

	result := cityRepository.findSuggestionsFor("SomeRandomCityInTheMiddleOfNowhere")

	assert.New(t).Empty(result.Suggestions)
	assert.New(t).NotNil(result.Suggestions)
}

func TestFindSuggestionsForEmptyName(t *testing.T) {
	cityRepository := createCityRepository()

	result := cityRepository.findSuggestionsFor("")
	
	assert.New(t).Empty(result.Suggestions)
	assert.New(t).NotNil(result.Suggestions)
}

func TestFindSuggestionsForExactLowerCaseName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("québec")
	assert.New(t).Contains(result.Suggestions[0].Name, "Québec")
}

func TestFindSuggestionsForPartialName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("lond")
	assert.New(t).Contains(result.Suggestions[0].Name, "London")
}

func TestFindSuggestionsForAsciiName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("levis")
	assert.New(t).Contains(result.Suggestions[0].Name, "Lévis")
}

func TestFindSuggestionsForPartialAlternateName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("udad ti Que")
	assert.New(t).Contains(result.Suggestions[0].Name, "Québec")
}

func TestFindSuggestionsForAppendsAdmin1LevelAndCountryCodeToTheName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("Québec")
	assert.New(t).Contains(result.Suggestions[0].Name, "Québec, 10, CA")
}

func TestCreateCityRepositoryForParseTsvFileData(t *testing.T) {
	cityRepository := createCityRepository()
	assert.New(t).NotEmpty(cityRepository.records)
}

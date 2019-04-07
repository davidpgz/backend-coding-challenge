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

func TestFindSuggestionsFor_ShouldAppendsAdmin1LevelAndCountryCodeToTheName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("Québec")
	assert.New(t).Contains(result.Suggestions[0].Name, "Québec, 10, CA")
}

func TestCreateCityRepositoryFor_ShouldParseTsvFileData(t *testing.T) {
	cityRepository := createCityRepository()
	assert.New(t).NotEmpty(cityRepository.records)
}

func TestFindSuggestionsForExactName_ShouldHaveScoreOf1(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("Québec")
	assert.New(t).Equal(float32(1.0), result.Suggestions[0].Score)
}

func TestFindSuggestionsForPartialName_ShouldHaveScoreRatioMatchingNumberOfChar(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("qué")
	assert.New(t).Equal(float32(3.0/6.0), result.Suggestions[0].Score)
}

func TestFindSuggestionsForPartialAsciiName_ShouldHaveScoreRatioMatchingNumberOfChar(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("queb")
	assert.New(t).Equal(float32(4.0/6.0), result.Suggestions[0].Score)
}

func TestFindSuggestionsForPartialAlternateName_ShouldHaveScoreRatioMatchingNumberOfChar(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("udad ti Que")
	assert.New(t).Equal(float32(11.0/16.0), result.Suggestions[0].Score)
}

func TestFindSuggestionsForLatitude(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("Québec")
	assert.New(t).Equal(46.81228, result.Suggestions[0].Latitude)
}

func TestFindSuggestionsForLongitude(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor("Québec")
	assert.New(t).Equal(-71.21454, result.Suggestions[0].Longitude)
}

func TestFindRankedSuggestionsForPartialName_ShouldSortSuggestionsByDescendingOrder(t *testing.T) {
	cityRepository := createCityRepository()

	result := cityRepository.FindRankedSuggestionsFor("king")

	assert := assert.New(t)
	assert.True(result.Suggestions[0].Score >= result.Suggestions[1].Score)
	assert.True(result.Suggestions[1].Score >= result.Suggestions[2].Score)
	assert.True(result.Suggestions[2].Score >= result.Suggestions[3].Score)
}

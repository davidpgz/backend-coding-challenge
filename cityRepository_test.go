package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createCityRepository() cityRepository {
	cityRepository, _ := createCityRepositoryFor("./data/cities_canada-usa.tsv")
	return cityRepository
}

func queryFor(cityName string) cityQuery {
	return cityQuery{
		name: cityName,
	}
}

func TestFindSuggestionsForExactName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("Québec"))
	assert.New(t).Contains(result.Suggestions[0].Name, "Québec")
}

func TestFindSuggestionsForInexistantName(t *testing.T) {
	cityRepository := createCityRepository()

	result := cityRepository.findSuggestionsFor(queryFor("SomeRandomCityInTheMiddleOfNowhere"))

	assert.New(t).Empty(result.Suggestions)
	assert.New(t).NotNil(result.Suggestions)
}

func TestFindSuggestionsForEmptyName(t *testing.T) {
	cityRepository := createCityRepository()

	result := cityRepository.findSuggestionsFor(queryFor(""))

	assert.New(t).Empty(result.Suggestions)
	assert.New(t).NotNil(result.Suggestions)
}

func TestFindSuggestionsForExactLowerCaseName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("québec"))
	assert.New(t).Contains(result.Suggestions[0].Name, "Québec")
}

func TestFindSuggestionsForPartialName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("lond"))
	assert.New(t).Contains(result.Suggestions[0].Name, "London")
}

func TestFindSuggestionsForAsciiName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("levis"))
	assert.New(t).Contains(result.Suggestions[0].Name, "Lévis")
}

func TestFindSuggestionsForPartialAlternateName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("udad ti Que"))
	assert.New(t).Contains(result.Suggestions[0].Name, "Québec")
}

func TestFindSuggestionsFor_ShouldAppendsAdmin1LevelAndCountryCodeToTheName(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("Québec"))
	assert.New(t).Contains(result.Suggestions[0].Name, "Québec, 10, CA")
}

func TestCreateCityRepositoryFor_ShouldParseTsvFileData(t *testing.T) {
	cityRepository := createCityRepository()
	assert.New(t).NotEmpty(cityRepository.records)
}

func TestFindSuggestionsForExactName_ShouldHaveScoreOf1(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("Québec"))
	assert.New(t).Equal(float32(1.0), result.Suggestions[0].Score)
}

func TestFindSuggestionsForPartialName_ShouldHaveScoreRatioMatchingNumberOfChar(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("qué"))
	assert.New(t).Equal(float32(3.0/6.0), result.Suggestions[0].Score)
}

func TestFindSuggestionsForPartialAsciiName_ShouldHaveScoreRatioMatchingNumberOfChar(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("queb"))
	assert.New(t).Equal(float32(4.0/6.0), result.Suggestions[0].Score)
}

func TestFindSuggestionsForPartialAlternateName_ShouldHaveScoreRatioMatchingNumberOfChar(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("udad ti Que"))
	assert.New(t).Equal(float32(11.0/16.0), result.Suggestions[0].Score)
}

func TestFindSuggestionsForLatitude(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("Québec"))
	assert.New(t).Equal(46.81228, result.Suggestions[0].Latitude)
}

func TestFindSuggestionsForLongitude(t *testing.T) {
	cityRepository := createCityRepository()
	result := cityRepository.findSuggestionsFor(queryFor("Québec"))
	assert.New(t).Equal(-71.21454, result.Suggestions[0].Longitude)
}

func TestFindSuggestionsForLatitude_ShouldChangeScore(t *testing.T) {
	cityRepository := cityRepository{records: [][]string{[]string{"", "somecity", "", "", "45.0", "-90.0", "", "", "", "", ""}}}
	result := cityRepository.findSuggestionsFor(cityQuery{name: "city", latitude: "0.0"})
	assert.New(t).Equal(float32(4.0/8.0*(1.0-45.0/180.0)), result.Suggestions[0].Score)
}

func TestFindSuggestionsForLongitude_ShouldChangeScore(t *testing.T) {
	cityRepository := cityRepository{records: [][]string{[]string{"", "somecity", "", "", "45.0", "-90.0", "", "", "", "", ""}}}
	result := cityRepository.findSuggestionsFor(cityQuery{name: "city", longitude: "0.0"})
	assert.New(t).Equal(float32(4.0/8.0*(1.0-90.0/360.0)), result.Suggestions[0].Score)
}

func TestFindRankedSuggestionsForPartialName_ShouldSortSuggestionsByDescendingOrder(t *testing.T) {
	cityRepository := createCityRepository()

	result := cityRepository.FindRankedSuggestionsFor(queryFor("king"))

	assert := assert.New(t)
	assert.True(result.Suggestions[0].Score >= result.Suggestions[1].Score)
	assert.True(result.Suggestions[1].Score >= result.Suggestions[2].Score)
	assert.True(result.Suggestions[2].Score >= result.Suggestions[3].Score)
}

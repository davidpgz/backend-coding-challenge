package challenge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCityRepository struct {
	mockFindRankedSuggestionsFor func(cityQuery) suggestions
}

func (mock *mockCityRepository) FindRankedSuggestionsFor(query cityQuery) suggestions {
	if mock.mockFindRankedSuggestionsFor != nil {
		return mock.mockFindRankedSuggestionsFor(query)
	}
	return suggestions{}
}

func (app *App) serveGetSuggestions(path string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", path, nil)
	app.router.ServeHTTP(recorder, request)
	return recorder
}

func (app *App) mockFindRankedSuggestionsFor(funcMock func(cityQuery) suggestions) {
	app.cityRepository = &mockCityRepository{funcMock}
}

func newInitializedApp() *App {
	app := App{}
	app.Initialize("./../../../data/cities_canada-usa.tsv")
	return &app
}

func TestInitializeCanServeGetSuggestions(t *testing.T) {
	app := newInitializedApp()

	recorder := app.serveGetSuggestions(suggestionsPath)
	assert.New(t).Equal(http.StatusOK, recorder.Code)
}

func TestGetSuggestionsForCityName_ShouldParseQueryCityName(t *testing.T) {
	var wasCalledWithParsedQuery bool
	app := newInitializedApp()
	app.mockFindRankedSuggestionsFor(func(query cityQuery) suggestions {
		wasCalledWithParsedQuery = query.name == "something"
		return suggestions{}
	})

	app.serveGetSuggestions(suggestionsPath + "?q=something")

	assert.New(t).True(wasCalledWithParsedQuery)
}

func TestGetSuggestionsForCityNameAndLatitude_ShouldParseQueryCityNameAndLatitude(t *testing.T) {
	var wasCalledWithParsedQuery bool
	app := newInitializedApp()
	app.mockFindRankedSuggestionsFor(func(query cityQuery) suggestions {
		wasCalledWithParsedQuery = query.name == "something" && query.latitude == "12.345"
		return suggestions{}
	})

	app.serveGetSuggestions(suggestionsPath + "?q=something&latitude=12.345")

	assert.New(t).True(wasCalledWithParsedQuery)
}

func TestGetSuggestionsForCityNameAndLongitude_ShouldParseQueryCityNameAndLongitude(t *testing.T) {
	var wasCalledWithParsedQuery bool
	app := newInitializedApp()
	app.mockFindRankedSuggestionsFor(func(query cityQuery) suggestions {
		wasCalledWithParsedQuery = query.name == "something" && query.longitude == "12.345"
		return suggestions{}
	})

	app.serveGetSuggestions(suggestionsPath + "?q=something&longitude=12.345")

	assert.New(t).True(wasCalledWithParsedQuery)
}

func TestGetSuggestionsForCityName_ShouldServeJsonData(t *testing.T) {
	expectedSuggestions := suggestions{[]match{match{Name: "something", Longitude: -1.2, Latitude: 3.4, Score: 0.5}}}
	app := newInitializedApp()
	app.mockFindRankedSuggestionsFor(func(query cityQuery) suggestions {
		return expectedSuggestions
	})

	recorder := app.serveGetSuggestions(suggestionsPath + "?q=something")

	var result suggestions
	json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.New(t).Equal(expectedSuggestions, result)
}

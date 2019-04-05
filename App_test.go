package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyMockedObject struct {
	mockFindSuggestionsFor func(string) suggestions
}

func (mock *MyMockedObject) findSuggestionsFor(query string) suggestions {
	if mock.mockFindSuggestionsFor != nil {
		return mock.mockFindSuggestionsFor(query)
	}
	return suggestions{}
}

func (app *App) serveGetSuggestions(path string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", path, nil)
	app.router.ServeHTTP(recorder, request)
	return recorder
}

func TestInitializeCanServeGetSuggestions(t *testing.T) {
	assert := assert.New(t)

	app := App{}
	app.Initialize()
	recorder := app.serveGetSuggestions(suggestionsPath)

	assert.Equal(http.StatusOK, recorder.Code)
}

func TestGetSuggestionsParseRawQuery(t *testing.T) {
	assert := assert.New(t)

	var wasCalledWithParsedQuery bool
	app := App{}
	app.Initialize()
	app.cityRepository = &MyMockedObject{
		mockFindSuggestionsFor: func(query string) suggestions {
			wasCalledWithParsedQuery = query == "something"
			return suggestions{}
		},
	}

	app.serveGetSuggestions(suggestionsPath + "?q=something")

	assert.True(wasCalledWithParsedQuery)
}

func TestGetSuggestionsShouldServeJsonData(t *testing.T) {
	assert := assert.New(t)

	expectedSuggestions := suggestions{[]match{match{Name: "something", Longitude: -1.2, Latitude: 3.4, Score: 0.5}}}
	app := App{}
	app.Initialize()
	app.cityRepository = &MyMockedObject{
		mockFindSuggestionsFor: func(query string) suggestions {
			return expectedSuggestions
		},
	}

	recorder := app.serveGetSuggestions(suggestionsPath + "?q=something")

	var result suggestions
	json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.Equal(expectedSuggestions, result)
}

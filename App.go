package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	suggestionsPath = "/suggestions"
)

// App root struct used to configure the REST web API
type App struct {
	router         *gin.Engine
	cityRepository cityRepositoryInterface
}

// Initialize the App struct before the Run function is called
func (a *App) Initialize() error {
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET(suggestionsPath, func(context *gin.Context) {
		query := parseCityQuery(context.Request.URL.Query())
		suggestions := a.cityRepository.FindRankedSuggestionsFor(query)
		context.JSON(http.StatusOK, suggestions)
	})

	cityRepository, err := createCityRepositoryFor("./data/cities_canada-usa.tsv")
	if err != nil {
		return err
	}

	a.router = router
	a.cityRepository = &cityRepository
	return nil
}

func parseCityQuery(queryValues url.Values) cityQuery {
	return cityQuery{
		name:      queryValues.Get("q"),
		latitude:  queryValues.Get("latitude"),
		longitude: queryValues.Get("longitude"),
	}
}

// Run the web service.
// The App struct have to be initialized before calling this function.
func (a *App) Run(port string) {
	a.router.Run(":" + port)
}

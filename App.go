package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// App root struct used to configure the REST web API
type App struct {
	router *gin.Engine
}

// Initialize the App struct before the Run function is called
func (a *App) Initialize() {
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	a.router = router
}

// Launch the web service
func (a *App) Run(port string) {
	a.router.Run(":" + port)
}

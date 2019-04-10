package main

import (
	"log"
	"os"
	
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/heroku/go-getting-started/internal/app/challenge"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	app := challenge.App{}
	err := app.Initialize("./data/cities_canada-usa.tsv")
	if err != nil {
		log.Fatal(err)
	}
	app.Run(port)
}

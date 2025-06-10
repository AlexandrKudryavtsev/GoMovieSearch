// @title           AlexandrKudryavtsev/GoMovieSearch
// @version         1.0
// @description     movie search and autocomplete

// @BasePath  /api
// @schemes https http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	"github.com/AlexandrKudryavtsev/GoMovieSearch/config"
	"github.com/AlexandrKudryavtsev/GoMovieSearch/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("can't init config: %s", err)
	}

	app.Run(cfg)
}

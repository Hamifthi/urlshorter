package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"urlshortner/RandomGenerator"
	"urlshortner/config"
	"urlshortner/handlers"
	"urlshortner/postgres"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("must provide a valid env file to read")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s database=%s sslmode=%s",
		config.DatabaseHost,
		config.DatabaseUser,
		config.DatabasePass,
		config.DatabaseName,
		config.DatabaseSSLMode,
	))

	defer db.Close()

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	rg := RandomGenerator.RandomService{
		NumberOfDigits: 6,
	}

	urlService := &postgres.UrlService{
		DB:              db,
		RandomGenerator: &rg,
	}

	handler := &handlers.HTTPHandler{Service: urlService}

	// localhost:8080/create_short_link?url=https://google.com
	http.HandleFunc("/create_short_link", handler.ShortenLink)
	//localhost:8080/get_original_url?key=hwa54c
	http.HandleFunc("/get_original_url", handler.GetOriginalUrl)
	http.ListenAndServe(":8080", nil)
}

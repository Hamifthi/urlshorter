package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"urlshortner/RandomGenerator"
	"urlshortner/config"
	"urlshortner/handlers"
	"urlshortner/postgres"
)

type UrlShortnerSuite struct {
	suite.Suite

	db *sql.DB

	urlService *postgres.UrlService

	handler *handlers.HTTPHandler

	key string
}

func (suite *UrlShortnerSuite) SetupSuite() {
	err := godotenv.Load("../../.env")
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
	suite.db = db

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
	suite.urlService = urlService

	handler := &handlers.HTTPHandler{Service: urlService}
	suite.handler = handler
}

func (suite *UrlShortnerSuite) TearDownTest() {
	_, err := suite.urlService.DeleteShortLink(suite.key)
	if err != nil {
		suite.T().Fatal(err)
	}
	_, err = suite.urlService.DeleteAllMatchingUrls("testing.com")
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *UrlShortnerSuite) TestShortenLink() {
	req, err := http.NewRequest("POST", "/create_short_link?url=testing.com", nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.handler.ShortenLink)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(res, req)

	// Check the status code is what we expect.
	if status := res.Code; status != http.StatusCreated {
		suite.T().Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	r, _ := regexp.Compile("key=\\w{6}")
	found := r.MatchString(res.Body.String())

	if !found {
		suite.T().Errorf("handler returned unexpected response: got %v", res.Body.String())
	}

	match := r.FindString(res.Body.String())
	if match == "" {
		suite.T().Errorf("cannot find key in %v", res.Body.String())
	}
	suite.key = strings.Split(match, "=")[1]
}

func (suite *UrlShortnerSuite) TestGetOriginalUrl() {
	key, err := suite.urlService.CreateShortLink("testing.com")
	suite.key = key
	if err != nil {
		suite.T().Fatal(err)
	}
	url := fmt.Sprintf("/get_original_url?key=%s", key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.handler.GetOriginalUrl)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(res, req)

	// Check the status code is what we expect.
	if status := res.Code; status != http.StatusOK {
		suite.T().Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if !strings.Contains(res.Body.String(), "testing.com") {
		suite.T().Errorf("handler returned unexpected response: got %v", res.Body.String())
	}
}

func TestUrlShortnerTestSuite(t *testing.T) {
	suite.Run(t, new(UrlShortnerSuite))
}

package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"

	"urlshortner"
)

type UrlService struct {
	DB              *sql.DB
	RandomGenerator urlshortner.RandomService
}

func (url *UrlService) GetUrl(key string) (*urlshortner.Url, error) {
	var u urlshortner.Url
	row := url.DB.QueryRow(`SELECT key, url FROM urls WHERE key = $1`, key)
	if err := row.Scan(&u.KEY, &u.URL); err != nil {
		return nil, err
	}
	return &u, nil
}

func (url *UrlService) CreateShortLink(urlPart string) (string, error) {
	key := url.RandomGenerator.GenerateRandomString()
	_, err := url.DB.Exec(`INSERT into urls (key, url) VALUES ($1, $2)`, key, urlPart)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (url *UrlService) DeleteShortLink(key string) (string, error) {
	_, err := url.DB.Exec(`DELETE from urls WHERE key = $1`, key)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func (url *UrlService) DeleteAllMatchingUrls(proposedUrl string) (string, error) {
	_, err := url.DB.Exec(`DELETE from urls WHERE url = $1`, proposedUrl)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

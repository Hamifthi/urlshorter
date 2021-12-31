package config

import (
	"os"
)

func getEnv(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

var DatabaseHost = getEnv("DATABASE_HOST", "urlshortner_postgres")
var DatabaseUser = getEnv("DATABASE_USER", "postgres")
var DatabasePass = getEnv("DATABASE_PASS", "postgres")
var DatabaseName = getEnv("DATABASE_NAME", "urlshortner")
var DatabaseSSLMode = getEnv("DATABASE_SSL_MODE", "disable")
var REDIRECT_URL = getEnv("REDIRECT_URL", "localhost:8080/get_original_url?key=")

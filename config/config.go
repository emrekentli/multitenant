package config

import (
	"app/src/general/util/secret_reader"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

var (
	// DATABASE
	DbHost     = os.Getenv("DB_HOST")
	DbPort     = os.Getenv("DB_PORT")
	DbUsername = os.Getenv("DB_USERNAME")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbDatabase = os.Getenv("DB_DATABASE")
	DbSchema   = os.Getenv("DB_SCHEMA")
	DbSslMode  = os.Getenv("DB_SSL_MODE")

	// CACHE
	CacheHost     = os.Getenv("CACHE_HOST")
	CachePort     = os.Getenv("CACHE_PORT")
	CacheDb, _    = strconv.Atoi(os.Getenv("CACHE_DB"))
	CachePassword = os.Getenv("CACHE_PASSWORD")

	// SERVER
	AppName = os.Getenv("APP_NAME")
	Port    = ""
	PreFork = os.Getenv("PreFork") == "true"

	JwtSecretByte = []byte(secret_reader.ReadSecret(os.Getenv("jwt_secret")))

	AmpOptimizerUrl = os.Getenv("AMP_OPTIMIZER_URL")
)

const (
	Development = "development"
	Production  = "production"
)

var (
	IsProduction  = false
	IsDevelopment = false
)

func init() {
	port, err := strconv.Atoi(os.Getenv("Port"))
	if err != nil {
		Port = ":8080"
	} else {
		Port = ":" + strconv.Itoa(port)
	}

	env := os.Getenv("ENV")
	if env == "" || env == Development {
		IsDevelopment = true
	} else {
		if env == Production {
			IsProduction = true
		}
	}

}

package config

import (
	"os"
	"freelance_elite/database"
)

func InitDB() {
	var dbUser, dbPassword, dbName string
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	if os.Getenv("APP_ENV") == "test" {
		dbUser = os.Getenv("TEST_DB_USER")
		dbPassword = os.Getenv("TEST_DB_PASSWORD")
		dbName = os.Getenv("TEST_DB_NAME")
	} else {
		dbUser = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbName = os.Getenv("DB_NAME")
	}

	database.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)
}

package config

import (
	"os"
	"freelance_elite/database"
)

func InitDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	database.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)
}

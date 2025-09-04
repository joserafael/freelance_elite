package config

import (
	"os"
	"testing"
	"freelance_elite/database"
	"github.com/joho/godotenv"
)

// InitDB initializes the database connection based on environment
func InitDB() {
	var dbUser, dbPassword, dbName string
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	switch os.Getenv("APP_ENV") {
	case "test":
		dbUser = os.Getenv("TEST_DB_USER")
		dbPassword = os.Getenv("TEST_DB_PASSWORD")
		dbName = os.Getenv("TEST_DB_NAME")
	case "production":
		dbUser = os.Getenv("PROD_DB_USER")
		dbPassword = os.Getenv("PROD_DB_PASSWORD")
		dbName = os.Getenv("PROD_DB_NAME")
	default: // development
		dbUser = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbName = os.Getenv("DB_NAME")
	}

	database.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)
}

// SetupTestDB initializes the test database connection
// This function loads environment variables and connects to the test database
func SetupTestDB(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	
	// Load environment variables from .env file
	loadErr := godotenv.Load("../.env")
	if loadErr != nil {
		t.Fatal("Error loading .env file", loadErr)
	}

	// Use the existing InitDB function which already handles test environment
	InitDB()
}

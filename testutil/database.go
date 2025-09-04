// Package testutil provides shared test utilities that can be used across different test packages in the project.
package testutil

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"freelance_elite/database"
)

// SetupTestDB initializes the test database connection
// This function loads environment variables and connects to the test database
func SetupTestDB(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	
	// Load environment variables from .env file
	loadErr := godotenv.Load("../.env")
	if loadErr != nil {
		t.Fatal("Error loading .env file", loadErr)
	}

	// Get database configuration from environment variables
	var dbUser, dbPassword, dbHost, dbPort, dbName string
	dbUser = os.Getenv("TEST_DB_USER")
	dbPassword = os.Getenv("TEST_DB_PASSWORD")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("TEST_DB_NAME")
	
	// Initialize database connection
	database.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)
}
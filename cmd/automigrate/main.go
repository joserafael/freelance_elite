package main

import (
	"log"
	"os"

	"freelance_elite/database"
	"freelance_elite/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Initialize database connection
	database.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)

	// Auto-migrate the schema
	err = database.DB.AutoMigrate(&models.User{}, &models.BlacklistedToken{})
	if err != nil {
		log.Fatalf("failed to auto-migrate database: %v", err)
	}

	log.Println("Database migration completed successfully.")
}

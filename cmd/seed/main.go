package main

import (
	"log"
	"os"

	"freelance_elite/database"

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

	// Seed gender data
	seedGenders()

	// Seed countries data
	seedCountries()

	log.Println("Database seeding completed successfully.")
}
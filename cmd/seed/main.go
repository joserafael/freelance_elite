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

	// Seed gender data
	seedGenders()

	log.Println("Database seeding completed successfully.")
}

func seedGenders() {
	genders := []models.Gender{
		{
			Name:        "Male",
			Description: "Male gender",
			IsActive:    true,
		},
		{
			Name:        "Female",
			Description: "Female gender",
			IsActive:    true,
		},
		{
			Name:        "Non-binary",
			Description: "Non-binary gender identity",
			IsActive:    true,
		},
		{
			Name:        "Other",
			Description: "Other gender identity",
			IsActive:    true,
		},
		{
			Name:        "Prefer not to say",
			Description: "Prefer not to disclose gender",
			IsActive:    true,
		},
	}

	for _, gender := range genders {
		// Check if gender already exists
		var existingGender models.Gender
		result := database.DB.Where("name = ?", gender.Name).First(&existingGender)
		if result.Error != nil {
			// Gender doesn't exist, create it
			if err := database.DB.Create(&gender).Error; err != nil {
				log.Printf("Failed to create gender %s: %v", gender.Name, err)
			} else {
				log.Printf("Created gender: %s", gender.Name)
			}
		} else {
			log.Printf("Gender %s already exists, skipping", gender.Name)
		}
	}
}
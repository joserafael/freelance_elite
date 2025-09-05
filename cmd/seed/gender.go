package main

import (
	"log"

	"freelance_elite/database"
	"freelance_elite/models"
)

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

	// Insert genders, checking for duplicates
	for _, gender := range genders {
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
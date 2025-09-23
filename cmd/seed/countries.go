package main

import (
	"log"

	"freelance_elite/db"
	"freelance_elite/models"
)

func seedCountries() {
	countries := []models.Country{
		// North America
		{
			Name:       "United States",
			Code:       "USA",
			Region:     "Americas",
			Subregion:  "North America",
			Capital:    "Washington, D.C.",
			Population: 331900000,
			Area:       9833517,
			Currency:   "US Dollar",
			Language:   "English",
			IsActive:   true,
		},
		{
			Name:       "Canada",
			Code:       "CAN",
			Region:     "Americas",
			Subregion:  "North America",
			Capital:    "Ottawa",
			Population: 38000000,
			Area:       9984670,
			Currency:   "Canadian Dollar",
			Language:   "English, French",
			IsActive:   true,
		},
		{
			Name:       "Mexico",
			Code:       "MEX",
			Region:     "Americas",
			Subregion:  "North America",
			Capital:    "Mexico City",
			Population: 128900000,
			Area:       1964375,
			Currency:   "Mexican Peso",
			Language:   "Spanish",
			IsActive:   true,
		},
		// Central America
		{
			Name:       "Guatemala",
			Code:       "GTM",
			Region:     "Americas",
			Subregion:  "Central America",
			Capital:    "Guatemala City",
			Population: 17900000,
			Area:       108889,
			Currency:   "Guatemalan Quetzal",
			Language:   "Spanish",
			IsActive:   true,
		},
		{
			Name:       "Costa Rica",
			Code:       "CRI",
			Region:     "Americas",
			Subregion:  "Central America",
			Capital:    "San José",
			Population: 5100000,
			Area:       51100,
			Currency:   "Costa Rican Colón",
			Language:   "Spanish",
			IsActive:   true,
		},
		{
			Name:       "Panama",
			Code:       "PAN",
			Region:     "Americas",
			Subregion:  "Central America",
			Capital:    "Panama City",
			Population: 4300000,
			Area:       75417,
			Currency:   "Panamanian Balboa",
			Language:   "Spanish",
			IsActive:   true,
		},
		// South America
		{
			Name:       "Brazil",
			Code:       "BRA",
			Region:     "Americas",
			Subregion:  "South America",
			Capital:    "Brasília",
			Population: 215300000,
			Area:       8515767,
			Currency:   "Brazilian Real",
			Language:   "Portuguese",
			IsActive:   true,
		},
		{
			Name:       "Argentina",
			Code:       "ARG",
			Region:     "Americas",
			Subregion:  "South America",
			Capital:    "Buenos Aires",
			Population: 45400000,
			Area:       2780400,
			Currency:   "Argentine Peso",
			Language:   "Spanish",
			IsActive:   true,
		},
		{
			Name:       "Colombia",
			Code:       "COL",
			Region:     "Americas",
			Subregion:  "South America",
			Capital:    "Bogotá",
			Population: 51000000,
			Area:       1141748,
			Currency:   "Colombian Peso",
			Language:   "Spanish",
			IsActive:   true,
		},
		{
			Name:       "Peru",
			Code:       "PER",
			Region:     "Americas",
			Subregion:  "South America",
			Capital:    "Lima",
			Population: 33000000,
			Area:       1285216,
			Currency:   "Peruvian Sol",
			Language:   "Spanish",
			IsActive:   true,
		},
		{
			Name:       "Chile",
			Code:       "CHL",
			Region:     "Americas",
			Subregion:  "South America",
			Capital:    "Santiago",
			Population: 19100000,
			Area:       756102,
			Currency:   "Chilean Peso",
			Language:   "Spanish",
			IsActive:   true,
		},
		{
			Name:       "Ecuador",
			Code:       "ECU",
			Region:     "Americas",
			Subregion:  "South America",
			Capital:    "Quito",
			Population: 17600000,
			Area:       283561,
			Currency:   "US Dollar",
			Language:   "Spanish",
			IsActive:   true,
		},
		{
			Name:       "Venezuela",
			Code:       "VEN",
			Region:     "Americas",
			Subregion:  "South America",
			Capital:    "Caracas",
			Population: 28400000,
			Area:       916445,
			Currency:   "Venezuelan Bolívar",
			Language:   "Spanish",
			IsActive:   true,
		},
		// Caribbean
		{
			Name:       "Cuba",
			Code:       "CUB",
			Region:     "Americas",
			Subregion:  "Caribbean",
			Capital:    "Havana",
			Population: 11300000,
			Area:       109884,
			Currency:   "Cuban Peso",
			Language:   "Spanish",
			IsActive:   true,
		},
		{
			Name:       "Jamaica",
			Code:       "JAM",
			Region:     "Americas",
			Subregion:  "Caribbean",
			Capital:    "Kingston",
			Population: 2900000,
			Area:       10991,
			Currency:   "Jamaican Dollar",
			Language:   "English",
			IsActive:   true,
		},
	}

	// Insert countries, checking for duplicates
	for _, country := range countries {
		var existingCountry models.Country
		result := db.DB.Where("code = ?", country.Code).First(&existingCountry)
		if result.Error != nil {
			// Country doesn't exist, create it
			if err := db.DB.Create(&country).Error; err != nil {
				log.Printf("Failed to create country %s: %v", country.Name, err)
			} else {
				log.Printf("Created country: %s (%s)", country.Name, country.Code)
			}
		} else {
			log.Printf("Country %s (%s) already exists, skipping", country.Name, country.Code)
		}
	}
}
package models

import (
	"os"
	"testing"

	"freelance_elite/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CountryTestSuite struct {
	suite.Suite
}

func (s *CountryTestSuite) SetupSuite() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		s.T().Fatal("Error loading .env file")
	}

	// Get database configuration from environment
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Initialize database connection
	db.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)

	// Auto-migrate the schema
	db.DB.AutoMigrate(&Country{})
}

func (s *CountryTestSuite) TearDownSuite() {
	// Clean up test database after all tests are done
	db.DB.Exec("DELETE FROM countries")
}

func (s *CountryTestSuite) SetupTest() {
	// Clean the countries table before each test
	db.DB.Exec("DELETE FROM countries")
}

func (s *CountryTestSuite) TestCreateCountry() {
	country := Country{
		Name:     "United States",
		Code:     "US",
		IsActive: true,
	}

	result := db.DB.Create(&country)
	assert.NoError(s.T(), result.Error)
	assert.NotZero(s.T(), country.ID)
	assert.Equal(s.T(), "United States", country.Name)
	assert.Equal(s.T(), "US", country.Code)
	assert.True(s.T(), country.IsActive)
}

func (s *CountryTestSuite) TestCreateCountryWithDuplicateName() {
	// Create first country
	country1 := Country{
		Name:     "Canada",
		Code:     "CA",
		IsActive: true,
	}
	result1 := db.DB.Create(&country1)
	assert.NoError(s.T(), result1.Error)

	// Try to create second country with same name
	country2 := Country{
		Name:     "Canada",
		Code:     "CA2",
		IsActive: true,
	}
	result2 := db.DB.Create(&country2)
	assert.Error(s.T(), result2.Error)
}

func (s *CountryTestSuite) TestCreateCountryWithDuplicateCode() {
	// Create first country
	country1 := Country{
		Name:     "Mexico",
		Code:     "MX",
		IsActive: true,
	}
	result1 := db.DB.Create(&country1)
	assert.NoError(s.T(), result1.Error)

	// Try to create second country with same code
	country2 := Country{
		Name:     "Mexico 2",
		Code:     "MX",
		IsActive: true,
	}
	result2 := db.DB.Create(&country2)
	assert.Error(s.T(), result2.Error)
}

func (s *CountryTestSuite) TestFindCountry() {
	// Create a country first
	country := Country{
		Name:     "Brazil",
		Code:     "BR",
		IsActive: true,
	}
	db.DB.Create(&country)

	// Find the country
	var foundCountry Country
	result := db.DB.First(&foundCountry, country.ID)
	assert.NoError(s.T(), result.Error)
	assert.Equal(s.T(), country.ID, foundCountry.ID)
	assert.Equal(s.T(), "Brazil", foundCountry.Name)
	assert.Equal(s.T(), "BR", foundCountry.Code)
}

func (s *CountryTestSuite) TestUpdateCountry() {
	// Create a country first
	country := Country{
		Name:     "Argentina",
		Code:     "AR",
		IsActive: true,
	}
	db.DB.Create(&country)

	// Update the country
	country.Name = "Argentina Updated"
	country.IsActive = false
	result := db.DB.Save(&country)
	assert.NoError(s.T(), result.Error)

	// Verify the update
	var updatedCountry Country
	db.DB.First(&updatedCountry, country.ID)
	assert.Equal(s.T(), "Argentina Updated", updatedCountry.Name)
	assert.False(s.T(), updatedCountry.IsActive)
}

func (s *CountryTestSuite) TestDeleteCountry() {
	// Create a country first
	country := Country{
		Name:     "Temporary",
		Code:     "TMP",
		IsActive: true,
	}
	db.DB.Create(&country)

	// Delete the country
	result := db.DB.Delete(&country)
	assert.NoError(s.T(), result.Error)

	// Verify deletion
	var deletedCountry Country
	err := db.DB.First(&deletedCountry, country.ID).Error
	assert.Error(s.T(), err)
}

func (s *CountryTestSuite) TestFindCountries() {
	// Create multiple countries
	countries := []Country{
		{Name: "France", Code: "FR", IsActive: true},
		{Name: "Germany", Code: "DE", IsActive: true},
		{Name: "Italy", Code: "IT", IsActive: false},
	}

	for _, country := range countries {
		err := db.DB.Create(&country).Error
		assert.NoError(s.T(), err)
	}

	// Find all countries
	var foundCountries []Country
	err := db.DB.Find(&foundCountries).Error
	assert.NoError(s.T(), err)
	assert.Len(s.T(), foundCountries, 3)

	// Find active countries only
	var activeCountries []Country
	err = db.DB.Where("is_active = ?", true).Find(&activeCountries).Error
	assert.NoError(s.T(), err)
	assert.Len(s.T(), activeCountries, 2)
}

func (s *CountryTestSuite) TestTableName() {
	country := Country{}
	assert.Equal(s.T(), "countries", country.TableName())
}

func TestCountryTestSuite(t *testing.T) {
	suite.Run(t, new(CountryTestSuite))
}
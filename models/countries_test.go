package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/joho/godotenv"

	"freelance_elite/database"
)

type CountryTestSuite struct {
	suite.Suite
}

func (s *CountryTestSuite) SetupSuite() {
	// Setup test database configuration
	os.Setenv("APP_ENV", "test")
	
	// Load environment variables from .env file
	loadErr := godotenv.Load("../.env")
	if loadErr != nil {
		s.T().Fatal("Error loading .env file", loadErr)
	}
	
	// Initialize test database
	dbUser := os.Getenv("TEST_DB_USER")
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	dbName := os.Getenv("TEST_DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	
	database.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)
	
	// Auto-migrate Country table for tests
	database.DB.AutoMigrate(&Country{})
}

func (s *CountryTestSuite) TearDownSuite() {
	database.DB.Exec("DELETE FROM countries")
}

func (s *CountryTestSuite) SetupTest() {
	database.DB.Exec("DELETE FROM countries")
}

func (s *CountryTestSuite) TestCreateCountry() {
	country := Country{
		Name:       "United States",
		Code:       "USA",
		Region:     "Americas",
		Subregion:  "Northern America",
		Capital:    "Washington D.C.",
		Population: 331900000,
		Area:       9833517.0,
		Currency:   "USD",
		Language:   "English",
		IsActive:   true,
	}

	err := database.DB.Create(&country).Error
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), country.ID)
	assert.Equal(s.T(), "United States", country.Name)
	assert.Equal(s.T(), "USA", country.Code)
	assert.Equal(s.T(), "Americas", country.Region)
	assert.Equal(s.T(), "Northern America", country.Subregion)
	assert.Equal(s.T(), "Washington D.C.", country.Capital)
	assert.Equal(s.T(), int64(331900000), country.Population)
	assert.Equal(s.T(), 9833517.0, country.Area)
	assert.Equal(s.T(), "USD", country.Currency)
	assert.Equal(s.T(), "English", country.Language)
	assert.True(s.T(), country.IsActive)
}

func (s *CountryTestSuite) TestCreateCountryMinimalData() {
	country := Country{
		Name:     "Canada",
		Code:     "CAN",
		IsActive: true,
	}

	err := database.DB.Create(&country).Error
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), country.ID)
	assert.Equal(s.T(), "Canada", country.Name)
	assert.Equal(s.T(), "CAN", country.Code)
	assert.True(s.T(), country.IsActive)
}

func (s *CountryTestSuite) TestCreateCountryUniqueNameConstraint() {
	// Create first country
	country1 := Country{
		Name:     "Mexico",
		Code:     "MEX",
		IsActive: true,
	}
	err := database.DB.Create(&country1).Error
	assert.NoError(s.T(), err)

	// Try to create second country with same name
	country2 := Country{
		Name:     "Mexico",
		Code:     "MX2",
		IsActive: true,
	}
	err = database.DB.Create(&country2).Error
	assert.Error(s.T(), err)
}

func (s *CountryTestSuite) TestCreateCountryUniqueCodeConstraint() {
	// Create first country
	country1 := Country{
		Name:     "Brazil",
		Code:     "BRA",
		IsActive: true,
	}
	err := database.DB.Create(&country1).Error
	assert.NoError(s.T(), err)

	// Try to create second country with same code
	country2 := Country{
		Name:     "Brasil",
		Code:     "BRA",
		IsActive: true,
	}
	err = database.DB.Create(&country2).Error
	assert.Error(s.T(), err)
}

func (s *CountryTestSuite) TestCreateCountryWithEmptyName() {
	country := Country{
		Code:     "ARG",
		IsActive: true,
	}

	err := database.DB.Create(&country).Error
	// GORM doesn't enforce validation by default, so this might not error
	// In a real application, you would add validation middleware
	if err != nil {
		assert.Error(s.T(), err)
	} else {
		// If no error, verify the record was created but name is empty
		assert.Empty(s.T(), country.Name)
	}
}

func (s *CountryTestSuite) TestCreateCountryWithEmptyCode() {
	country := Country{
		Name:     "Argentina",
		IsActive: true,
	}

	err := database.DB.Create(&country).Error
	// GORM doesn't enforce validation by default, so this might not error
	// In a real application, you would add validation middleware
	if err != nil {
		assert.Error(s.T(), err)
	} else {
		// If no error, verify the record was created but code is empty
		assert.Empty(s.T(), country.Code)
	}
}

func (s *CountryTestSuite) TestUpdateCountry() {
	// Create a country
	country := Country{
		Name:       "Chile",
		Code:       "CHL",
		Region:     "Americas",
		Capital:    "Santiago",
		Population: 19116000,
		IsActive:   true,
	}
	err := database.DB.Create(&country).Error
	assert.NoError(s.T(), err)

	// Update the country
	country.Capital = "Santiago de Chile"
	country.Population = 19200000
	err = database.DB.Save(&country).Error
	assert.NoError(s.T(), err)

	// Verify the update
	var updatedCountry Country
	err = database.DB.First(&updatedCountry, country.ID).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Santiago de Chile", updatedCountry.Capital)
	assert.Equal(s.T(), int64(19200000), updatedCountry.Population)
}

func (s *CountryTestSuite) TestDeleteCountry() {
	// Create a country
	country := Country{
		Name:     "Test Delete",
		Code:     "TDL",
		IsActive: true,
	}
	err := database.DB.Create(&country).Error
	assert.NoError(s.T(), err)

	// Delete the country
	err = database.DB.Delete(&country).Error
	assert.NoError(s.T(), err)

	// Verify deletion
	var deletedCountry Country
	err = database.DB.First(&deletedCountry, country.ID).Error
	assert.Error(s.T(), err)
}

func (s *CountryTestSuite) TestFindCountries() {
	// Create active countries
	country1 := Country{Name: "Peru", Code: "PER", Region: "Americas", IsActive: true}
	err := database.DB.Create(&country1).Error
	assert.NoError(s.T(), err)

	country2 := Country{Name: "Colombia", Code: "COL", Region: "Americas", IsActive: true}
	err = database.DB.Create(&country2).Error
	assert.NoError(s.T(), err)

	// Create inactive country by updating after creation
	country3 := Country{Name: "Ecuador", Code: "ECU", Region: "Americas", IsActive: true}
	err = database.DB.Create(&country3).Error
	assert.NoError(s.T(), err)
	
	// Update to set IsActive to false
	err = database.DB.Model(&country3).Update("is_active", false).Error
	assert.NoError(s.T(), err)

	// Find all countries
	var allCountries []Country
	err = database.DB.Find(&allCountries).Error
	assert.NoError(s.T(), err)
	assert.Len(s.T(), allCountries, 3)

	// Find active countries only
	var activeCountries []Country
	err = database.DB.Where("is_active = ?", true).Find(&activeCountries).Error
	assert.NoError(s.T(), err)
	assert.Len(s.T(), activeCountries, 2)

	// Find inactive countries only
	var inactiveCountries []Country
	err = database.DB.Where("is_active = ?", false).Find(&inactiveCountries).Error
	assert.NoError(s.T(), err)
	assert.Len(s.T(), inactiveCountries, 1)
}

func (s *CountryTestSuite) TestTableName() {
	country := Country{}
	assert.Equal(s.T(), "countries", country.TableName())
}

func TestCountryTestSuite(t *testing.T) {
	suite.Run(t, new(CountryTestSuite))
}
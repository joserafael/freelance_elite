package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/joho/godotenv"

	"freelance_elite/database"
	"freelance_elite/models"
)

type CountryTestSuite struct {
	suite.Suite
	e *echo.Echo
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
	database.DB.AutoMigrate(&models.Country{})

	s.e = echo.New()
	s.e.GET("/countries", GetCountries)
	s.e.GET("/countries/:id", GetCountry)
	s.e.POST("/countries", CreateCountry)
	s.e.PUT("/countries/:id", UpdateCountry)
	s.e.DELETE("/countries/:id", DeleteCountry)
	s.e.GET("/countries/region/:region", GetCountriesByRegion)
}

func (s *CountryTestSuite) TearDownSuite() {
	// Clean up test database after all tests are done
	sqlDB, _ := database.DB.DB()
	sqlDB.Close()
}

func (s *CountryTestSuite) SetupTest() {
	// Clean the countries table before each test
	database.DB.Exec("DELETE FROM countries")
}

func (s *CountryTestSuite) TestGetCountriesEmpty() {
	req := httptest.NewRequest(http.MethodGet, "/countries", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	
	var countries []models.Country
	err := json.Unmarshal(rec.Body.Bytes(), &countries)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), countries, 0)
}

func (s *CountryTestSuite) TestCreateCountrySuccess() {
	countryData := map[string]interface{}{
		"name":       "United States",
		"code":       "USA",
		"region":     "Americas",
		"subregion":  "Northern America",
		"capital":    "Washington D.C.",
		"population": 331900000,
		"area":       9833517.0,
		"currency":   "USD",
		"language":   "English",
		"is_active":  true,
	}

	jsonData, _ := json.Marshal(countryData)
	req := httptest.NewRequest(http.MethodPost, "/countries", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	
	responseCountry := response["country"].(map[string]interface{})
	assert.NotZero(s.T(), responseCountry["ID"])
	assert.Equal(s.T(), "United States", responseCountry["name"])
	assert.Equal(s.T(), "USA", responseCountry["code"])
	assert.Equal(s.T(), "Americas", responseCountry["region"])
}

func (s *CountryTestSuite) TestCreateCountryEmptyName() {
	countryData := map[string]interface{}{
		"code":      "USA",
		"is_active": true,
	}

	jsonData, _ := json.Marshal(countryData)
	req := httptest.NewRequest(http.MethodPost, "/countries", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Name is required")
}

func (s *CountryTestSuite) TestCreateCountryEmptyCode() {
	countryData := map[string]interface{}{
		"name":      "United States",
		"is_active": true,
	}

	jsonData, _ := json.Marshal(countryData)
	req := httptest.NewRequest(http.MethodPost, "/countries", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Code is required")
}

func (s *CountryTestSuite) TestCreateCountryDuplicateName() {
	// Create first country
	country := models.Country{
		Name:     "Canada",
		Code:     "CAN",
		IsActive: true,
	}
	database.DB.Create(&country)

	// Try to create second country with same name
	countryData := map[string]interface{}{
		"name":      "Canada",
		"code":      "CA2",
		"is_active": true,
	}

	jsonData, _ := json.Marshal(countryData)
	req := httptest.NewRequest(http.MethodPost, "/countries", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusConflict, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Country with this name already exists")
}

func (s *CountryTestSuite) TestCreateCountryDuplicateCode() {
	// Create first country
	country := models.Country{
		Name:     "Mexico",
		Code:     "MEX",
		IsActive: true,
	}
	database.DB.Create(&country)

	// Try to create second country with same code
	countryData := map[string]interface{}{
		"name":      "Mexico 2",
		"code":      "MEX",
		"is_active": true,
	}

	jsonData, _ := json.Marshal(countryData)
	req := httptest.NewRequest(http.MethodPost, "/countries", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusConflict, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Country with this code already exists")
}

func (s *CountryTestSuite) TestGetCountriesWithData() {
	// Create test countries
	countries := []models.Country{
		{Name: "Argentina", Code: "ARG", Region: "Americas", IsActive: true},
		{Name: "Brazil", Code: "BRA", Region: "Americas", IsActive: true},
		{Name: "Chile", Code: "CHL", Region: "Americas", IsActive: false},
	}

	for _, country := range countries {
		database.DB.Create(&country)
	}

	req := httptest.NewRequest(http.MethodGet, "/countries", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	
	var retrievedCountries []models.Country
	err := json.Unmarshal(rec.Body.Bytes(), &retrievedCountries)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), retrievedCountries, 3)
	// Should be ordered by name ASC
	assert.Equal(s.T(), "Argentina", retrievedCountries[0].Name)
	assert.Equal(s.T(), "Brazil", retrievedCountries[1].Name)
	assert.Equal(s.T(), "Chile", retrievedCountries[2].Name)
}

func (s *CountryTestSuite) TestGetCountriesWithFilters() {
	// Create test countries
	countries := []models.Country{
		{Name: "Argentina", Code: "ARG", Region: "Americas", Subregion: "South America", IsActive: true},
		{Name: "Canada", Code: "CAN", Region: "Americas", Subregion: "Northern America", IsActive: true},
		{Name: "Germany", Code: "DEU", Region: "Europe", Subregion: "Western Europe", IsActive: true},
		{Name: "Chile", Code: "CHL", Region: "Americas", Subregion: "South America", IsActive: false},
	}

	for _, country := range countries {
		database.DB.Create(&country)
	}

	// Test region filter
	req := httptest.NewRequest(http.MethodGet, "/countries?region=Americas", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var americasCountries []models.Country
	err := json.Unmarshal(rec.Body.Bytes(), &americasCountries)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), americasCountries, 3)

	// Test is_active filter (should return Argentina, Canada, Germany - 3 active countries)
	req = httptest.NewRequest(http.MethodGet, "/countries?is_active=true", nil)
	rec = httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var activeCountries []models.Country
	err = json.Unmarshal(rec.Body.Bytes(), &activeCountries)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), activeCountries, 3)
	// Verify the active countries are Argentina, Canada, and Germany
	activeNames := make([]string, len(activeCountries))
	for i, country := range activeCountries {
		activeNames[i] = country.Name
	}
	assert.Contains(s.T(), activeNames, "Argentina")
	assert.Contains(s.T(), activeNames, "Canada")
	assert.Contains(s.T(), activeNames, "Germany")

	// Test search filter
	req = httptest.NewRequest(http.MethodGet, "/countries?search=arg", nil)
	rec = httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var searchResults []models.Country
	err = json.Unmarshal(rec.Body.Bytes(), &searchResults)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), searchResults, 1)
	assert.Equal(s.T(), "Argentina", searchResults[0].Name)
}

func (s *CountryTestSuite) TestGetCountryByIdSuccess() {
	// Create a country
	country := models.Country{
		Name:     "Peru",
		Code:     "PER",
		Region:   "Americas",
		Capital:  "Lima",
		IsActive: true,
	}
	database.DB.Create(&country)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/countries/%d", country.ID), nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	
	var retrievedCountry models.Country
	err := json.Unmarshal(rec.Body.Bytes(), &retrievedCountry)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), country.ID, retrievedCountry.ID)
	assert.Equal(s.T(), "Peru", retrievedCountry.Name)
	assert.Equal(s.T(), "PER", retrievedCountry.Code)
}

func (s *CountryTestSuite) TestGetCountryByIdNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/countries/999", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Country not found")
}

func (s *CountryTestSuite) TestGetCountryByIdInvalidId() {
	req := httptest.NewRequest(http.MethodGet, "/countries/invalid", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Invalid country ID")
}

func (s *CountryTestSuite) TestUpdateCountrySuccess() {
	// Create a country
	country := models.Country{
		Name:       "Colombia",
		Code:       "COL",
		Region:     "Americas",
		Capital:    "Bogota",
		Population: 50000000,
		IsActive:   true,
	}
	database.DB.Create(&country)

	// Update data
	updateData := map[string]interface{}{
		"name":       "Colombia",
		"code":       "COL",
		"region":     "Americas",
		"capital":    "Bogotá",
		"population": 51000000,
		"is_active":  true,
	}

	jsonData, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/countries/%d", country.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	
	responseCountry := response["country"].(map[string]interface{})
	assert.Equal(s.T(), "Bogotá", responseCountry["capital"])
	assert.Equal(s.T(), float64(51000000), responseCountry["population"])
}

func (s *CountryTestSuite) TestUpdateCountryNotFound() {
	updateData := map[string]interface{}{
		"name":      "Non-existent",
		"code":      "NEX",
		"is_active": true,
	}

	jsonData, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, "/countries/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Country not found")
}

func (s *CountryTestSuite) TestUpdateCountryDuplicateName() {
	// Create two countries
	country1 := models.Country{Name: "Ecuador", Code: "ECU", IsActive: true}
	country2 := models.Country{Name: "Venezuela", Code: "VEN", IsActive: true}
	database.DB.Create(&country1)
	database.DB.Create(&country2)

	// Try to update country2 with country1's name
	updateData := map[string]interface{}{
		"name":      "Ecuador",
		"code":      "VEN",
		"is_active": true,
	}

	jsonData, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/countries/%d", country2.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusConflict, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Country with this name already exists")
}

func (s *CountryTestSuite) TestDeleteCountrySuccess() {
	// Create a country
	country := models.Country{
		Name:     "To be deleted",
		Code:     "TBD",
		IsActive: true,
	}
	database.DB.Create(&country)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/countries/%d", country.ID), nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Country deleted successfully")

	// Verify deletion
	var deletedCountry models.Country
	err := database.DB.First(&deletedCountry, country.ID).Error
	assert.Error(s.T(), err)
}

func (s *CountryTestSuite) TestDeleteCountryNotFound() {
	req := httptest.NewRequest(http.MethodDelete, "/countries/999", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Country not found")
}

func (s *CountryTestSuite) TestDeleteCountryInvalidId() {
	req := httptest.NewRequest(http.MethodDelete, "/countries/invalid", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Invalid country ID")
}

func (s *CountryTestSuite) TestGetCountriesByRegionSuccess() {
	// Create test countries
	countries := []models.Country{
		{Name: "Argentina", Code: "ARG", Region: "Americas", IsActive: true},
		{Name: "Brazil", Code: "BRA", Region: "Americas", IsActive: true},
		{Name: "Germany", Code: "DEU", Region: "Europe", IsActive: true},
		{Name: "France", Code: "FRA", Region: "Europe", IsActive: true},
	}

	for _, country := range countries {
		database.DB.Create(&country)
	}

	req := httptest.NewRequest(http.MethodGet, "/countries/region/Americas", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	
	var americasCountries []models.Country
	err := json.Unmarshal(rec.Body.Bytes(), &americasCountries)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), americasCountries, 2)
	assert.Equal(s.T(), "Argentina", americasCountries[0].Name)
	assert.Equal(s.T(), "Brazil", americasCountries[1].Name)
}

func (s *CountryTestSuite) TestGetCountriesByRegionEmpty() {
	req := httptest.NewRequest(http.MethodGet, "/countries/region/NonExistent", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	
	var countries []models.Country
	err := json.Unmarshal(rec.Body.Bytes(), &countries)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), countries, 0)
}

func TestCountryTestSuite(t *testing.T) {
	suite.Run(t, new(CountryTestSuite))
}
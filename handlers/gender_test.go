package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"freelance_elite/db"
	"freelance_elite/models"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GenderTestSuite struct {
	suite.Suite
	e *echo.Echo
}

func (s *GenderTestSuite) SetupSuite() {
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
	
	db.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)
	
	// Auto-migrate Gender table for tests
	db.DB.AutoMigrate(&models.Gender{})

	s.e = echo.New()
	s.e.GET("/genders", GetGenders)
	s.e.GET("/genders/:id", GetGender)
	s.e.POST("/genders", CreateGender)
	s.e.PUT("/genders/:id", UpdateGender)
	s.e.DELETE("/genders/:id", DeleteGender)
}

func (s *GenderTestSuite) TearDownSuite() {
	// Clean up test database after all tests are done
	sqlDB, _ := db.DB.DB()
	sqlDB.Close()
}

func (s *GenderTestSuite) SetupTest() {
	// Clean the genders table before each test
	db.DB.Exec("DELETE FROM genders")
}

func (s *GenderTestSuite) TestGetGendersEmpty() {
	req := httptest.NewRequest(http.MethodGet, "/genders", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	
	var genders []models.Gender
	err := json.Unmarshal(rec.Body.Bytes(), &genders)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), genders, 0)
}

func (s *GenderTestSuite) TestCreateGenderSuccess() {
	gender := models.Gender{
		Name:        "Male",
		Description: "Male gender",
		IsActive:    true,
	}
	jsonGender, _ := json.Marshal(gender)

	req := httptest.NewRequest(http.MethodPost, "/genders", bytes.NewBuffer(jsonGender))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Gender created successfully")
	assert.Contains(s.T(), rec.Body.String(), "Male")

	// Verify gender is in the database
	var createdGender models.Gender
	err := db.DB.Where("name = ?", gender.Name).First(&createdGender).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), gender.Name, createdGender.Name)
	assert.Equal(s.T(), gender.Description, createdGender.Description)
	assert.Equal(s.T(), gender.IsActive, createdGender.IsActive)
}

func (s *GenderTestSuite) TestCreateGenderEmptyName() {
	gender := models.Gender{
		Name:        "",
		Description: "Empty name test",
		IsActive:    true,
	}
	jsonGender, _ := json.Marshal(gender)

	req := httptest.NewRequest(http.MethodPost, "/genders", bytes.NewBuffer(jsonGender))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Name is required")
}

func (s *GenderTestSuite) TestCreateGenderDuplicateName() {
	// Create first gender
	gender1 := models.Gender{
		Name:        "Female",
		Description: "Female gender",
		IsActive:    true,
	}
	jsonGender1, _ := json.Marshal(gender1)

	req := httptest.NewRequest(http.MethodPost, "/genders", bytes.NewBuffer(jsonGender1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	// Try to create another gender with the same name
	gender2 := models.Gender{
		Name:        "Female",
		Description: "Another female description",
		IsActive:    false,
	}
	jsonGender2, _ := json.Marshal(gender2)

	req = httptest.NewRequest(http.MethodPost, "/genders", bytes.NewBuffer(jsonGender2))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusConflict, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Gender with this name already exists")
}

func (s *GenderTestSuite) TestGetGendersWithData() {
	// Create test genders
	genders := []models.Gender{
		{Name: "Male", Description: "Male gender", IsActive: true},
		{Name: "Female", Description: "Female gender", IsActive: true},
		{Name: "Non-binary", Description: "Non-binary gender", IsActive: false},
	}

	for _, gender := range genders {
		db.DB.Create(&gender)
	}

	req := httptest.NewRequest(http.MethodGet, "/genders", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	
	var retrievedGenders []models.Gender
	err := json.Unmarshal(rec.Body.Bytes(), &retrievedGenders)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), retrievedGenders, 3)
}

func (s *GenderTestSuite) TestGetGenderByIdSuccess() {
	// Create a gender
	gender := models.Gender{
		Name:        "Other",
		Description: "Other gender",
		IsActive:    true,
	}
	db.DB.Create(&gender)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/genders/%d", gender.ID), nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	
	var retrievedGender models.Gender
	err := json.Unmarshal(rec.Body.Bytes(), &retrievedGender)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), gender.Name, retrievedGender.Name)
	assert.Equal(s.T(), gender.Description, retrievedGender.Description)
	assert.Equal(s.T(), gender.IsActive, retrievedGender.IsActive)
}

func (s *GenderTestSuite) TestGetGenderByIdNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/genders/999", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Gender not found")
}

func (s *GenderTestSuite) TestGetGenderByIdInvalidId() {
	req := httptest.NewRequest(http.MethodGet, "/genders/invalid", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Invalid gender ID")
}

func (s *GenderTestSuite) TestUpdateGenderSuccess() {
	// Create a gender
	gender := models.Gender{
		Name:        "Prefer not to say",
		Description: "Original description",
		IsActive:    true,
	}
	db.DB.Create(&gender)

	// Update the gender
	updateData := models.Gender{
		Name:        "Prefer not to disclose",
		Description: "Updated description",
		IsActive:    false,
	}
	jsonUpdate, _ := json.Marshal(updateData)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/genders/%d", gender.ID), bytes.NewBuffer(jsonUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Gender updated successfully")
	assert.Contains(s.T(), rec.Body.String(), "Prefer not to disclose")

	// Verify the update in database
	var updatedGender models.Gender
	db.DB.First(&updatedGender, gender.ID)
	assert.Equal(s.T(), "Prefer not to disclose", updatedGender.Name)
	assert.Equal(s.T(), "Updated description", updatedGender.Description)
	assert.False(s.T(), updatedGender.IsActive)
}

func (s *GenderTestSuite) TestUpdateGenderNotFound() {
	updateData := models.Gender{
		Name:        "Non-existent",
		Description: "This should fail",
		IsActive:    true,
	}
	jsonUpdate, _ := json.Marshal(updateData)

	req := httptest.NewRequest(http.MethodPut, "/genders/999", bytes.NewBuffer(jsonUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Gender not found")
}

func (s *GenderTestSuite) TestUpdateGenderDuplicateName() {
	// Create two genders
	gender1 := models.Gender{Name: "Male", Description: "Male gender", IsActive: true}
	gender2 := models.Gender{Name: "Female", Description: "Female gender", IsActive: true}
	db.DB.Create(&gender1)
	db.DB.Create(&gender2)

	// Try to update gender2 to have the same name as gender1
	updateData := models.Gender{Name: "Male"}
	jsonUpdate, _ := json.Marshal(updateData)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/genders/%d", gender2.ID), bytes.NewBuffer(jsonUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusConflict, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Gender with this name already exists")
}

func (s *GenderTestSuite) TestDeleteGenderSuccess() {
	// Create a gender
	gender := models.Gender{
		Name:        "To be deleted",
		Description: "This gender will be deleted",
		IsActive:    true,
	}
	db.DB.Create(&gender)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/genders/%d", gender.ID), nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Gender deleted successfully")

	// Verify deletion
	var deletedGender models.Gender
	err := db.DB.First(&deletedGender, gender.ID).Error
	assert.Error(s.T(), err)
}

func (s *GenderTestSuite) TestDeleteGenderNotFound() {
	req := httptest.NewRequest(http.MethodDelete, "/genders/999", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Gender not found")
}

func (s *GenderTestSuite) TestDeleteGenderInvalidId() {
	req := httptest.NewRequest(http.MethodDelete, "/genders/invalid", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Invalid gender ID")
}

func TestGenderTestSuite(t *testing.T) {
	suite.Run(t, new(GenderTestSuite))
}
package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/joho/godotenv"

	"freelance_elite/database"
)

type GenderTestSuite struct {
	suite.Suite
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
	
	database.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)
	
	// Auto-migrate Gender table for tests
	database.DB.AutoMigrate(&Gender{})
}

func (s *GenderTestSuite) TearDownSuite() {
	database.DB.Exec("DELETE FROM genders")
}

func (s *GenderTestSuite) SetupTest() {
	database.DB.Exec("DELETE FROM genders")
}

func (s *GenderTestSuite) TestCreateGender() {
	gender := Gender{
		Name:        "Male",
		Description: "Male gender",
		IsActive:    true,
	}

	err := database.DB.Create(&gender).Error
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), gender.ID)
	assert.Equal(s.T(), "Male", gender.Name)
	assert.Equal(s.T(), "Male gender", gender.Description)
	assert.True(s.T(), gender.IsActive)
}

func (s *GenderTestSuite) TestCreateGenderWithoutDescription() {
	gender := Gender{
		Name:     "Female",
		IsActive: true,
	}

	err := database.DB.Create(&gender).Error
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), gender.ID)
	assert.Equal(s.T(), "Female", gender.Name)
	assert.Empty(s.T(), gender.Description)
	assert.True(s.T(), gender.IsActive)
}

func (s *GenderTestSuite) TestCreateGenderUniqueConstraint() {
	// Create first gender
	gender1 := Gender{
		Name:        "Non-binary",
		Description: "Non-binary gender",
		IsActive:    true,
	}
	err := database.DB.Create(&gender1).Error
	assert.NoError(s.T(), err)

	// Try to create another gender with the same name
	gender2 := Gender{
		Name:        "Non-binary",
		Description: "Another non-binary description",
		IsActive:    false,
	}
	err = database.DB.Create(&gender2).Error
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "Duplicate")
}

func (s *GenderTestSuite) TestCreateGenderWithEmptyName() {
	gender := Gender{
		Name:        "",
		Description: "Empty name test",
		IsActive:    true,
	}

	err := database.DB.Create(&gender).Error
	assert.Error(s.T(), err)
}

func (s *GenderTestSuite) TestUpdateGender() {
	// Create a gender
	gender := Gender{
		Name:        "Other",
		Description: "Other gender",
		IsActive:    true,
	}
	err := database.DB.Create(&gender).Error
	assert.NoError(s.T(), err)

	// Update the gender
	gender.Description = "Updated description"
	gender.IsActive = false
	err = database.DB.Save(&gender).Error
	assert.NoError(s.T(), err)

	// Verify the update
	var updatedGender Gender
	err = database.DB.First(&updatedGender, gender.ID).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated description", updatedGender.Description)
	assert.False(s.T(), updatedGender.IsActive)
}

func (s *GenderTestSuite) TestDeleteGender() {
	// Create a gender
	gender := Gender{
		Name:        "Test Delete",
		Description: "Gender to be deleted",
		IsActive:    true,
	}
	err := database.DB.Create(&gender).Error
	assert.NoError(s.T(), err)

	// Delete the gender
	err = database.DB.Delete(&gender).Error
	assert.NoError(s.T(), err)

	// Verify deletion
	var deletedGender Gender
	err = database.DB.First(&deletedGender, gender.ID).Error
	assert.Error(s.T(), err)
}

func (s *GenderTestSuite) TestFindGenders() {
	// Create multiple genders
	genders := []Gender{
		{Name: "Male", Description: "Male gender", IsActive: true},
		{Name: "Female", Description: "Female gender", IsActive: true},
		{Name: "Non-binary", Description: "Non-binary gender", IsActive: false},
	}

	for _, gender := range genders {
		err := database.DB.Create(&gender).Error
		assert.NoError(s.T(), err)
	}

	// Find all genders
	var foundGenders []Gender
	err := database.DB.Find(&foundGenders).Error
	assert.NoError(s.T(), err)
	assert.Len(s.T(), foundGenders, 3)

	// Find active genders only
	var activeGenders []Gender
	err = database.DB.Where("is_active = ?", true).Find(&activeGenders).Error
	assert.NoError(s.T(), err)
	assert.Len(s.T(), activeGenders, 2)
}

func TestGenderTestSuite(t *testing.T) {
	suite.Run(t, new(GenderTestSuite))
}
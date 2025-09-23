package models

import (
	"os"
	"testing"

	"freelance_elite/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GenderTestSuite struct {
	suite.Suite
}

func (s *GenderTestSuite) SetupSuite() {
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
	db.DB.AutoMigrate(&Gender{})
}

func (s *GenderTestSuite) TearDownSuite() {
	// Clean up test database after all tests are done
	db.DB.Exec("DELETE FROM genders")
}

func (s *GenderTestSuite) SetupTest() {
	// Clean the genders table before each test
	db.DB.Exec("DELETE FROM genders")
}

func (s *GenderTestSuite) TestCreateGender() {
	gender := Gender{
		Name:        "Male",
		Description: "Male gender",
		IsActive:    true,
	}

	result := db.DB.Create(&gender)
	assert.NoError(s.T(), result.Error)
	assert.NotZero(s.T(), gender.ID)
	assert.Equal(s.T(), "Male", gender.Name)
	assert.Equal(s.T(), "Male gender", gender.Description)
	assert.True(s.T(), gender.IsActive)
}

func (s *GenderTestSuite) TestCreateGenderWithDuplicateName() {
	// Create first gender
	gender1 := Gender{
		Name:        "Female",
		Description: "Female gender",
		IsActive:    true,
	}
	result1 := db.DB.Create(&gender1)
	assert.NoError(s.T(), result1.Error)

	// Try to create second gender with same name
	gender2 := Gender{
		Name:        "Female",
		Description: "Another female gender",
		IsActive:    true,
	}
	result2 := db.DB.Create(&gender2)
	assert.Error(s.T(), result2.Error)
}

func (s *GenderTestSuite) TestFindGender() {
	// Create a gender first
	gender := Gender{
		Name:        "Non-binary",
		Description: "Non-binary gender",
		IsActive:    true,
	}
	db.DB.Create(&gender)

	// Find the gender
	var foundGender Gender
	result := db.DB.First(&foundGender, gender.ID)
	assert.NoError(s.T(), result.Error)
	assert.Equal(s.T(), gender.ID, foundGender.ID)
	assert.Equal(s.T(), "Non-binary", foundGender.Name)
}

func (s *GenderTestSuite) TestUpdateGender() {
	// Create a gender first
	gender := Gender{
		Name:        "Other",
		Description: "Other gender",
		IsActive:    true,
	}
	db.DB.Create(&gender)

	// Update the gender
	gender.Description = "Updated description"
	gender.IsActive = false
	result := db.DB.Save(&gender)
	assert.NoError(s.T(), result.Error)

	// Verify the update
	var updatedGender Gender
	db.DB.First(&updatedGender, gender.ID)
	assert.Equal(s.T(), "Updated description", updatedGender.Description)
	assert.False(s.T(), updatedGender.IsActive)
}

func (s *GenderTestSuite) TestDeleteGender() {
	// Create a gender first
	gender := Gender{
		Name:        "Temporary",
		Description: "Temporary gender",
		IsActive:    true,
	}
	db.DB.Create(&gender)

	// Delete the gender
	result := db.DB.Delete(&gender)
	assert.NoError(s.T(), result.Error)

	// Verify deletion
	var deletedGender Gender
	err := db.DB.First(&deletedGender, gender.ID).Error
	assert.Error(s.T(), err)
}

func (s *GenderTestSuite) TestFindGenders() {
	// Create multiple genders individually
	gender1 := Gender{Name: "Male", Description: "Male gender", IsActive: true}
	err := db.DB.Create(&gender1).Error
	assert.NoError(s.T(), err)

	gender2 := Gender{Name: "Female", Description: "Female gender", IsActive: true}
	err = db.DB.Create(&gender2).Error
	assert.NoError(s.T(), err)

	gender3 := Gender{Name: "Non-binary", Description: "Non-binary gender", IsActive: true}
	err = db.DB.Create(&gender3).Error
	assert.NoError(s.T(), err)

	// Update one gender to be inactive
	err = db.DB.Model(&gender3).Update("is_active", false).Error
	assert.NoError(s.T(), err)

	// Find all genders
	var foundGenders []Gender
	err = db.DB.Find(&foundGenders).Error
	assert.NoError(s.T(), err)
	assert.Len(s.T(), foundGenders, 3)

	// Find active genders only
	var activeGenders []Gender
	err = db.DB.Where("is_active = ?", true).Find(&activeGenders).Error
	assert.NoError(s.T(), err)
	assert.Len(s.T(), activeGenders, 2)
}

func TestGenderTestSuite(t *testing.T) {
	suite.Run(t, new(GenderTestSuite))
}
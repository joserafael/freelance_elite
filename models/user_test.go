package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/joho/godotenv"

	"freelance_elite/database"
)

type UserTestSuite struct {
	suite.Suite
}

func (s *UserTestSuite) SetupSuite() {
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
}

func (s *UserTestSuite) TearDownSuite() {
	database.DB.Exec("DELETE FROM users")
}

func (s *UserTestSuite) SetupTest() {
	database.DB.Exec("DELETE FROM users")
}

func (s *UserTestSuite) TestCreateUser() {
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	result := database.DB.Create(&user)
	assert.NoError(s.T(), result.Error)
	assert.NotZero(s.T(), user.ID)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
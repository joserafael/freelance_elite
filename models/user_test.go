package models

import (
	"os"
	"testing"

	"freelance_elite/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	
	db.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)
}

func (s *UserTestSuite) TearDownSuite() {
	db.DB.Exec("DELETE FROM users")
}

func (s *UserTestSuite) SetupTest() {
	db.DB.Exec("DELETE FROM users")
}

func (s *UserTestSuite) TestCreateUser() {
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	result := db.DB.Create(&user)
	assert.NoError(s.T(), result.Error)
	assert.NotZero(s.T(), user.ID)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
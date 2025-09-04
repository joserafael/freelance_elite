package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"freelance_elite/database"
	"freelance_elite/testutil"
)

type UserTestSuite struct {
	suite.Suite
}

func (s *UserTestSuite) SetupSuite() {
	// Setup test database configuration
	testutil.SetupTestDB(s.T())
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
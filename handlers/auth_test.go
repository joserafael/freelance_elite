package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"freelance_elite/database"
	"freelance_elite/models"
)

type AuthTestSuite struct {
	suite.Suite
	e *echo.Echo
}

func (s *AuthTestSuite) SetupSuite() {
	// Load test environment variables
	os.Setenv("APP_ENV", "test")
	envMap, err := godotenv.Read("../.env.test")
	if err != nil {
		s.Fail("Error reading .env.test file", err)
	}

	dbUser := envMap["DB_USER"]
	dbPassword := envMap["DB_PASSWORD"]
	dbHost := envMap["DB_HOST"]
	dbPort := envMap["DB_PORT"]
	dbName := envMap["DB_NAME"]

	log.Printf("Test DB_USER: %s", dbUser)
	log.Printf("Test DB_PASSWORD: %s", dbPassword)
	log.Printf("Test DB_HOST: %s", dbHost)
	log.Printf("Test DB_PORT: %s", dbPort)
	log.Printf("Test DB_NAME: %s", dbName)

	// Initialize the test database
	database.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)

	// Create a new Echo instance for testing
	s.e = echo.New()
	s.e.POST("/register", Register)
	s.e.POST("/login", Login)
	s.e.POST("/logout", Logout)

	p := s.e.Group("/profile")
	p.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))
	p.Use(CheckBlacklist)
	p.GET("", Profile)
}

func (s *AuthTestSuite) TearDownSuite() {
	// Clean up test database after all tests are done
	// Close the database connection
	sqlDB, _ := database.DB.DB()
	sqlDB.Close()
}

func (s *AuthTestSuite) SetupTest() {
	// Clean the users table before each test
	database.DB.Exec("TRUNCATE TABLE users")
}

func (s *AuthTestSuite) TestRegisterSuccess() {
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonUser, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "User created successfully")

	// Verify user is in the database
	var registeredUser models.User
	err := database.DB.Where("email = ?", user.Email).First(&registeredUser).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Email, registeredUser.Email)
}

func (s *AuthTestSuite) TestRegisterDuplicateEmail() {
	// Register once
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonUser, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	// Register again with same email
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusConflict, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Username or email already exists")
}

func (s *AuthTestSuite) TestLoginSuccess() {
	// Register a user first
	user := models.User{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: "password123",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	database.DB.Create(&user)

	// Attempt to login
	loginPayload := models.User{
		Email:    "login@example.com",
		Password: "password123",
	}
	jsonLogin, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "token")
}

func (s *AuthTestSuite) TestLoginInvalidCredentials() {
	// Attempt to login without registering
	loginPayload := models.User{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}
	jsonLogin, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Invalid email or password")
}

func (s *AuthTestSuite) TestLogoutSuccess() {
	// Register a user
	user := models.User{
		Username: "logoutuser",
		Email:    "logout@example.com",
		Password: "password123",
	}
	jsonUser, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	// Login the user to get a token
	loginPayload := models.User{
		Email:    "logout@example.com",
		Password: "password123",
	}
	jsonLogin, _ := json.Marshal(loginPayload)

	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code)

	var loginResponse map[string]string
	json.Unmarshal(rec.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]
	assert.NotEmpty(s.T(), token)

	// Logout the user
	req = httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Logged out successfully")

	// Try to access a protected route with the blacklisted token
	req = httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "Token has been revoked")
}


func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
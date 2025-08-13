package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"freelance_elite/database"
	"freelance_elite/models"
)

func Register(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}
	user.Password = string(hashedPassword)

	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Username or email already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func Login(c echo.Context) error {
	payload := new(models.User)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func Logout(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header missing"})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Parse the token to get claims and expiration time
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expiration not found"})
	}
	expiresAt := time.Unix(int64(exp), 0)

	blacklistedToken := models.BlacklistedToken{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Create(&blacklistedToken).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to blacklist token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func CheckBlacklist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return next(c) // Let echojwt.JWT handle missing token
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		var blacklistedToken models.BlacklistedToken
		if err := database.DB.Where("token = ?", tokenString).First(&blacklistedToken).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return next(c) // Token not blacklisted, proceed
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check blacklist"})
		}

		// Token found in blacklist
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token has been revoked"})
	}
}


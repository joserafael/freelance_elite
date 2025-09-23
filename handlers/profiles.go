package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"freelance_elite/database"
	"freelance_elite/models"

	"github.com/labstack/echo/v4"
)

// CreateProfile creates a new profile
func CreateProfile(c echo.Context) error {
	profile := new(models.Profile)
	if err := c.Bind(profile); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate required fields
	if strings.TrimSpace(profile.Name) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
	}
	if strings.TrimSpace(profile.LastName) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Last name is required"})
	}
	if profile.DateBirth.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Date of birth is required"})
	}
	if profile.UserID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	if profile.GenderID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Gender ID is required"})
	}
	if profile.CountryID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Country ID is required"})
	}

	// Check if user exists
	var user models.User
	if err := database.DB.First(&user, profile.UserID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not found"})
	}

	// Check if gender exists
	var gender models.Gender
	if err := database.DB.First(&gender, profile.GenderID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Gender not found"})
	}

	// Check if country exists
	var country models.Country
	if err := database.DB.First(&country, profile.CountryID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Country not found"})
	}

	// Create the profile
	if err := database.DB.Create(&profile).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create profile"})
	}

	// Load relationships
	database.DB.Preload("User").Preload("Gender").Preload("Country").First(&profile, profile.ID)

	return c.JSON(http.StatusCreated, profile)
}

// GetProfiles retrieves all profiles with optional filters
func GetProfiles(c echo.Context) error {
	var profiles []models.Profile
	query := database.DB.Preload("User").Preload("Gender").Preload("Country")

	// Filter by gender
	if genderID := c.QueryParam("gender_id"); genderID != "" {
		query = query.Where("gender_id = ?", genderID)
	}

	// Filter by country
	if countryID := c.QueryParam("country_id"); countryID != "" {
		query = query.Where("country_id = ?", countryID)
	}

	// Search by name or last name
	if search := c.QueryParam("search"); search != "" {
		query = query.Where("name LIKE ? OR last_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Filter by age range (calculated from date_birth)
	if minAge := c.QueryParam("min_age"); minAge != "" {
		if age, err := strconv.Atoi(minAge); err == nil {
			maxDate := time.Now().AddDate(-age, 0, 0)
			query = query.Where("date_birth <= ?", maxDate)
		}
	}

	if maxAge := c.QueryParam("max_age"); maxAge != "" {
		if age, err := strconv.Atoi(maxAge); err == nil {
			minDate := time.Now().AddDate(-age-1, 0, 0)
			query = query.Where("date_birth > ?", minDate)
		}
	}

	if err := query.Find(&profiles).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve profiles"})
	}

	return c.JSON(http.StatusOK, profiles)
}

// GetProfileByID retrieves a profile by ID
func GetProfileByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Profile ID is required"})
	}

	// Validate ID format
	if _, err := strconv.Atoi(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid profile ID format"})
	}

	var profile models.Profile
	if err := database.DB.Preload("User").Preload("Gender").Preload("Country").First(&profile, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
	}

	return c.JSON(http.StatusOK, profile)
}

// UpdateProfile updates an existing profile
func UpdateProfile(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Profile ID is required"})
	}

	// Validate ID format
	if _, err := strconv.Atoi(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid profile ID format"})
	}

	// Check if profile exists
	var existingProfile models.Profile
	if err := database.DB.First(&existingProfile, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
	}

	// Bind the updated data
	updatedProfile := new(models.Profile)
	if err := c.Bind(updatedProfile); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate required fields
	if strings.TrimSpace(updatedProfile.Name) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
	}
	if strings.TrimSpace(updatedProfile.LastName) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Last name is required"})
	}
	if updatedProfile.DateBirth.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Date of birth is required"})
	}
	if updatedProfile.UserID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	if updatedProfile.GenderID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Gender ID is required"})
	}
	if updatedProfile.CountryID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Country ID is required"})
	}

	// Check if user exists
	var user models.User
	if err := database.DB.First(&user, updatedProfile.UserID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not found"})
	}

	// Check if gender exists
	var gender models.Gender
	if err := database.DB.First(&gender, updatedProfile.GenderID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Gender not found"})
	}

	// Check if country exists
	var country models.Country
	if err := database.DB.First(&country, updatedProfile.CountryID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Country not found"})
	}

	// Update the profile
	existingProfile.Name = updatedProfile.Name
	existingProfile.LastName = updatedProfile.LastName
	existingProfile.DateBirth = updatedProfile.DateBirth
	existingProfile.About = updatedProfile.About
	existingProfile.UserID = updatedProfile.UserID
	existingProfile.GenderID = updatedProfile.GenderID
	existingProfile.CountryID = updatedProfile.CountryID

	if err := database.DB.Save(&existingProfile).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update profile"})
	}

	// Load relationships
	database.DB.Preload("User").Preload("Gender").Preload("Country").First(&existingProfile, existingProfile.ID)

	return c.JSON(http.StatusOK, existingProfile)
}

// DeleteProfile deletes a profile by ID
func DeleteProfile(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Profile ID is required"})
	}

	// Validate ID format
	if _, err := strconv.Atoi(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid profile ID format"})
	}

	// Check if profile exists
	var profile models.Profile
	if err := database.DB.First(&profile, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
	}

	// Delete the profile (soft delete)
	if err := database.DB.Delete(&profile).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete profile"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile deleted successfully"})
}
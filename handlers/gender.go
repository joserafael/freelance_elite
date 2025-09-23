package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"freelance_elite/db"
	"freelance_elite/models"
)

// GetGenders retrieves all genders
func GetGenders(c echo.Context) error {
	var genders []models.Gender
	if err := db.DB.Find(&genders).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve genders"})
	}
	return c.JSON(http.StatusOK, genders)
}

// GetGender retrieves a specific gender by ID
func GetGender(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid gender ID"})
	}

	var gender models.Gender
	if err := db.DB.First(&gender, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Gender not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve gender"})
	}
	return c.JSON(http.StatusOK, gender)
}

// CreateGender creates a new gender
func CreateGender(c echo.Context) error {
	gender := new(models.Gender)
	if err := c.Bind(gender); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Validate required fields
	if gender.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
	}

	// Check if gender with same name already exists
	var existingGender models.Gender
	if err := db.DB.Where("name = ?", gender.Name).First(&existingGender).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Gender with this name already exists"})
	}

	if err := db.DB.Create(&gender).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create gender"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Gender created successfully",
		"gender":  gender,
	})
}

// UpdateGender updates an existing gender
func UpdateGender(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid gender ID"})
	}

	var gender models.Gender
	if err := db.DB.First(&gender, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Gender not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve gender"})
	}

	updateData := new(models.Gender)
	if err := c.Bind(updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Check if another gender with the same name exists (excluding current gender)
	if updateData.Name != "" && updateData.Name != gender.Name {
		var existingGender models.Gender
		if err := db.DB.Where("name = ? AND id != ?", updateData.Name, id).First(&existingGender).Error; err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Gender with this name already exists"})
		}
	}

	// Update fields
	if updateData.Name != "" {
		gender.Name = updateData.Name
	}
	if updateData.Description != "" {
		gender.Description = updateData.Description
	}
	// Update IsActive field (including false values)
	gender.IsActive = updateData.IsActive

	if err := db.DB.Save(&gender).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update gender"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Gender updated successfully",
		"gender":  gender,
	})
}

// DeleteGender deletes a gender
func DeleteGender(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid gender ID"})
	}

	var gender models.Gender
	if err := db.DB.First(&gender, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Gender not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve gender"})
	}

	if err := db.DB.Delete(&gender).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete gender"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Gender deleted successfully"})
}
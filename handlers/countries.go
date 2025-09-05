package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"freelance_elite/database"
	"freelance_elite/models"
)

// GetCountries retrieves all countries
func GetCountries(c echo.Context) error {
	var countries []models.Country

	// Get query parameters for filtering
	region := c.QueryParam("region")
	subregion := c.QueryParam("subregion")
	isActive := c.QueryParam("is_active")
	search := c.QueryParam("search")

	// Build query
	query := database.DB

	// Apply filters
	if region != "" {
		query = query.Where("region = ?", region)
	}
	if subregion != "" {
		query = query.Where("subregion LIKE ?", "%"+subregion+"%")
	}
	if isActive != "" {
		if isActive == "true" {
			query = query.Where("is_active = ?", true)
		} else if isActive == "false" {
			query = query.Where("is_active = ?", false)
		}
	}
	if search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR capital LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Execute query with ordering
	if err := query.Order("name ASC").Find(&countries).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve countries",
		})
	}

	return c.JSON(http.StatusOK, countries)
}

// GetCountry retrieves a single country by ID
func GetCountry(c echo.Context) error {
	id := c.Param("id")
	countryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid country ID",
		})
	}

	var country models.Country
	if err := database.DB.First(&country, countryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Country not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve country",
		})
	}

	return c.JSON(http.StatusOK, country)
}

// CreateCountry creates a new country
func CreateCountry(c echo.Context) error {
	var country models.Country

	if err := c.Bind(&country); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if strings.TrimSpace(country.Name) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Name is required",
		})
	}
	if strings.TrimSpace(country.Code) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Code is required",
		})
	}

	// Convert code to uppercase for consistency
	country.Code = strings.ToUpper(strings.TrimSpace(country.Code))
	country.Name = strings.TrimSpace(country.Name)

	// Check if country with same name or code already exists
	var existingCountry models.Country
	if err := database.DB.Where("name = ? OR code = ?", country.Name, country.Code).First(&existingCountry).Error; err == nil {
		if existingCountry.Name == country.Name {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Country with this name already exists",
			})
		}
		if existingCountry.Code == country.Code {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Country with this code already exists",
			})
		}
	}

	// Create the country
	if err := database.DB.Create(&country).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create country",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Country created successfully",
		"country": country,
	})
}

// UpdateCountry updates an existing country
func UpdateCountry(c echo.Context) error {
	id := c.Param("id")
	countryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid country ID",
		})
	}

	// Check if country exists
	var existingCountry models.Country
	if err := database.DB.First(&existingCountry, countryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Country not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve country",
		})
	}

	// Bind the updated data
	var updateData models.Country
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate and clean data if provided
	if updateData.Name != "" {
		updateData.Name = strings.TrimSpace(updateData.Name)
		if updateData.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Name cannot be empty",
			})
		}
	}
	if updateData.Code != "" {
		updateData.Code = strings.ToUpper(strings.TrimSpace(updateData.Code))
		if updateData.Code == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Code cannot be empty",
			})
		}
	}

	// Check for duplicates if name or code is being updated
	if updateData.Name != "" && updateData.Name != existingCountry.Name {
		var duplicateCountry models.Country
		if err := database.DB.Where("name = ? AND id != ?", updateData.Name, countryID).First(&duplicateCountry).Error; err == nil {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Country with this name already exists",
			})
		}
	}
	if updateData.Code != "" && updateData.Code != existingCountry.Code {
		var duplicateCountry models.Country
		if err := database.DB.Where("code = ? AND id != ?", updateData.Code, countryID).First(&duplicateCountry).Error; err == nil {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Country with this code already exists",
			})
		}
	}

	// Update the country
	if err := database.DB.Model(&existingCountry).Updates(updateData).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update country",
		})
	}

	// Fetch the updated country
	var updatedCountry models.Country
	database.DB.First(&updatedCountry, countryID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Country updated successfully",
		"country": updatedCountry,
	})
}

// DeleteCountry deletes a country
func DeleteCountry(c echo.Context) error {
	id := c.Param("id")
	countryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid country ID",
		})
	}

	// Check if country exists
	var country models.Country
	if err := database.DB.First(&country, countryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Country not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve country",
		})
	}

	// Delete the country
	if err := database.DB.Delete(&country).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete country",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Country deleted successfully",
	})
}

// GetCountriesByRegion retrieves countries filtered by region
func GetCountriesByRegion(c echo.Context) error {
	region := c.Param("region")
	if region == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Region parameter is required",
		})
	}

	var countries []models.Country
	if err := database.DB.Where("region LIKE ?", "%"+region+"%").Order("name ASC").Find(&countries).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve countries",
		})
	}

	return c.JSON(http.StatusOK, countries)
}
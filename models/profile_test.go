package models

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProfileTableName(t *testing.T) {
	profile := Profile{}
	assert.Equal(t, "profiles", profile.TableName())
}

func TestProfileValidation(t *testing.T) {
	// Test valid profile
	validProfile := Profile{
		Name:      "John",
		LastName:  "Doe",
		DateBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		About:     "Software developer",
		UserID:    1,
		GenderID:  1,
		CountryID: 1,
	}

	// Test that all fields are properly set
	assert.Equal(t, "John", validProfile.Name)
	assert.Equal(t, "Doe", validProfile.LastName)
	assert.Equal(t, time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), validProfile.DateBirth)
	assert.Equal(t, "Software developer", validProfile.About)
	assert.Equal(t, 1, validProfile.UserID)
	assert.Equal(t, uint(1), validProfile.GenderID)
	assert.Equal(t, uint(1), validProfile.CountryID)
}

func TestProfileJSONTags(t *testing.T) {
	// This test ensures that the JSON tags are properly set
	// by checking the struct field tags
	profile := Profile{}
	profileType := reflect.TypeOf(profile)

	// Check Name field
	nameField, _ := profileType.FieldByName("Name")
	assert.Contains(t, string(nameField.Tag), `json:"name"`)

	// Check LastName field
	lastNameField, _ := profileType.FieldByName("LastName")
	assert.Contains(t, string(lastNameField.Tag), `json:"last_name"`)

	// Check DateBirth field
	dateBirthField, _ := profileType.FieldByName("DateBirth")
	assert.Contains(t, string(dateBirthField.Tag), `json:"date_birth"`)

	// Check About field
	aboutField, _ := profileType.FieldByName("About")
	assert.Contains(t, string(aboutField.Tag), `json:"about"`)

	// Check UserID field
	userIDField, _ := profileType.FieldByName("UserID")
	assert.Contains(t, string(userIDField.Tag), `json:"user_id"`)

	// Check GenderID field
	genderIDField, _ := profileType.FieldByName("GenderID")
	assert.Contains(t, string(genderIDField.Tag), `json:"gender_id"`)

	// Check CountryID field
	countryIDField, _ := profileType.FieldByName("CountryID")
	assert.Contains(t, string(countryIDField.Tag), `json:"country_id"`)
}

func TestProfileGormTags(t *testing.T) {
	// This test ensures that the GORM tags are properly set
	profile := Profile{}
	profileType := reflect.TypeOf(profile)

	// Check Name field GORM tags
	nameField, _ := profileType.FieldByName("Name")
	assert.Contains(t, string(nameField.Tag), `gorm:"not null;size:100"`)

	// Check LastName field GORM tags
	lastNameField, _ := profileType.FieldByName("LastName")
	assert.Contains(t, string(lastNameField.Tag), `gorm:"not null;size:100"`)

	// Check DateBirth field GORM tags
	dateBirthField, _ := profileType.FieldByName("DateBirth")
	assert.Contains(t, string(dateBirthField.Tag), `gorm:"not null"`)

	// Check About field GORM tags
	aboutField, _ := profileType.FieldByName("About")
	assert.Contains(t, string(aboutField.Tag), `gorm:"type:text"`)

	// Check UserID field GORM tags
	userIDField, _ := profileType.FieldByName("UserID")
	assert.Contains(t, string(userIDField.Tag), `gorm:"type:int;not null"`)

	// Check GenderID field GORM tags
	genderIDField, _ := profileType.FieldByName("GenderID")
	assert.Contains(t, string(genderIDField.Tag), `gorm:"not null"`)

	// Check CountryID field GORM tags
	countryIDField, _ := profileType.FieldByName("CountryID")
	assert.Contains(t, string(countryIDField.Tag), `gorm:"not null"`)
}

func TestProfileRelationships(t *testing.T) {
	// Test that relationship fields exist
	profile := Profile{}
	profileType := reflect.TypeOf(profile)

	// Check User relationship
	userField, exists := profileType.FieldByName("User")
	assert.True(t, exists)
	assert.Equal(t, "User", userField.Type.Name())
	assert.Contains(t, string(userField.Tag), `gorm:"foreignKey:UserID;references:ID"`)

	// Check Gender relationship
	genderField, exists := profileType.FieldByName("Gender")
	assert.True(t, exists)
	assert.Equal(t, "Gender", genderField.Type.Name())
	assert.Contains(t, string(genderField.Tag), `gorm:"foreignKey:GenderID;references:ID"`)

	// Check Country relationship
	countryField, exists := profileType.FieldByName("Country")
	assert.True(t, exists)
	assert.Equal(t, "Country", countryField.Type.Name())
	assert.Contains(t, string(countryField.Tag), `gorm:"foreignKey:CountryID;references:ID"`)
}

func TestProfileZeroValues(t *testing.T) {
	// Test zero values
	profile := Profile{}

	assert.Equal(t, "", profile.Name)
	assert.Equal(t, "", profile.LastName)
	assert.True(t, profile.DateBirth.IsZero())
	assert.Equal(t, "", profile.About)
	assert.Equal(t, 0, profile.UserID)
	assert.Equal(t, uint(0), profile.GenderID)
	assert.Equal(t, uint(0), profile.CountryID)
}
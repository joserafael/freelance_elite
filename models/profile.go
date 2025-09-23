package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null;size:100" validate:"required"`
	LastName    string    `json:"last_name" gorm:"not null;size:100" validate:"required"`
	DateBirth   time.Time `json:"date_birth" gorm:"not null" validate:"required"`
	About       string    `json:"about" gorm:"type:text"`
	UserID      int       `json:"user_id" gorm:"type:int;not null" validate:"required"`
	GenderID    uint      `json:"gender_id" gorm:"not null" validate:"required"`
	CountryID   uint      `json:"country_id" gorm:"not null" validate:"required"`

	// Relationships
	User    User    `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Gender  Gender  `json:"gender" gorm:"foreignKey:GenderID;references:ID"`
	Country Country `json:"country" gorm:"foreignKey:CountryID;references:ID"`
}

// TableName specifies the table name for the Profile model
func (Profile) TableName() string {
	return "profiles"
}
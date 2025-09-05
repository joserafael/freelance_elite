package models

import (
	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique;not null;size:100" validate:"required"`
	Code        string `json:"code" gorm:"unique;not null;size:3" validate:"required"`
	Region      string `json:"region" gorm:"size:100"`
	Subregion   string `json:"subregion" gorm:"size:100"`
	Capital     string `json:"capital" gorm:"size:100"`
	Population  int64  `json:"population"`
	Area        float64 `json:"area"`
	Currency    string `json:"currency" gorm:"size:50"`
	Language    string `json:"language" gorm:"size:100"`
	IsActive    bool   `json:"is_active"`
}

// TableName specifies the table name for the Country model
func (Country) TableName() string {
	return "countries"
}
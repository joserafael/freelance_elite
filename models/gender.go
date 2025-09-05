package models

import "gorm.io/gorm"

type Gender struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
}
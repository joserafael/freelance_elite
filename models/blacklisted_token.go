package models

import (
	"time"

	"gorm.io/gorm"
)

type BlacklistedToken struct {
	gorm.Model
	Token     string    `gorm:"type:varchar(500);uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
}

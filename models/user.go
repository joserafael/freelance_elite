package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID                   int    `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
	Username             string `json:"username"`
	Email                string `json:"email" gorm:"unique"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation" gorm:"-"`
}

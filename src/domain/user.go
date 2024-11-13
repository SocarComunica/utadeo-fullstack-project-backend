package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
	DNI      string `gorm:"unique;not null"`
	Type     string `gorm:"not null"`
}

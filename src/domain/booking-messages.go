package domain

import "gorm.io/gorm"

type BookingMessage struct {
	gorm.Model
	BookingID uint    `gorm:"not null"`
	Booking   Booking `gorm:"not null"`
	Message   string  `gorm:"not null"`
}

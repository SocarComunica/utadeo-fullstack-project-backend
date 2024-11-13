package domain

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	Status           string  `gorm:"not null"`
	BrandModel       string  `gorm:"not null"`
	Brand            string  `gorm:"not null"`
	TransmissionType string  `gorm:"not null"`
	Year             int     `gorm:"not null"`
	Type             string  `gorm:"not null"`
	HourlyFare       float64 `gorm:"not null"`
	Bookings         []Booking
}

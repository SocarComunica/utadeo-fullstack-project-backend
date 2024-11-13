package domain

import (
	"gorm.io/gorm"
	"time"
)

type Booking struct {
	gorm.Model
	Status          string `gorm:"not null"`
	UserID          uint   `gorm:"not null"`
	User            User
	VehicleID       uint `gorm:"not null"`
	Vehicle         Vehicle
	Observations    *string
	Rating          *int
	Feedback        *string
	StartDate       time.Time `gorm:"not null"`
	EndDate         time.Time
	PickUpLocation  string  `gorm:"not null"`
	DropOffLocation string  `gorm:"not null"`
	HourlyFare      float64 `gorm:"not null"`
	Messages        []BookingMessage
}

package sql

import (
	"errors"
	"time"

	"backend/src/commons"
	"backend/src/domain"
	"backend/src/services"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client interface {
	services.UsersDatabase
	services.BookingsDatabase
}

type client struct {
	DB *gorm.DB
}

func NewClient(dsn string) Client {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Here should be added migrations
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Booking{},
		&domain.BookingMessage{},
		&domain.Vehicle{},
	); err != nil {
		log.Error(err)
	}

	return &client{
		DB: db,
	}
}

func (c client) GetUserById(id uint) (*domain.User, error) {
	var user domain.User
	result := c.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return &user, nil
}

func (c client) GetVehicleById(id uint) (*domain.Vehicle, error) {
	var vehicle domain.Vehicle
	result := c.DB.First(&vehicle, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return &vehicle, nil
}

func (c client) CreateBooking(booking domain.Booking) (*domain.Booking, error) {
	result := c.DB.Create(&booking)
	if result.Error != nil {
		return nil, result.Error
	}

	return &booking, nil
}

func (c client) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := c.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return &user, nil
}

func (c client) CreateUser(user domain.User) (*domain.User, error) {
	result := c.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (c client) GetAvailableVehicles(from time.Time, to time.Time) ([]domain.Vehicle, error) {
	var vehicles []domain.Vehicle

	subQuery := c.DB.Model(&domain.Booking{}).
		Select("vehicle_id").
		Where("start_date < ? AND end_date > ? AND status IN (?, ?)",
			to,
			from,
			commons.BookingStatusReserved,
			commons.BookingStatusConfirmed)
	result := c.DB.Model(&domain.Vehicle{}).
		Where("id NOT IN (?)", subQuery).
		Find(&vehicles)
	if result.Error != nil {
		return nil, result.Error
	}

	return vehicles, nil
}

func (c client) GetBookingById(id uint) (*domain.Booking, error) {
	var booking domain.Booking
	result := c.DB.
		Preload("Vehicle").
		Preload("Messages").
		First(&booking, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return &booking, nil
}

func (c client) GetBookingsByUserID(userID uint) ([]domain.Booking, error) {
	var bookings []domain.Booking
	result := c.DB.
		Preload("Vehicle").
		Preload("Messages").
		Where("user_id = ?", userID).
		Order("start_date asc").
		Find(&bookings)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}
	return bookings, nil
}

func (c client) GetAdminBookings() ([]domain.Booking, error) {
	var bookings []domain.Booking
	result := c.DB.
		Preload("Vehicle").
		Preload("Messages").
		Order("start_date asc").
		Find(&bookings)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}
	return bookings, nil
}

func (c client) UpdateBooking(booking domain.Booking) (*domain.Booking, error) {
	result := c.DB.Save(&booking)
	if result.Error != nil {
		return nil, result.Error
	}

	return &booking, nil
}

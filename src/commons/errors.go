package commons

import (
	"errors"
)

var (
	ErrUserAlreadyExists          = errors.New("user already exists")
	ErrUserNotFound               = errors.New("user not found")
	ErrInvalidCredentials         = errors.New("invalid credentials")
	ErrVehicleNotAvailable        = errors.New("vehicle not available")
	ErrVehicleNotFound            = errors.New("vehicle not found")
	ErrBookingNotFound            = errors.New("booking not found")
	ErrBookingAlreadyCancelled    = errors.New("booking already cancelled")
	ErrBookingAlreadyFinished     = errors.New("booking already finished")
	ErrBookingAlreadyStarted      = errors.New("booking already started")
	ErrBookingNotStarted          = errors.New("booking not started")
	ErrBookingNotFinished         = errors.New("booking not finished")
	ErrBookingAlreadyHaveFeedback = errors.New("booking already have feedback")
)

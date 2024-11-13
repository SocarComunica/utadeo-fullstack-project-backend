package requests

import "time"

type AvailableVehiclesResponse struct {
	ID               uint    `json:"id"`
	Status           string  `json:"status"`
	BrandModel       string  `json:"brand_model"`
	Brand            string  `json:"brand"`
	TransmissionType string  `json:"transmission_type"`
	Year             int     `json:"year"`
	Type             string  `json:"type"`
	HourlyFare       float64 `json:"hourly_fare"`
}

type CreateBookingRequest struct {
	UserID          uint      `json:"user_id" validate:"required"`
	VehicleID       uint      `json:"vehicle_id" validate:"required"`
	StartDate       time.Time `json:"start_date" validate:"required"`
	EndDate         time.Time `json:"end_date" validate:"required"`
	PickUpLocation  string    `json:"pick_up_location" validate:"required"`
	DropOffLocation string    `json:"drop_off_location" validate:"required"`
}

type BookingResponse struct {
	ID              uint                      `json:"id"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	Status          string                    `json:"status"`
	UserID          uint                      `json:"user_id"`
	Vehicle         AvailableVehiclesResponse `json:"vehicle"`
	Observations    *string                   `json:"observations"`
	Feedback        *string                   `json:"feedback"`
	Rating          *int                      `json:"rating"`
	StarDate        time.Time                 `json:"start_date"`
	EndDate         time.Time                 `json:"end_date"`
	PickUpLocation  string                    `json:"pick_up_location"`
	DropOffLocation string                    `json:"drop_off_location"`
	HourlyFare      float64                   `json:"hourly_fare"`
	TotalAmount     float64                   `json:"total_amount"`
	Messages        []MessagesResponse        `json:"messages"`
}

type MessagesResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	BookingID uint      `json:"booking_id"`
	Message   string    `json:"message"`
}

type CancelBookingRequest struct {
	ID     uint `json:"id" validate:"required"`
	UserID uint `json:"user_id" validate:"required"`
}

type ConfirmBookingRequest struct {
	ID     uint `json:"id" validate:"required"`
	UserID uint `json:"user_id" validate:"required"`
}

type FinishBookingRequest struct {
	ID     uint `json:"id" validate:"required"`
	UserID uint `json:"user_id" validate:"required"`
}

type AddFeedbackBookingRequest struct {
	ID       uint   `json:"id" validate:"required"`
	UserID   uint   `json:"user_id" validate:"required"`
	Feedback string `json:"feedback" validate:"required"`
}

type RateBookingRequest struct {
	ID     uint `json:"id" validate:"required"`
	UserID uint `json:"user_id" validate:"required"`
	Rating int  `json:"rating" validate:"required"`
}

type AddMessageToBookingRequest struct {
	ID      uint   `json:"id" validate:"required"`
	UserID  uint   `json:"user_id" validate:"required"`
	Message string `json:"message" validate:"required"`
}

package services

import (
	"time"

	"backend/src/commons"
	"backend/src/domain"
	"backend/src/handlers/requests"
)

type BookingsDatabase interface {
	GetAvailableVehicles(from time.Time, to time.Time) ([]domain.Vehicle, error)
	GetUserById(id uint) (*domain.User, error)
	GetVehicleById(id uint) (*domain.Vehicle, error)
	GetBookingById(id uint) (*domain.Booking, error)
	CreateBooking(booking domain.Booking) (*domain.Booking, error)
	UpdateBooking(booking domain.Booking) (*domain.Booking, error)
	GetBookingsByUserID(userID uint) ([]domain.Booking, error)
	GetAdminBookings() ([]domain.Booking, error)
}

type BookingsService struct {
	Database BookingsDatabase
}

func NewBookingsService(database BookingsDatabase) *BookingsService {
	return &BookingsService{
		Database: database,
	}
}

func (b *BookingsService) GetAvailableVehicles(from time.Time, to time.Time) ([]domain.Vehicle, error) {
	return b.Database.GetAvailableVehicles(from, to)
}

func (b *BookingsService) CreateBooking(request requests.CreateBookingRequest) (*domain.Booking, error) {
	user, err := b.Database.GetUserById(request.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, commons.ErrUserNotFound
	}

	vehicle, err := b.Database.GetVehicleById(request.VehicleID)
	if err != nil {
		return nil, err
	}

	if vehicle == nil {
		return nil, commons.ErrVehicleNotFound
	}

	booking := &domain.Booking{
		Status:          commons.BookingStatusReserved,
		UserID:          request.UserID,
		User:            *user,
		VehicleID:       request.VehicleID,
		Vehicle:         *vehicle,
		StartDate:       request.StartDate,
		EndDate:         request.EndDate,
		PickUpLocation:  request.PickUpLocation,
		DropOffLocation: request.DropOffLocation,
		HourlyFare:      vehicle.HourlyFare,
	}

	return b.Database.CreateBooking(*booking)
}

func (b *BookingsService) CancelBooking(request requests.CancelBookingRequest) (*domain.Booking, error) {
	user, err := b.Database.GetUserById(request.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, commons.ErrInvalidCredentials
	}

	booking, err := b.Database.GetBookingById(request.ID)
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, commons.ErrBookingNotFound
	}

	if booking.UserID != request.UserID && user.Type != commons.UserTypeAdmin {
		return nil, commons.ErrInvalidCredentials
	}

	if booking.Status == commons.BookingStatusCancelled {
		return nil, commons.ErrBookingAlreadyCancelled
	}

	if booking.Status == commons.BookingStatusFinished {
		return nil, commons.ErrBookingAlreadyFinished
	}

	if booking.Status == commons.BookingStatusConfirmed {
		return nil, commons.ErrBookingAlreadyStarted
	}

	if booking.Status == commons.BookingStatusReserved {
		booking.Status = commons.BookingStatusCancelled
	}

	var message string
	if user.Type == commons.UserTypeAdmin {
		message = "Booking cancelled by admin"
	} else {
		message = "Booking cancelled by user"
	}
	booking.Observations = &message

	return b.Database.UpdateBooking(*booking)
}

func (b *BookingsService) ConfirmBooking(request requests.ConfirmBookingRequest) (*domain.Booking, error) {
	user, err := b.Database.GetUserById(request.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, commons.ErrInvalidCredentials
	}

	booking, err := b.Database.GetBookingById(request.ID)
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, commons.ErrBookingNotFound
	}

	if booking.UserID != request.UserID && user.Type != commons.UserTypeAdmin {
		return nil, commons.ErrInvalidCredentials
	}

	if booking.Status == commons.BookingStatusCancelled {
		return nil, commons.ErrBookingAlreadyCancelled
	}

	if booking.Status == commons.BookingStatusFinished {
		return nil, commons.ErrBookingAlreadyFinished
	}

	if booking.Status == commons.BookingStatusConfirmed {
		return nil, commons.ErrBookingAlreadyStarted
	}

	if booking.Status == commons.BookingStatusReserved {
		booking.Status = commons.BookingStatusConfirmed
	}

	var message string
	if user.Type == commons.UserTypeAdmin {
		message = "Booking confirmed by admin"
	} else {
		message = "Booking confirmed by user"
	}
	booking.Observations = &message

	return b.Database.UpdateBooking(*booking)
}

func (b *BookingsService) FinishBooking(request requests.FinishBookingRequest) (*domain.Booking, error) {
	user, err := b.Database.GetUserById(request.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, commons.ErrInvalidCredentials
	}

	booking, err := b.Database.GetBookingById(request.ID)
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, commons.ErrBookingNotFound
	}

	if booking.UserID != request.UserID && user.Type != commons.UserTypeAdmin {
		return nil, commons.ErrInvalidCredentials
	}

	if booking.Status == commons.BookingStatusCancelled {
		return nil, commons.ErrBookingAlreadyCancelled
	}

	if booking.Status == commons.BookingStatusFinished {
		return nil, commons.ErrBookingAlreadyFinished
	}

	if booking.Status == commons.BookingStatusReserved {
		return nil, commons.ErrBookingNotStarted
	}

	if booking.Status == commons.BookingStatusConfirmed {
		booking.Status = commons.BookingStatusFinished
	}

	var message string
	if user.Type == commons.UserTypeAdmin {
		message = "Booking finished by admin"
	} else {
		message = "Booking finished by user"
	}
	booking.Observations = &message

	return b.Database.UpdateBooking(*booking)
}

func (b *BookingsService) AddFeedbackBooking(request requests.AddFeedbackBookingRequest) (*domain.Booking, error) {
	user, err := b.Database.GetUserById(request.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, commons.ErrInvalidCredentials
	}

	booking, err := b.Database.GetBookingById(request.ID)
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, commons.ErrBookingNotFound
	}

	if booking.UserID != request.UserID || user.Type != commons.UserTypeClient {
		return nil, commons.ErrInvalidCredentials
	}

	if booking.Status == commons.BookingStatusReserved {
		return nil, commons.ErrBookingNotStarted
	}

	if booking.Status == commons.BookingStatusConfirmed {
		return nil, commons.ErrBookingNotFinished
	}

	if booking.Status == commons.BookingStatusCancelled {
		return nil, commons.ErrBookingAlreadyCancelled
	}

	if booking.Status != commons.BookingStatusFinished {
		return nil, commons.ErrBookingNotFinished
	}

	if booking.Feedback != nil && *booking.Feedback != "" {
		return nil, commons.ErrBookingAlreadyHaveFeedback
	}

	booking.Feedback = &request.Feedback

	return b.Database.UpdateBooking(*booking)
}

func (b *BookingsService) RateBooking(request requests.RateBookingRequest) (*domain.Booking, error) {
	user, err := b.Database.GetUserById(request.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, commons.ErrInvalidCredentials
	}

	booking, err := b.Database.GetBookingById(request.ID)
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, commons.ErrBookingNotFound
	}

	if booking.UserID != request.UserID || user.Type != commons.UserTypeClient {
		return nil, commons.ErrInvalidCredentials
	}

	if booking.Status == commons.BookingStatusReserved {
		return nil, commons.ErrBookingNotStarted
	}

	if booking.Status == commons.BookingStatusConfirmed {
		return nil, commons.ErrBookingNotFinished
	}

	if booking.Status == commons.BookingStatusCancelled {
		return nil, commons.ErrBookingAlreadyCancelled
	}

	if booking.Status != commons.BookingStatusFinished {
		return nil, commons.ErrBookingNotFinished
	}

	booking.Rating = &request.Rating

	return b.Database.UpdateBooking(*booking)
}

func (b *BookingsService) AddMessageToBooking(request requests.AddMessageToBookingRequest) (*domain.Booking, error) {
	user, err := b.Database.GetUserById(request.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, commons.ErrInvalidCredentials
	}

	booking, err := b.Database.GetBookingById(request.ID)
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, commons.ErrBookingNotFound
	}

	if booking.UserID != request.UserID || user.Type != commons.UserTypeClient {
		return nil, commons.ErrInvalidCredentials
	}

	message := domain.BookingMessage{
		BookingID: booking.ID,
		Message:   request.Message,
	}
	booking.Messages = append(booking.Messages, message)

	return b.Database.UpdateBooking(*booking)
}

func (b *BookingsService) GetBookingsByUserID(userID uint) ([]domain.Booking, error) {
	return b.Database.GetBookingsByUserID(userID)
}

func (b *BookingsService) GetBookingByID(bookingID uint) (*domain.Booking, error) {
	booking, err := b.Database.GetBookingById(bookingID)
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, commons.ErrBookingNotFound
	}

	return booking, nil
}

func (b *BookingsService) GetAdminBookings() ([]domain.Booking, error) {
	return b.Database.GetAdminBookings()
}

package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"backend/src/commons"
	"backend/src/domain"
	"backend/src/handlers/requests"

	"github.com/labstack/echo/v4"
)

type BookingsService interface {
	GetAvailableVehicles(from time.Time, to time.Time) ([]domain.Vehicle, error)
	CreateBooking(request requests.CreateBookingRequest) (*domain.Booking, error)
	CancelBooking(request requests.CancelBookingRequest) (*domain.Booking, error)
	ConfirmBooking(request requests.ConfirmBookingRequest) (*domain.Booking, error)
	FinishBooking(request requests.FinishBookingRequest) (*domain.Booking, error)
	AddFeedbackBooking(request requests.AddFeedbackBookingRequest) (*domain.Booking, error)
	RateBooking(request requests.RateBookingRequest) (*domain.Booking, error)
	AddMessageToBooking(request requests.AddMessageToBookingRequest) (*domain.Booking, error)
	GetBookingsByUserID(userID uint) ([]domain.Booking, error)
	GetBookingByID(bookingID uint) (*domain.Booking, error)
	GetAdminBookings() ([]domain.Booking, error)
}

type BookingsHandler struct {
	service BookingsService
}

func NewBookingsHandler(service BookingsService) *BookingsHandler {
	return &BookingsHandler{
		service: service,
	}
}

func (h *BookingsHandler) AddRoutes(router *echo.Router) {
	router.Add(echo.GET, "/bookings", h.getBookings)
	router.Add(echo.GET, "/bookings/admin", h.getAdminBookings)
	router.Add(echo.GET, "/bookings/available-vehicles", h.getAvailableVehicles)
	router.Add(echo.POST, "/bookings", h.createBooking)
	router.Add(echo.POST, "/bookings/message", h.addMessageToBooking)
	router.Add(echo.PATCH, "/bookings/cancel", h.cancelBooking)
	router.Add(echo.PATCH, "/bookings/confirm", h.confirmBooking)
	router.Add(echo.PATCH, "/bookings/finish", h.finishBooking)
	router.Add(echo.PATCH, "/bookings/feedback", h.addFeedbackBooking)
	router.Add(echo.PATCH, "/bookings/rate", h.rateBooking)
}

func (h *BookingsHandler) getAvailableVehicles(c echo.Context) error {
	from, err := time.Parse(time.RFC3339, c.QueryParam("from"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid from date",
		})
	}
	to, err := time.Parse(time.RFC3339, c.QueryParam("to"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid to date",
		})
	}

	vehicles, err := h.service.GetAvailableVehicles(from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "error getting available vehicles",
		})
	}

	response := make([]*requests.AvailableVehiclesResponse, 0)
	for _, vehicle := range vehicles {
		response = append(response, mapVehicleToResponse(vehicle))
	}

	return c.JSON(http.StatusOK, response)
}

func (h *BookingsHandler) createBooking(c echo.Context) error {
	r := new(requests.CreateBookingRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if r.StartDate.Before(time.Now()) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "start date cannot be in the past",
		})
	}

	if r.StartDate.After(r.EndDate) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "start date cannot be after end date",
		})
	}

	booking, err := h.service.CreateBooking(*r)
	if err != nil {
		if errors.Is(err, commons.ErrVehicleNotAvailable) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, mapBookingToResponse(*booking))
}

func (h *BookingsHandler) cancelBooking(c echo.Context) error {
	r := new(requests.CancelBookingRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	booking, err := h.service.CancelBooking(*r)
	if err != nil {
		if errors.Is(err, commons.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyCancelled) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyFinished) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyStarted) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapBookingToResponse(*booking))
}

func (h *BookingsHandler) confirmBooking(c echo.Context) error {
	r := new(requests.ConfirmBookingRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	booking, err := h.service.ConfirmBooking(*r)
	if err != nil {
		if errors.Is(err, commons.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyCancelled) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyFinished) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyStarted) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapBookingToResponse(*booking))
}

func (h *BookingsHandler) finishBooking(c echo.Context) error {
	r := new(requests.FinishBookingRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	booking, err := h.service.FinishBooking(*r)
	if err != nil {
		if errors.Is(err, commons.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyCancelled) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyFinished) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingNotStarted) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapBookingToResponse(*booking))
}

func (h *BookingsHandler) addFeedbackBooking(c echo.Context) error {
	r := new(requests.AddFeedbackBookingRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	booking, err := h.service.AddFeedbackBooking(*r)
	if err != nil {
		if errors.Is(err, commons.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyCancelled) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingNotFinished) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyHaveFeedback) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapBookingToResponse(*booking))
}

func (h *BookingsHandler) rateBooking(c echo.Context) error {
	r := new(requests.RateBookingRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	booking, err := h.service.RateBooking(*r)
	if err != nil {
		if errors.Is(err, commons.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingAlreadyCancelled) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingNotFinished) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapBookingToResponse(*booking))
}

func (h *BookingsHandler) addMessageToBooking(c echo.Context) error {
	r := new(requests.AddMessageToBookingRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	booking, err := h.service.AddMessageToBooking(*r)
	if err != nil {
		if errors.Is(err, commons.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": err.Error(),
			})
		}

		if errors.Is(err, commons.ErrBookingNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapBookingToResponse(*booking))
}

func (h *BookingsHandler) getBookings(c echo.Context) error {
	// Tiene prioridad la busqueda de booking por ID
	bookingIDStr := c.QueryParam("booking_id")
	userIDStr := c.QueryParam("user_id")

	if bookingIDStr == "" && userIDStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid query params, required booking_id or user_id",
		})
	}

	if bookingIDStr != "" {
		bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid booking id",
			})
		}

		if bookingID != 0 {
			booking, err := h.service.GetBookingByID(uint(bookingID))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": err.Error(),
				})
			}

			if booking == nil {
				return c.JSON(http.StatusNotFound, echo.Map{
					"message": "booking not found",
				})
			}

			return c.JSON(http.StatusOK, mapBookingToResponse(*booking))
		}
	}

	if userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid user id",
			})
		}

		if userID != 0 {
			bookings, err := h.service.GetBookingsByUserID(uint(userID))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": err.Error(),
				})
			}

			response := make([]*requests.BookingResponse, 0)
			for _, booking := range bookings {
				response = append(response, mapBookingToResponse(booking))
			}

			return c.JSON(http.StatusOK, response)
		}
	}

	return c.JSON(http.StatusBadRequest, echo.Map{
		"message": "invalid query params",
	})
}

func (h *BookingsHandler) getAdminBookings(c echo.Context) error {
	bookings, err := h.service.GetAdminBookings()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	response := make([]*requests.BookingResponse, 0)
	for _, booking := range bookings {
		response = append(response, mapBookingToResponse(booking))
	}

	return c.JSON(http.StatusOK, response)
}

func mapVehicleToResponse(vehicle domain.Vehicle) *requests.AvailableVehiclesResponse {
	return &requests.AvailableVehiclesResponse{
		ID:               vehicle.ID,
		Status:           vehicle.Status,
		BrandModel:       vehicle.BrandModel,
		Brand:            vehicle.Brand,
		TransmissionType: vehicle.TransmissionType,
		Year:             vehicle.Year,
		Type:             vehicle.Type,
		HourlyFare:       vehicle.HourlyFare,
	}
}

func mapMessageToResponse(message domain.BookingMessage) *requests.MessagesResponse {
	return &requests.MessagesResponse{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
		BookingID: message.BookingID,
		Message:   message.Message,
	}
}

func mapBookingToResponse(booking domain.Booking) *requests.BookingResponse {
	vehicle := mapVehicleToResponse(booking.Vehicle)
	totalAmount := booking.EndDate.Sub(booking.StartDate).Hours() * booking.HourlyFare

	messages := make([]requests.MessagesResponse, 0)
	for _, message := range booking.Messages {
		messages = append(messages, *mapMessageToResponse(message))
	}

	return &requests.BookingResponse{
		ID:              booking.ID,
		CreatedAt:       booking.CreatedAt,
		UpdatedAt:       booking.UpdatedAt,
		Status:          booking.Status,
		UserID:          booking.UserID,
		Vehicle:         *vehicle,
		Feedback:        booking.Feedback,
		Observations:    booking.Observations,
		Rating:          booking.Rating,
		StarDate:        booking.StartDate,
		EndDate:         booking.EndDate,
		PickUpLocation:  booking.PickUpLocation,
		DropOffLocation: booking.DropOffLocation,
		HourlyFare:      booking.HourlyFare,
		TotalAmount:     totalAmount,
		Messages:        messages,
	}
}

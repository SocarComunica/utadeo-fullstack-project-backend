package handlers

import (
	"errors"
	"net/http"

	"backend/src/commons"
	"backend/src/domain"
	"backend/src/handlers/requests"

	"github.com/labstack/echo/v4"
)

type UsersService interface {
	RegisterUser(context echo.Context, request requests.RegisterUserRequest) (*domain.User, error)
	LoginUser(context echo.Context, request requests.LoginUserRequest) (*domain.User, error)
}

type UsersHandler struct {
	service UsersService
}

func NewUsersHandler(service UsersService) *UsersHandler {
	return &UsersHandler{
		service: service,
	}
}

func (u *UsersHandler) AddRoutes(router *echo.Router) {
	router.Add(echo.POST, "/users", u.RegisterUser)
	router.Add(echo.POST, "/users/login", u.LoginUser)
}

func (u *UsersHandler) RegisterUser(c echo.Context) error {
	r := new(requests.RegisterUserRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"layer": "RegisterUserError Handler",
			"error": err.Error(),
		})
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"layer": "RegisterUserError Handler",
			"error": err.Error(),
		})
	}

	user, err := u.service.RegisterUser(c, *r)
	if err != nil {
		if errors.Is(err, commons.ErrUserAlreadyExists) {
			return c.JSON(http.StatusConflict, echo.Map{
				"layer": "RegisterUserError Handler",
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": "RegisterUserError Handler",
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, mapUserToResponse(*user))
}

func (u *UsersHandler) LoginUser(c echo.Context) error {
	r := new(requests.LoginUserRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"layer": "LoginUserError Handler",
			"error": err.Error(),
		})
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"layer": "LoginUserError Handler",
			"error": err.Error(),
		})
	}

	user, err := u.service.LoginUser(c, *r)
	if err != nil {
		if errors.Is(err, commons.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"layer": "LoginUserError Handler",
				"error": err.Error(),
			})
		}
		if errors.Is(err, commons.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"layer": "LoginUserError Handler",
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": "LoginUserError Handler",
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapUserToResponse(*user))
}

func mapUserToResponse(user domain.User) *requests.UserResponse {
	return &requests.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		DNI:   user.DNI,
		Type:  user.Type,
	}
}

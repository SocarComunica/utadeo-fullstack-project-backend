package services

import (
	"backend/src/commons"
	"backend/src/domain"
	"backend/src/handlers/requests"

	"github.com/labstack/echo/v4"
)

type UsersDatabase interface {
	GetUserByEmail(email string) (*domain.User, error)
	CreateUser(user domain.User) (*domain.User, error)
}

type UsersService struct {
	Database UsersDatabase
}

func NewUsersService(database UsersDatabase) *UsersService {
	return &UsersService{
		Database: database,
	}
}

func (u *UsersService) RegisterUser(context echo.Context, request requests.RegisterUserRequest) (*domain.User, error) {
	user, err := u.Database.GetUserByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, commons.ErrUserAlreadyExists
	}

	user = mapRegisterUserRequestToUser(request)

	return u.Database.CreateUser(*user)
}

func (u *UsersService) LoginUser(context echo.Context, request requests.LoginUserRequest) (*domain.User, error) {
	user, err := u.Database.GetUserByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, commons.ErrUserNotFound
	}

	if user.Password != request.Password {
		return nil, commons.ErrInvalidCredentials
	}

	return user, nil
}

func mapRegisterUserRequestToUser(request requests.RegisterUserRequest) *domain.User {
	return &domain.User{
		Email:    request.Email,
		Name:     request.Name,
		Password: request.Password,
		DNI:      request.DNI,
		Type:     commons.UserTypeClient,
	}
}

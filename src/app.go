package src

import (
	"log"
	"os"

	handlers2 "backend/src/handlers"
	"backend/src/services"
	"backend/src/sql"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Enable cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Custom validator
	e.Validator = NewCustomValidator()

	// Setup database
	database := sql.NewClient(os.Getenv("DATABASE_DSN"))

	// Users handler
	usersService := services.NewUsersService(database)
	usersHandler := handlers2.NewUsersHandler(usersService)

	// Bookings handler
	bookingsService := services.NewBookingsService(database)
	bookingsHandler := handlers2.NewBookingsHandler(bookingsService)

	handlers := []handlers2.Handler{
		usersHandler,
		bookingsHandler,
	}

	for _, handler := range handlers {
		handler.AddRoutes(e.Router())
	}

	return e.Start(os.Getenv("SERVER_PORT"))
}

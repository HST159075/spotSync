package server

import (
	"spotsync/internal/config"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "SpotSync API is running! 🚗")
	})

	// Dependency Injection
	userRepo := user.NewRepository(config.DB)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	zoneRepo := zone.NewRepository(config.DB)
	zoneService := zone.NewService(zoneRepo)
	zoneHandler := zone.NewHandler(zoneService)

	resRepo := reservation.NewRepository(config.DB)
	resService := reservation.NewService(resRepo)
	resHandler := reservation.NewHandler(resService)

	// Routes
	api := e.Group("/api/v1")

	// Auth (public)
	auth := api.Group("/auth")
	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)

	// Zones
	zones := api.Group("/zones")
	zones.GET("", zoneHandler.GetAllZones)
	zones.GET("/:id", zoneHandler.GetZoneByID)
	zones.POST("", zoneHandler.CreateZone, middlewares.JWTMiddleware)
	zones.PUT("/:id", zoneHandler.UpdateZone, middlewares.JWTMiddleware)
	zones.DELETE("/:id", zoneHandler.DeleteZone, middlewares.JWTMiddleware)

	// Reservations
	res := api.Group("/reservations", middlewares.JWTMiddleware)
	res.POST("", resHandler.CreateReservation)
	res.GET("/my-reservations", resHandler.GetMyReservations)
	res.DELETE("/:id", resHandler.CancelReservation)
	res.GET("", resHandler.GetAllReservations)

	return e
}
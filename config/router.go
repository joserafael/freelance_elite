package config

import (
	"os"
	"freelance_elite/handlers"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)
	e.POST("/logout", handlers.Logout)

	p := e.Group("/profile")
	p.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))
	p.Use(handlers.CheckBlacklist)
	p.GET("", handlers.Profile)
}

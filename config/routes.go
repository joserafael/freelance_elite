package config

import (
	"os"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"freelance_elite/handlers"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)
	e.POST("/logout", handlers.Logout)

	// Gender routes (public access)
	genders := e.Group("/genders")
	genders.GET("", handlers.GetGenders)
	genders.GET("/:id", handlers.GetGender)
	genders.POST("", handlers.CreateGender)
	genders.PUT("/:id", handlers.UpdateGender)
	genders.DELETE("/:id", handlers.DeleteGender)

	// Countries routes (public access)
	countries := e.Group("/countries")
	countries.GET("", handlers.GetCountries)
	countries.GET("/:id", handlers.GetCountry)
	countries.POST("", handlers.CreateCountry)
	countries.PUT("/:id", handlers.UpdateCountry)
	countries.DELETE("/:id", handlers.DeleteCountry)
	countries.GET("/region/:region", handlers.GetCountriesByRegion)

	p := e.Group("/profile")
	p.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))
	p.Use(handlers.CheckBlacklist)
	p.GET("", handlers.Profile)
}

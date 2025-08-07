package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"

	"freelance_elite/database"
	"freelance_elite/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

	p := e.Group("/profile")
	p.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))
	p.GET("", handlers.Profile)

	e.Logger.Fatal(e.Start(":1323"))
}

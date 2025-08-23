package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"freelance_elite/config"
)

func main() {
	envFile := ".env"
	if os.Getenv("APP_ENV") == "test" {
		envFile = ".env.test"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	config.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}

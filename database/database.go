package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Connect to MySQL server without specifying a database to create it if it doesn't exist
	initialDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort)

	sqlDBForCreation, err := sql.Open("mysql", initialDSN)
	if err != nil {
		log.Fatalf("failed to connect to mysql server for creation: %v", err)
	}
	defer sqlDBForCreation.Close()

	// Create the database if it doesn't exist
	_, err = sqlDBForCreation.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	// Now, connect to the specific database for GORM and migrations
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName)

	DB, err = gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database '%s': %v", dbName, err)
	}

	// Get the underlying sql.DB for migrations
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from GORM: %v", err)
	}
	// Do NOT close sqlDB here; GORM manages its lifecycle. It will be closed when the application exits.

	// Run migrations
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName,
		driver,
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database: %v", err)
	}
}

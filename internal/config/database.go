package config

import (
	"fmt"
	"log"
	"os"

	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&user.Model{},
		&zone.Model{},
		&reservation.Model{},
	)
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	log.Println("Database connected and migrated successfully")
	DB = db
}
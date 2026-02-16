package database

import (
	"fmt"
	"log"

	"saasproject/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg := config.Get()
	sslMode := cfg.Database.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		sslMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
		return
	}
	DB = db
}

func Connection() *gorm.DB {
	return DB
}

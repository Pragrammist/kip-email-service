package database

import (
	"fmt"
	"go-email/config"
	"go-email/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DbName,
		cfg.Database.Port,
		cfg.Database.SslMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Cannot connect to the database %s", err.Error())
	}

	// migrate model
	err = db.AutoMigrate(&models.EmailModel{})

	if err != nil {
		log.Fatalf("Cannot migrate models to the database: %s", err.Error())
	}

	return db
}

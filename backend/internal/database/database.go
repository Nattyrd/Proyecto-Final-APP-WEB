package database

import (
	"fmt"
	"log"

	"github.com/grupo5/ecommerce-api/internal/config"
	"github.com/grupo5/ecommerce-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Info
	if cfg.AppEnv == "production" {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("conectar a PostgreSQL: %w", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Receipt{},
		&models.ReceiptItem{},
	); err != nil {
		return nil, fmt.Errorf("ejecutar migraciones: %w", err)
	}

	log.Println("Base de datos conectada y migrada correctamente")
	return db, nil
}

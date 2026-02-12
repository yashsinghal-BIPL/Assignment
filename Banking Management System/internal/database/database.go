package database

import (
	"bank-management-system/internal/config"
	"bank-management-system/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) error {
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Database connected successfully")

	// Auto migrate all models
	err = db.AutoMigrate(
		&models.Bank{},
		&models.Branch{},
		&models.Customer{},
		&models.Account{},
		&models.CustomerAccount{},
		&models.Loan{},
		&models.Transaction{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

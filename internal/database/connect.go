package database

import (
	"fmt"
	"log"
	"strconv"

	config "app/configs"
	"app/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connects to the Neon PostgreSQL database
func ConnectDB() error {
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse database port: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		config.Config("DB_HOST"),
		uint(port),
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database Migrated: Users, Products tables created")
	// Run migrations
	err = DB.AutoMigrate(&model.User{}, &model.Product{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("Database Migrated")
	return nil
}

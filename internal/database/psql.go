package database

import (
	"auth_project/config"
	entities "auth_project/internal/entities/user"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(config config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			config.Host,
			config.DBPort,
			config.Username,
			config.DBName,
			config.SSLMode,
			config.Password)),
		&gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entities.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

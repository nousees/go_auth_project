package database

import (
	"auth_project/internal/domain"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConnectionInfo struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func NewPostgresConnection(info ConnectionInfo) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s", info.Host, info.Port, info.Username, info.DBName, info.SSLMode, info.Password)),
		&gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

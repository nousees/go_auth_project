package repository

import (
	"auth_project/internal/domain"

	"gorm.io/gorm"
)

type Users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{db}
}

func (us *Users) Create(user *domain.User) error {
	return us.db.Create(user).Error
}

func (us *Users) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := us.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

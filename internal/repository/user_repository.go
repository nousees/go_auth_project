package repository

import (
	entities "auth_project/internal/entities/user"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
}

type Users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{db}
}

func (us *Users) Create(user *entities.User) error {
	return us.db.Create(user).Error
}

func (us *Users) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := us.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

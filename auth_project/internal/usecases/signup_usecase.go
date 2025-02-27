package usecases

import (
	"auth_project/internal/domain"
	"auth_project/internal/repository"
	"auth_project/pkg/hash"
)

type SignUpUsecase struct {
	Users *repository.Users
}

func NewSignUpUsecase(users *repository.Users) *SignUpUsecase {
	return &SignUpUsecase{users}
}

func (s *SignUpUsecase) SignUp(sup domain.SignUpInput) error {
	hashPass, err := hash.GenerateHash(sup.Password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    sup.Email,
		Password: hashPass,
	}

	return s.Users.Create(user)
}

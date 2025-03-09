package usecases

import (
	entities "auth_project/internal/entities/user"
	"auth_project/internal/repository"
	"auth_project/pkg/hash"
)

type SignUpUsecase struct {
	Users repository.UserRepository
}

func NewSignUpUsecase(users repository.UserRepository) *SignUpUsecase {
	return &SignUpUsecase{users}
}

func (s *SignUpUsecase) SignUp(sup entities.SignUpInput) error {
	hashPass, err := hash.GenerateHash(sup.Password)
	if err != nil {
		return err
	}

	user := &entities.User{
		Email:    sup.Email,
		Password: hashPass,
	}

	return s.Users.Create(user)
}

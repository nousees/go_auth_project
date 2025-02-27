package usecases

import (
	"auth_project/internal/domain"
	"auth_project/internal/repository"
	"auth_project/pkg/hash"
	"auth_project/pkg/jwt"
	"errors"
)

type SignInUsecase struct {
	Users *repository.Users
}

func NewSignInUsecase(users *repository.Users) *SignInUsecase {
	return &SignInUsecase{users}
}

func (s *SignInUsecase) SignIn(sinInput domain.SignInInput) (string, error) {
	user, err := s.Users.GetUserByEmail(sinInput.Email)
	if err != nil {
		return "", errors.New("user no found")
	}

	if !hash.CompareHash(user.Password, sinInput.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

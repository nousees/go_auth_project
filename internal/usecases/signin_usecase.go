package usecases

import (
	entities "auth_project/internal/entities/user"
	"auth_project/internal/repository"
	"auth_project/pkg/hash"
	"auth_project/pkg/jwt"
	"errors"
)

type SignInUsecase struct {
	Users repository.UserRepository
}

func NewSignInUsecase(users repository.UserRepository) *SignInUsecase {
	return &SignInUsecase{users}
}

func (s *SignInUsecase) SignIn(sinInput entities.SignInInput) (string, error) {
	user, err := s.Users.GetUserByEmail(sinInput.Email)
	if err != nil {
		return "", errors.New("user not found")
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

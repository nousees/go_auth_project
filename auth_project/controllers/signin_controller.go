package controllers

import (
	entities "auth_project/internal/entities/user"
	"auth_project/internal/usecases"

	"github.com/gin-gonic/gin"
)

type SignInController struct {
	SignInUsecase usecases.SignInUsecase
}

func NewSignInController(sin usecases.SignInUsecase) *SignInController {
	return &SignInController{sin}
}

func (sc *SignInController) SignIn(c *gin.Context) {
	var user entities.SignInInput
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := sc.SignInUsecase.SignIn(user)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "authorizzation successfully", "token": token})
}

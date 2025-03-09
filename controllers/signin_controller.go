package controllers

import (
	entities "auth_project/internal/entities/user"
	"auth_project/internal/usecases"
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "bad request", "error": err.Error()})
		return
	}

	token, err := sc.SignInUsecase.SignIn(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "authorization error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "authorization successfully", "token": token})
}

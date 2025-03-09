package controllers

import (
	entities "auth_project/internal/entities/user"
	"auth_project/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignUpController struct {
	SignUpUsecase usecases.SignUpUsecase
}

func NewSignUpController(sup usecases.SignUpUsecase) *SignUpController {
	return &SignUpController{sup}
}

func (sc *SignUpController) SignUp(c *gin.Context) {
	var user entities.SignUpInput
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "bad request", "error": err.Error()})
		return
	}

	err := sc.SignUpUsecase.SignUp(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "registration error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "registration successfully"})
}

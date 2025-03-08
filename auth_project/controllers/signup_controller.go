package controllers

import (
	entities "auth_project/internal/entities/user"
	"auth_project/internal/usecases"

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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := sc.SignUpUsecase.SignUp(user)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "register successfully"})
}

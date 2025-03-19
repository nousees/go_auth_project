package main

import (
	"auth_project/config"
	"auth_project/controllers"
	"auth_project/internal/database"
	"auth_project/internal/repository"
	"auth_project/internal/usecases"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgresConnection(cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}

	users := repository.NewUsers(db)
	signUpUsecase := usecases.NewSignUpUsecase(users)
	signInUsecase := usecases.NewSignInUsecase(users)

	signInController := controllers.NewSignInController(*signInUsecase)
	signUpController := controllers.NewSignUpController(*signUpUsecase)

	router := gin.Default()

	router.POST("/sign-up", signUpController.SignUp)
	router.POST("/sign-in", signInController.SignIn)

	router.Run(":" + cfg.Server.Port)
}

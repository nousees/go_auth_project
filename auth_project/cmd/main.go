package main

import (
	"auth_project/internal/controllers"
	"auth_project/internal/database"
	"auth_project/internal/repository"
	"auth_project/internal/usecases"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "go_db",
		SSLMode:  "disable",
		Password: "135267984",
	})
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}

	users := repository.NewUsers(db)
	signUpUsecase := usecases.NewSignUpUsecase(users)
	signIpUsecase := usecases.NewSignInUsecase(users)

	signInController := controllers.NewSignInUsecase(*signIpUsecase)
	signUpController := controllers.NewSignUpController(*signUpUsecase)

	router := gin.Default()

	router.POST("/sign-up", signUpController.SignUp)
	router.POST("/sign-in", signInController.SignIn)

	router.Run(":8080")
}

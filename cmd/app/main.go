package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/peruccii/roadmap-go-backend/internal/api"
	"github.com/peruccii/roadmap-go-backend/internal/db"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	database, err := db.InitDB("sqlite", "test.db")
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// Migrar estrutura do banco
	database.AutoMigrate(&models.User{}, &models.Robots{})

	userRepo := repository.NewUserRepository(database)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	stripeService := services.NewStripeService()

	r := gin.Default()
	api.SetupRoutes(r, userService, authService, stripeService)

	r.Run(":8080")
}

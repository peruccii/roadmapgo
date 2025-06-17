package app

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/peruccii/roadmap-go-backend/internal/db"
	"github.com/peruccii/roadmap-go-backend/internal/models"
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

	database.AutoMigrate(&models.User{}, &models.Courses{})
}

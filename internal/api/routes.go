package api

import (
	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/controller"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

func SetupRoutes(r *gin.Engine, userService services.UserService) {
	userController := controller.NewUserController(userService)

	users := r.Group("/users")
	{
		users.POST("", userController.Create)
		users.GET("/email/:email", userController.FindByEmail)
	}
}

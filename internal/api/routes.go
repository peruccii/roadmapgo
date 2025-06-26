package api

import (
	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/controller"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

func SetupRoutes(r *gin.Engine, userService services.UserService, authService services.AuthService, stripeService services.StripeService) {
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)
	stripeController := controller.NewStripeController(stripeService)

	r.POST("/login", authController.Login)

	r.POST("/webhook", stripeController.StripeWebhookController)

	users := r.Group("/users")
	{
		users.POST("", userController.Create)
		users.GET("/email/:email", userController.FindByEmail)
	}
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/api/middleware"
	"github.com/peruccii/roadmap-go-backend/internal/controller"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

func SetupRoutes(r *gin.Engine, userService services.UserService, authService services.AuthService, stripeService services.StripeService, robotService services.RobotService) {
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)
	stripeController := controller.NewStripeController(stripeService)
	robotController := controller.NewRobotController(robotService)

	r.POST("/login", authController.Login)

	r.POST("/webhook", stripeController.StripeWebhookController)

	users := r.Group("/users")
	{
		users.POST("", userController.Create)
		users.GET("/email/:email", userController.FindByEmail)
	}

	robots := r.Group("/robots").Use(middleware.AuthMiddleware(authService))
	{
		robots.POST("", robotController.Create)
		robots.GET("/:name", robotController.FindByName)
	}
}

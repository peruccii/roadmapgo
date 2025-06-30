package api

import (
	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/api/middleware"
	"github.com/peruccii/roadmap-go-backend/internal/controller"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
	"github.com/peruccii/roadmap-go-backend/internal/services"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	planRepo := repository.NewPlanRepository(db)
	robotRepo := repository.NewRobotRepository(db)

	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	planService := services.NewPlanService(planRepo)
	robotService := services.NewRobotService(robotRepo, planService)
	iaService := services.NewIAService()

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	robotController := controller.NewRobotController(robotService)
	conversaController := controller.NewConversaController(db, iaService)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.Create)
			auth.POST("/login", authController.Login)
		}

		robots := api.Group("/robots").Use(middleware.AuthMiddleware(authService))
		{
			robots.POST("", robotController.Create)
			robots.GET("/:name", robotController.FindByName)
			robots.POST("/:id/token", robotController.GenerateToken)
			robots.GET("", robotController.FindAll)
		}

		users := api.Group("/users").Use(middleware.AuthMiddleware(authService))
		{
			users.GET("", userController.FindAll)
		}

		api.POST("/conversa", middleware.RoboAuthMiddleware(authService), conversaController.Conversa)
	}

	return r
}

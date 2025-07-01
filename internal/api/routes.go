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

	// Repositórios
	userRepo := repository.NewUserRepository(db)
	planRepo := repository.NewPlanRepository(db)
	robotRepo := repository.NewRobotRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	// Serviços
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	planService := services.NewPlanService(planRepo)
	paymentService := services.NewPaymentService(paymentRepo, robotRepo)
	stripeService := services.NewStripeService(paymentRepo, subscriptionRepo, robotRepo, paymentService)
	robotService := services.NewRobotService(robotRepo, planService)
	iaService := services.NewIAService()

	// Controladores
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	robotController := controller.NewRobotController(robotService)
	paymentController := controller.NewPaymentController(stripeService, userService)
	stripeController := controller.NewStripeController(stripeService)
	conversaController := controller.NewConversaController(db, iaService)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.Create)
			auth.POST("/login", authController.Login)
		}

	// Grupo protegido por autenticação de usuário
	protected := api.Group("").Use(middleware.AuthMiddleware(authService))
	{
		// Endpoints de pagamento (substituem a criação direta de robôs)
		payments := protected.Group("/payments")
		{
			payments.POST("/robot", paymentController.CreateRobotPayment)
			payments.POST("/status", paymentController.CheckPaymentStatus)
		}

		// Endpoints de robôs (sem criação direta)
		robots := protected.Group("/robots")
		{
			robots.GET("/:name", robotController.FindByName)
			robots.POST("/:id/token", robotController.GenerateToken)
			robots.GET("", robotController.FindAll)
		}

		// Endpoints de usuários
		users := protected.Group("/users")
		{
			users.GET("", userController.FindAll)
		}
	}

	// Webhook do Stripe (sem autenticação)
	api.POST("/stripe/webhook", stripeController.StripeWebhookController)

	// Endpoint de conversa (protegido por autenticação de robô)
	api.POST("/conversa", middleware.RoboAuthMiddleware(authService, subscriptionRepo, robotRepo), conversaController.Conversa)
	}

	return r
}

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

type PaymentController interface {
	CreateRobotPayment(c *gin.Context)
	CheckPaymentStatus(c *gin.Context)
}

type paymentController struct {
	stripeService services.StripeService
	userService   services.UserService
}

func NewPaymentController(stripeService services.StripeService, userService services.UserService) PaymentController {
	return &paymentController{
		stripeService: stripeService,
		userService:   userService,
	}
}

type CreateRobotPaymentRequest struct {
	RobotName string `json:"robot_name" binding:"required"`
	PlanType  string `json:"plan_type" binding:"required"`
}

type PaymentStatusRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

func (ctrl *paymentController) CreateRobotPayment(c *gin.Context) {
	var req CreateRobotPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	user, err := ctrl.userService.FindByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user data"})
		return
	}

	validPlans := map[string]bool{
		"basic":      true,
		"premium":    true,
		"enterprise": true,
	}

	if !validPlans[req.PlanType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan type. Valid options: basic, premium, enterprise"})
		return
	}

	session, err := ctrl.stripeService.CreateCheckoutSessionForRobot(
		userID.(string),
		req.RobotName,
		req.PlanType,
		user.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment session: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":   session.ID,
		"checkout_url": session.URL,
		"message":      "Payment session created. Complete the payment to create your robot.",
	})
}

func (ctrl *paymentController) CheckPaymentStatus(c *gin.Context) {
	var req PaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// Aqui você pode implementar lógica para verificar o status do pagamento
	// Por enquanto, retornamos uma resposta simples
	c.JSON(http.StatusOK, gin.H{
		"message": "Check your robot status in the robots endpoint",
		"note":    "Payment status is updated automatically via webhooks",
	})
}

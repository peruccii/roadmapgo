package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

// RoboAuthMiddleware verifica um token JWT que representa um robô e se sua assinatura está ativa
func RoboAuthMiddleware(authService services.AuthService, subscriptionRepo repository.SubscriptionRepository, robotRepo repository.RobotRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format, must be Bearer."})
			c.Abort()
			return
		}

		claims, err := authService.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extrair robo_id do token
		roboIDStr, ok := claims["robo_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing 'robo_id' in token"})
			c.Abort()
			return
		}

		roboID, err := uuid.Parse(roboIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid robot ID format"})
			c.Abort()
			return
		}

		// Verificar se o robô existe e está ativo
		robot, err := robotRepo.FindById(roboID)
		if err != nil || robot == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Robot not found"})
			c.Abort()
			return
		}

		// Verificar se o robô está ativo
		if robot.Status != "active" {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Robot is not active. Please check your subscription"})
			c.Abort()
			return
		}

		// Verificar assinatura ativa
		subscription, err := subscriptionRepo.FindActiveByRobotID(roboID)
		if err != nil || subscription == nil || !subscription.IsActive() {
			// Se não há assinatura ativa, verificar se ainda está no período de validade do plano
			if robot.PlanValidUntil == nil || robot.PlanValidUntil.Before(time.Now()) {
				c.JSON(http.StatusPaymentRequired, gin.H{"error": "Subscription expired. Please renew your plan"})
				c.Abort()
				return
			}
		}

		// Se chegou até aqui, o robô está autorizado
		c.Set("robo_id", roboIDStr)
		c.Set("robot", robot)
		if subscription != nil {
			c.Set("subscription", subscription)
		}
		c.Next()
	}
}

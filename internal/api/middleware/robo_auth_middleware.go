package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

// RoboAuthMiddleware verifica um token JWT que representa um robô.
func RoboAuthMiddleware(authService services.AuthService) gin.HandlerFunc {
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

		// A mudança principal está aqui: procuramos por "robo_id"
		roboID, ok := claims["robo_id"].(float64) // O parser de JWT trata números como float64
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing 'robo_id' in token"})
			c.Abort()
			return
		}

		c.Set("robo_id", roboID) // Usamos a chave que o controller espera
		c.Next()
	}
}

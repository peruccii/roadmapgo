package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/dtos"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

type RobotController interface {
	Create(robot *models.Robot) error
	FindByName(name string) (*models.Robot, error)
}

type robotController struct {
	services services.RobotService
}

func (ctrl *robotController) Create(c *gin.Context) {
	var input services.CreateRobotInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	if err := ctrl.services.CreateRobot(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "[ robot ] - created")
}

func (ctrl *robotController) Active(c *gin.Context) {
	var req dtos.ActiveReqInputDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Invalid Data"})
		return
	}

	userID, ok := c.Get("usuario_id") // Pega o ID do usuário do JWT (via middleware)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "User not found"})
		return
	}

	existingRobot, err := ctrl.services.FindByName(req.Name)
	if err != nil {
	}

	if existingRobot == nil {
	}
}

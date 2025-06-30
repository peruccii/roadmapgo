package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/dtos"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

type robotController struct {
	services services.RobotService
}

func NewRobotController(service services.RobotService) RobotController {
	return &robotController{services: service}
}

type RobotController interface {
	Create(c *gin.Context)
	FindByName(c *gin.Context)
	GenerateToken(c *gin.Context)
	FindAll(c *gin.Context)
}

func (ctrl *robotController) FindAll(c *gin.Context) {
	robots, err := ctrl.services.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Converter para DTO de resposta
	response := make([]dtos.RobotResponseDTO, len(robots))
	for i, robot := range robots {
		response[i] = dtos.ConvertToRobotResponseDTO(robot)
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *robotController) GenerateToken(c *gin.Context) {
	robotID := c.Param("id")
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	token, err := ctrl.services.GenerateRobotToken(robotID, userID.(string))
	if err != nil {
		if err.Error() == "robot not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "plan expired" {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ctrl *robotController) Create(c *gin.Context) {
	var input services.CreateRobotInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	input.UserID = userID.(string)

	if err := ctrl.services.CreateRobot(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "[ robot ] - created")
}

func (ctrl *robotController) FindByName(c *gin.Context) {
	name := c.Param("name")
	robot, err := ctrl.services.FindByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if robot == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "robot not found"})
		return
	}

	c.JSON(http.StatusOK, robot)
}


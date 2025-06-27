package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
}

func (ctrl *robotController) Create(c *gin.Context) {
	var input *services.CreateRobotInput
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

<<<<<<< HEAD
	if robot == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "robot not found"})
		return
	}

	c.JSON(http.StatusOK, robot)
=======
	activate, err := ctrl.services.Active(existingRobot)
>>>>>>> 36a86502da1bde880b25e8ef2173dcb9fa6ff936
}

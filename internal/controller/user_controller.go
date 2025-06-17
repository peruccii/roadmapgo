package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

type UserController interface {
	Create(c *gin.Context)
	FindByEmail(c *gin.Context)
}

type userController struct {
	service services.UserService
}

func (ctrl *userController) Create(c *gin.Context) {
	var input services.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	if err := ctrl.service.CreateUser(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "[ user ] - created")
}

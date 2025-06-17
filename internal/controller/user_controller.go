package controller

import (
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
	}
}

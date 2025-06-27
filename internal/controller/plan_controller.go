package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/peruccii/roadmap-go-backend/internal/services"
)

type PlanController interface {
	Create(c *gin.Context)
}

type planController struct {
	service services.PlanService
}

func NewPlanController(service services.PlanService) PlanController {
	return &planController{service: service}
}

func (ctrl *planController) Create(c *gin.Context) {
}
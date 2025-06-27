package dtos

import "github.com/peruccii/roadmap-go-backend/internal/models"

type CreateRobotRequestInputDTO struct {
	 Robot  *models.Robot
	 *models.User
}

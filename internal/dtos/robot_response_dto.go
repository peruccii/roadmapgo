package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
)

type PlanResponseDTO struct {
	Type string `json:"type"`
}

type RobotResponseDTO struct {
	ID             uuid.UUID         `json:"ID"`
	Name           string            `json:"Name"`
	UserID         uuid.UUID         `json:"UserID"`
	User           *models.User      `json:"User"`
	ActivateIn     *time.Time        `json:"ActivateIn"`
	Status         models.RobotStatus `json:"Status"`
	PlanValidUntil *time.Time        `json:"PlanValidUntil"`
	LastPing       *time.Time        `json:"ultimo_ping"`
	CreatedAt      time.Time         `json:"CreatedAt"`
	Plans          []PlanResponseDTO `json:"Plans"`
}

// ConvertToRobotResponseDTO converte um modelo Robot para DTO de resposta
func ConvertToRobotResponseDTO(robot models.Robot) RobotResponseDTO {
	// Filtrar apenas planos ativos
	var activePlans []PlanResponseDTO
	for _, plan := range robot.Plans {
		if plan.Active {
			activePlans = append(activePlans, PlanResponseDTO{
				Type: string(plan.Type),
			})
		}
	}

	return RobotResponseDTO{
		ID:             robot.ID,
		Name:           robot.Name,
		UserID:         robot.UserID,
		User:           robot.User,
		ActivateIn:     robot.ActivateIn,
		Status:         robot.Status,
		PlanValidUntil: robot.PlanValidUntil,
		LastPing:       robot.LastPing,
		CreatedAt:      robot.CreatedAt,
		Plans:          activePlans,
	}
}

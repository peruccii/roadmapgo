package services

import (
	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
)

type PlanService interface {
	CreatePlan(robotID uuid.UUID, userID uuid.UUID) error
}

type planService struct {
	repo repository.PlanRepository
}

func NewPlanService(repo repository.PlanRepository) PlanService {
	return &planService{repo: repo}
}

func (s *planService) CreatePlan(robotID uuid.UUID, userID uuid.UUID) error {

	plan := &models.Plan{
		UserID:  userID,
		RobotID: robotID,
		Type:    models.BasicPlan,
		Active:  true,
	}

	return s.repo.CreatePlan(plan)
}
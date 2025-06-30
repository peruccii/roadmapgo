package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
)

type PlanService interface {
	CreatePlan(robotID uuid.UUID, userID uuid.UUID) error
	GetPlanByRobotID(robotID uuid.UUID) (*models.Plan, error)
}

type planService struct {
	repo repository.PlanRepository
}

func NewPlanService(repo repository.PlanRepository) PlanService {
	return &planService{repo: repo}
}

func (s *planService) GetPlanByRobotID(robotID uuid.UUID) (*models.Plan, error) {
	return s.repo.FindByRobotID(robotID)
}

func (s *planService) CreatePlan(robotID uuid.UUID, userID uuid.UUID) error {
	// Verificar se já existe um plano ativo para este robô
	existingPlan, err := s.repo.FindByRobotID(robotID)
	if err == nil && existingPlan != nil && existingPlan.Active {
		// Já existe um plano ativo, não criar outro
		return nil
	}

	// Desativar todos os planos antigos para este robô antes de criar um novo
	if err := s.repo.DeactivateOldPlans(robotID); err != nil {
		return err
	}

	plan := &models.Plan{
		UserID:    userID,
		RobotID:   robotID,
		Type:      models.BasicPlan,
		Active:    true,
		ExpiredIn: time.Now().Add(time.Hour * 24 * 30),
	}

	return s.repo.CreatePlan(plan)
}

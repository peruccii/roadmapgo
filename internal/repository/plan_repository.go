package repository

import (
	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"gorm.io/gorm"
)

type PlanRepository interface {
	CreatePlan(plan *models.Plan) error
	FindByRobotID(robotID uuid.UUID) (*models.Plan, error)
	DeactivateOldPlans(robotID uuid.UUID) error
}

func (r *planRepository) FindByRobotID(robotID uuid.UUID) (*models.Plan, error) {
	var plan models.Plan
	if err := r.db.Where("robot_id = ? AND active = ?", robotID, true).First(&plan).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

type planRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) PlanRepository {
	return &planRepository{db: db}
}

func (r *planRepository) CreatePlan(plan *models.Plan) error {
	return r.db.Create(plan).Error
}

func (r *planRepository) DeactivateOldPlans(robotID uuid.UUID) error {
	return r.db.Model(&models.Plan{}).Where("robot_id = ?", robotID).Update("active", false).Error
}

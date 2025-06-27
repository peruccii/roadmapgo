package repository

import (
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"gorm.io/gorm"
)

type PlanRepository interface {
	CreatePlan(plan *models.Plan) error
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
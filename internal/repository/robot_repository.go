package repository

import (
	"errors"

	"github.com/peruccii/roadmap-go-backend/internal/dtos"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"gorm.io/gorm"
)

type robotRepository struct{ db *gorm.DB }

<<<<<<< HEAD
func NewRobotRepository(db *gorm.DB) RobotRepository {
	return &robotRepository{db: db}
}
=======

>>>>>>> 36a86502da1bde880b25e8ef2173dcb9fa6ff936

type RobotRepository interface {
	Create(robot *models.Robot) error
	FindByName(name string) (*models.Robot, error)
}

func (r *robotRepository) Active(input *dtos.ActiveReqInputDTO) *models.Robot {
	var robo models.Robot
	if err := r.db.First(&robo, "id = ?", input.DeviceID).Error; err != nil {
		return &models.Robot{}
	}

	return &robo
}

func (r *robotRepository) FindByName(name string) (*models.Robot, error) {
	var robot models.Robot
	if err := r.db.Where("name = ?", name).First(&robot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &robot, nil
}

func (r *robotRepository) Create(robot *models.Robot) error {
	tx := r.db.Begin()
	result := tx.Create(robot)

	err := result.Error
	if err != nil {
		tx.Rollback() // transactional cancelled
		return errors.New("failed to create user:" + err.Error())
	}

	return tx.Commit().Error
}

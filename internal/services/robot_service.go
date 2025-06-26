package services

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
)

type CreateRobotInput struct {
	Name string
}

type robotRepository struct {
	repo repository.RobotRepository
}

type RobotService interface {
	CreateRobot(input CreateRobotInput) error
	FindByName(name string) (*models.Robot, error)
}

func (r *robotRepository) FindByName(name string) (*models.Robot, error) {
	return r.repo.FindByName(name)
}

func (r *robotRepository) CreateRobot(input CreateRobotInput) error {
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return errors.New("invalid input" + err.Error())
	}

	existingRobot, err := r.repo.FindByName(input.Name)
	if err != nil {
		return err
	}

	if existingRobot == nil {
		return errors.New("robot already exist")
	}

	return nil
}

package services

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
)

type CreateRobotInput struct {
	Name   string
	UserID string
}

type robotService struct {
	repo repository.RobotRepository
	planService PlanService
}

func NewRobotService(repo repository.RobotRepository, planService PlanService) RobotService {
	return &robotService{repo: repo, planService: planService}
}

type RobotService interface {
	CreateRobot(input CreateRobotInput) error
	FindByName(name string) (*models.Robot, error)
}

func (r *robotService) FindByName(name string) (*models.Robot, error) {
	return r.repo.FindByName(name)
}

func (r *robotService) CreateRobot(input CreateRobotInput) error {
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return errors.New("invalid input" + err.Error())
	}

	existingRobot, err := r.repo.FindByName(input.Name)
	if err != nil {
		return err
	}

	if existingRobot != nil {
		return errors.New("robot already exist")
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return errors.New("invalid user id")
	}

	robot := &models.Robot{
		Name:   input.Name,
		UserID: userID,
	}

	if err := r.repo.Create(robot); err != nil {
		return err
	}

	if err := r.planService.CreatePlan(robot.ID, userID); err != nil {
		return err
	}

	return nil
}
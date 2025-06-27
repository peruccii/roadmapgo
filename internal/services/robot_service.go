package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/dtos"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
	"github.com/peruccii/roadmap-go-backend/internal/utils"
)

type CreateRobotInput struct {
	Name string
	UserID uuid.UUID
}

type robotRepository struct {
	repo repository.RobotRepository
	repoPlan repository.PlanRepository
}

type RobotService interface {
	CreateRobot(input *CreateRobotInput) error
	FindByName(name string) (*models.Robot, error)
}

func (r *robotRepository) ActiveRobot(input dtos.ActiveReqInputDTO) error {
	if err := utils.ValidateFields(input); err != nil {
		return err
	}

	return nil
}

func (r *robotRepository) FindByName(name string) (*models.Robot, error) {
	return r.repo.FindByName(name)
}

func (r *robotRepository) CreateRobot(input *CreateRobotInput) error {
	if err := utils.ValidateFields(input); err != nil {
		return err
	}

	existingRobot, err := r.repo.FindByName(input.Name)
	if err != nil {
		return err
	}

	if existingRobot == nil {
		return errors.New("robot already exist")
	}

	robot := &models.Robot{
		Name:     input.Name,
		UserID:   input.UserID ,
	}

	if err := r.repo.Create(robot); err != nil {

	}

	plan := &models.Plan{
		ID:        uuid.New(),
		UserID:    input.UserID,
		RobotID:   robot.ID,
		Type:      models.BasicPlan,
		InitiateIn: time.Now(),
		ExpiredIn:  time.Now().Add(30 * 24 * time.Hour),
		Active:     true,
		PaymentID:  "",
	}

	if err := r.repoPlan.Create(plan); err != nil {

	}

	return nil
}

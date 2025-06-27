package services

import (
	"errors"
	"time"

<<<<<<< HEAD
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
=======
	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/dtos"
>>>>>>> 36a86502da1bde880b25e8ef2173dcb9fa6ff936
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
	"github.com/peruccii/roadmap-go-backend/internal/utils"
)

type CreateRobotInput struct {
<<<<<<< HEAD
	Name   string
	UserID string
}

type robotService struct {
	repo        repository.RobotRepository
	planService PlanService
}

func NewRobotService(repo repository.RobotRepository, planService PlanService) RobotService {
	return &robotService{repo: repo, planService: planService}
=======
	Name string
	UserID uuid.UUID
}

type robotRepository struct {
	repo repository.RobotRepository
	repoPlan repository.PlanRepository
>>>>>>> 36a86502da1bde880b25e8ef2173dcb9fa6ff936
}

type RobotService interface {
	CreateRobot(input *CreateRobotInput) error
	FindByName(name string) (*models.Robot, error)
}

<<<<<<< HEAD
func (r *robotService) FindByName(name string) (*models.Robot, error) {
	return r.repo.FindByName(name)
}

func (r *robotService) CreateRobot(input CreateRobotInput) error {
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return errors.New("invalid input" + err.Error())
=======
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
>>>>>>> 36a86502da1bde880b25e8ef2173dcb9fa6ff936
	}

	existingRobot, err := r.repo.FindByName(input.Name)
	if err != nil {
		return err
	}

	if existingRobot != nil {
		return errors.New("robot already exist")
	}

<<<<<<< HEAD
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
=======
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

>>>>>>> 36a86502da1bde880b25e8ef2173dcb9fa6ff936
	}

	return nil
}

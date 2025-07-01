package services

import (
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
)

type CreateRobotInput struct {
	Name   string
	UserID string
}

type robotService struct {
	repo        repository.RobotRepository
	planService PlanService
	secretKey   []byte
}

func NewRobotService(repo repository.RobotRepository, planService PlanService) RobotService {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		secret = "default-secret"
	}
	return &robotService{
		repo:        repo,
		planService: planService,
		secretKey:   []byte(secret),
	}
}

type RobotService interface {
	CreateRobot(input CreateRobotInput) error
	FindByName(name string) (*models.Robot, error)
	GenerateRobotToken(robotID, userID string) (string, error)
	FindAll() ([]models.Robot, error)
}

func (r *robotService) FindAll() ([]models.Robot, error) {
	robots, err := r.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Atualizar PlanValidUntil para cada robô baseado nos planos ativos
	for i := range robots {
		r.updateRobotPlanValidUntil(&robots[i])
	}

	return robots, nil
}

func (s *robotService) GenerateRobotToken(robotID, userID string) (string, error) {
	robot, err := s.repo.FindByIDAndUserID(robotID, userID)
	if err != nil {
		return "", err
	}
	if robot == nil {
		return "", errors.New("robot not found")
	}

	parsedRobotID, err := uuid.Parse(robotID)
	if err != nil {
		return "", errors.New("invalid robot id")
	}

	plan, err := s.planService.GetPlanByRobotID(parsedRobotID)
	if err != nil {
		return "", err
	}

	if plan.ExpiredIn.Before(time.Now()) {
		return "", errors.New("plan expired")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"robo_id": robot.ID,
			"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

	return token.SignedString(s.secretKey)
}

func (r *robotService) FindByName(name string) (*models.Robot, error) {
	return r.repo.FindByName(name)
}

// CreateRobot agora não pode criar robô diretamente - deve ser feito através do pagamento
func (r *robotService) CreateRobot(input CreateRobotInput) error {
	return errors.New("criação de robô deve ser feita através do pagamento. Use o endpoint de pagamento")
}

// CreateRobotWithPayment cria robô após confirmação de pagamento (usado internamente)
func (r *robotService) CreateRobotWithPayment(input CreateRobotInput, planValidUntil time.Time) error {
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
		Name:           input.Name,
		UserID:         userID,
		Status:         models.StatusActive, // Ativo após pagamento
		PlanValidUntil: &planValidUntil,
	}

	if err := r.repo.Create(robot); err != nil {
		return err
	}

	if err := r.planService.CreatePlan(robot.ID, userID); err != nil {
		return err
	}

	return nil
}

// updateRobotPlanValidUntil atualiza o campo PlanValidUntil do robô baseado nos planos ativos
func (r *robotService) updateRobotPlanValidUntil(robot *models.Robot) {
	if len(robot.Plans) == 0 {
		robot.PlanValidUntil = nil
		return
	}

	// Encontrar o plano ativo com data de expiração mais recente
	var latestExpiration *time.Time
	for _, plan := range robot.Plans {
		if plan.Active && (latestExpiration == nil || plan.ExpiredIn.After(*latestExpiration)) {
			latestExpiration = &plan.ExpiredIn
		}
	}

	robot.PlanValidUntil = latestExpiration
}

package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
)

type PaymentService interface {
	CreatePayment(payment *models.Payment) error
	HandlePaymentSuccess(paymentID, sessionID string) error
	HandlePaymentFailure(paymentID string) error
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	robotRepo   repository.RobotRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, robotRepo repository.RobotRepository) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		robotRepo:   robotRepo,
	}
}

func (s *paymentService) CreatePayment(payment *models.Payment) error {
	return s.paymentRepo.Create(payment)
}

func (s *paymentService) HandlePaymentSuccess(paymentID, sessionID string) error {
	payment, err := s.paymentRepo.FindByProviderSessionID(sessionID)
	if err != nil || payment == nil {
		return fmt.Errorf("payment not found for session: %s", sessionID)
	}

	payment.Status = models.PaymentCompleted
	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	// Ativar Robot se necessário
	if payment.RobotID != nil {
		robot, err := s.robotRepo.FindById(*payment.RobotID)
		if err == nil && robot != nil {
			robot.Status = models.StatusActive
			return s.robotRepo.Update(robot)
		}
	}

	return nil
}

func (s *paymentService) HandlePaymentFailure(paymentID string) error {
	id, err := uuid.Parse(paymentID)
	if err != nil {
		return fmt.Errorf("invalid payment ID: %s", paymentID)
	}

	payment, err := s.paymentRepo.FindByID(id)
	if err != nil || payment == nil {
		return fmt.Errorf("payment not found: %s", paymentID)
	}

	payment.Status = models.PaymentFailed
	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	// Suspender Robot se necessário
	if payment.RobotID != nil {
		robot, err := s.robotRepo.FindById(*payment.RobotID)
		if err == nil && robot != nil {
			robot.Status = models.StatusSuspense
			return s.robotRepo.Update(robot)
		}
	}

	return nil
}

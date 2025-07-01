package repository

import (
	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByID(id uuid.UUID) (*models.Payment, error)
	FindByProviderPaymentID(providerPaymentID string) (*models.Payment, error)
	FindByProviderSessionID(sessionID string) (*models.Payment, error)
	FindByUserID(userID uuid.UUID) ([]models.Payment, error)
	FindByRobotID(robotID uuid.UUID) ([]models.Payment, error)
	UpdateStatus(id uuid.UUID, status models.PaymentStatus) error
	Update(payment *models.Payment) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("User").Preload("Robot").Preload("Plan").First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByProviderPaymentID(providerPaymentID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("User").Preload("Robot").Preload("Plan").First(&payment, "provider_payment_id = ?", providerPaymentID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByProviderSessionID(sessionID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("User").Preload("Robot").Preload("Plan").First(&payment, "provider_session_id = ?", sessionID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByUserID(userID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Preload("User").Preload("Robot").Preload("Plan").Where("user_id = ?", userID).Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) FindByRobotID(robotID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Preload("User").Preload("Robot").Preload("Plan").Where("robot_id = ?", robotID).Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) UpdateStatus(id uuid.UUID, status models.PaymentStatus) error {
	return r.db.Model(&models.Payment{}).Where("id = ?", id).Update("status", status).Error
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

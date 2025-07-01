package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	FindByID(id uuid.UUID) (*models.Subscription, error)
	FindByRobotID(robotID uuid.UUID) (*models.Subscription, error)
	FindByUserID(userID uuid.UUID) ([]models.Subscription, error)
	FindByProviderSubscriptionID(providerSubscriptionID string) (*models.Subscription, error)
	FindActiveByRobotID(robotID uuid.UUID) (*models.Subscription, error)
	UpdateStatus(id uuid.UUID, status models.SubscriptionStatus) error
	Update(subscription *models.Subscription) error
	FindExpiringSubscriptions(days int) ([]models.Subscription, error)
	CancelSubscription(id uuid.UUID, cancelAtPeriodEnd bool) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *subscriptionRepository) FindByID(id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("User").Preload("Robot").Preload("Payments").First(&subscription, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) FindByRobotID(robotID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("User").Preload("Robot").Preload("Payments").
		Where("robot_id = ?", robotID).
		Order("created_at DESC").
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) FindByUserID(userID uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Preload("User").Preload("Robot").Preload("Payments").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) FindByProviderSubscriptionID(providerSubscriptionID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Preload("User").Preload("Robot").Preload("Payments").
		First(&subscription, "provider_subscription_id = ?", providerSubscriptionID).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) FindActiveByRobotID(robotID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	now := time.Now()
	err := r.db.Preload("User").Preload("Robot").Preload("Payments").
		Where("robot_id = ? AND status = ? AND current_period_start <= ? AND current_period_end > ?", 
			robotID, models.SubscriptionActive, now, now).
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) UpdateStatus(id uuid.UUID, status models.SubscriptionStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	if status == models.SubscriptionCanceled {
		updates["canceled_at"] = time.Now()
	}
	
	return r.db.Model(&models.Subscription{}).Where("id = ?", id).Updates(updates).Error
}

func (r *subscriptionRepository) Update(subscription *models.Subscription) error {
	return r.db.Save(subscription).Error
}

func (r *subscriptionRepository) FindExpiringSubscriptions(days int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	cutoffDate := time.Now().Add(time.Duration(days) * 24 * time.Hour)
	
	err := r.db.Preload("User").Preload("Robot").
		Where("status = ? AND current_period_end <= ? AND cancel_at_period_end = ?", 
			models.SubscriptionActive, cutoffDate, false).
		Find(&subscriptions).Error
	
	return subscriptions, err
}

func (r *subscriptionRepository) CancelSubscription(id uuid.UUID, cancelAtPeriodEnd bool) error {
	updates := map[string]interface{}{
		"cancel_at_period_end": cancelAtPeriodEnd,
	}
	
	if !cancelAtPeriodEnd {
		updates["status"] = models.SubscriptionCanceled
		updates["canceled_at"] = time.Now()
	}
	
	return r.db.Model(&models.Subscription{}).Where("id = ?", id).Updates(updates).Error
}

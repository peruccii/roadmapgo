package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionStatus representa o status de uma assinatura
type SubscriptionStatus string

const (
	SubscriptionActive   SubscriptionStatus = "active"
	SubscriptionInactive SubscriptionStatus = "inactive"
	SubscriptionCanceled SubscriptionStatus = "canceled"
	SubscriptionExpired  SubscriptionStatus = "expired"
	SubscriptionPending  SubscriptionStatus = "pending"
)

// Subscription representa uma assinatura de plano
type Subscription struct {
	ID                     uuid.UUID          `gorm:"type:uuid;primaryKey"`
	UserID                 uuid.UUID          `gorm:"type:uuid;not null"`
	User                   User               `gorm:"foreignKey:UserID"`
	RobotID                uuid.UUID          `gorm:"type:uuid;not null"`
	Robot                  Robot              `gorm:"foreignKey:RobotID"`
	PlanType               PlanType           `gorm:"type:text;not null"`
	Status                 SubscriptionStatus `gorm:"type:text;default:'pending'"`
	CurrentPeriodStart     time.Time          `gorm:"not null"`
	CurrentPeriodEnd       time.Time          `gorm:"not null"`
	ProviderSubscriptionID string             `gorm:"type:varchar(255)"` // Stripe subscription ID
	ProviderCustomerID     string             `gorm:"type:varchar(255)"` // Stripe customer ID
	CancelAtPeriodEnd      bool               `gorm:"default:false"`
	CanceledAt             *time.Time
	CreatedAt              time.Time `gorm:"autoCreateTime"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"`

	// Relacionamentos
	Payments []Payment `gorm:"foreignKey:ProviderSubscriptionID;references:ProviderSubscriptionID"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}

// IsActive verifica se a assinatura está ativa e dentro do período válido
func (s *Subscription) IsActive() bool {
	now := time.Now()
	return s.Status == SubscriptionActive && 
		   now.After(s.CurrentPeriodStart) && 
		   now.Before(s.CurrentPeriodEnd)
}

// ShouldRenew verifica se a assinatura deve ser renovada
func (s *Subscription) ShouldRenew() bool {
	return s.Status == SubscriptionActive && 
		   !s.CancelAtPeriodEnd && 
		   time.Now().After(s.CurrentPeriodEnd.Add(-24*time.Hour)) // 1 dia antes do vencimento
}

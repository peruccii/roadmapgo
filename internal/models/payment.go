package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentStatus representa o status de um pagamento
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
	PaymentCanceled  PaymentStatus = "canceled"
	PaymentRefunded  PaymentStatus = "refunded"
)

// PaymentProvider representa o provedor de pagamento
type PaymentProvider string

const (
	ProviderStripe PaymentProvider = "stripe"
)

// Payment representa um pagamento no sistema
type Payment struct {
	ID                   uuid.UUID       `gorm:"type:uuid;primaryKey"`
	UserID               uuid.UUID       `gorm:"type:uuid;not null"`
	User                 User            `gorm:"foreignKey:UserID"`
	RobotID              *uuid.UUID      `gorm:"type:uuid"`
	Robot                *Robot          `gorm:"foreignKey:RobotID"`
	PlanID               *uuid.UUID      `gorm:"type:uuid"`
	Plan                 *Plan           `gorm:"foreignKey:PlanID"`
	Amount               int64           `gorm:"not null"` // em centavos
	Currency             string          `gorm:"type:varchar(3);default:'BRL'"`
	Status               PaymentStatus   `gorm:"type:text;default:'pending'"`
	Provider             PaymentProvider `gorm:"type:text;default:'stripe'"`
	ProviderPaymentID    string          `gorm:"type:varchar(255)"` // ID do Stripe
	ProviderCustomerID   string          `gorm:"type:varchar(255)"` // Customer ID do Stripe
	ProviderSessionID    string          `gorm:"type:varchar(255)"` // Session ID do Stripe
	ProviderSubscriptionID string        `gorm:"type:varchar(255)"` // Subscription ID do Stripe
	Metadata             string          `gorm:"type:text"` // JSON com dados extras
	CreatedAt            time.Time       `gorm:"autoCreateTime"`
	UpdatedAt            time.Time       `gorm:"autoUpdateTime"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}

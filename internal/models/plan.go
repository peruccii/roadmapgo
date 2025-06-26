package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanType string

const (
	BasicPlan PlanType = "basic"
)

type Plan struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	User       User      `gorm:"foreignKey:UserID"`
	RobotID    uuid.UUID `gorm:"type:uuid;not null"`
	Robot      Robot     `gorm:"foreignKey:RobotID"`
	Type       PlanType  `gorm:"type:text;not null"`
	InitiateIn time.Time `gorm:"autoCreateTime"`
	ExpiredIn  time.Time
	Active     bool   `gorm:"default:true"`
	PaymentID  string // ID vindo do Stripe, MercadoPago, etc
}

func (p *Plan) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}

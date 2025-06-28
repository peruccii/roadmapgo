package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RobotStatus string

const (
	StatusPENDING  RobotStatus = "pending"
	StatusActive   RobotStatus = "active"
	StatusSuspense RobotStatus = "suspense"
)

type Robot struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"` // device_id
	Name           string
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	User           *User      `gorm:"foreignKey:UserID"`
	ActivateIn     *time.Time
	Status         RobotStatus `gorm:"type:text;default:'pending'"`
	PlanValidUntil *time.Time
	LastPing       *time.Time `json:"ultimo_ping"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`

	Plans []Plan `gorm:"foreignKey:RobotID"`
}

func (r *Robot) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}

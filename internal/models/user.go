package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name" db:"name" gorm:"type:varchar(255);not null"`
	Email     string    `json:"email" db:"email" gorm:"type:varchar(255);unique;not null"`
	Password  string    `json:"-" db:"password" gorm:"type:varchar(255);not null"` // hash, n exposto no JSON
	Robots    []Robot  `json:"robots" gorm:"manytoone:user_robots;"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`

	Robot []Robot `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

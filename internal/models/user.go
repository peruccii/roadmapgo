package models

import (
	"time"
)

// struct

type User struct {
	ID        int64     `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" db:"name" gorm:"type:varchar(255);not null"`
	Email     string    `json:"email" db:"email" gorm:"type:varchar(255);unique;not null"`
	Password  string    `json:"-" db:"password" gorm:"type:varchar(255);not null"` // hash, n exposto no JSON
	Contents  []Courses `json:"contents" gorm:"many2many:user_courses;"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

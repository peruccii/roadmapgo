package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UserOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserInputDTO struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=255"`
	Email    string `json:"email" validate:"omitempty,email,max=255"`
	Password string `json:"password" validate:"omitempty,min=4,max=255"`
}

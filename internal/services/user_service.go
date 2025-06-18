package services

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/peruccii/roadmap-go-backend/internal/dtos"
	"github.com/peruccii/roadmap-go-backend/internal/models"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Name     string `validate:"required,min=2,max=255"`
	Email    string `validate:"required,email,max=255"`
	Password string `validate:"required,min=8"`
}

type UserService interface {
	CreateUser(input UserInput) error
	FindByEmail(email string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func (s *userService) CreateUser(input UserInput) error {
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return errors.New("invalid input" + err.Error())
	}

	existingUser, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password" + err.Error())
	}

	user := &models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}
	return nil
}

func (s *userService) FindByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(email)
}

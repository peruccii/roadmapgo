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
	Delete(param *repository.DeleteUserParams) error
	Update(ID int64, input *dtos.UpdateUserInputDTO) error
	FindAll() ([]models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *userService) FindAll() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) Update(ID int64, input *dtos.UpdateUserInputDTO) error {
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return errors.New("invalid input" + err.Error())
	}

	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	// verificar se o novo email ( se fornecido ) já esta em uso por outro usuario

	if input.Email != "" && input.Email != user.Email {
		existingUser, err := s.repo.FindByEmail(input.Email)
		if err != nil {
			return err
		}

		if existingUser != nil {
			return errors.New("email already exists")
		}
	}

	// todo

	if input.Email != "" {
		user.Email = input.Email
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	user.UpdatedAt = time.Now()

	if input.Password != "" {
		hashPassword, err := hashPassword(input.Password)
		if err != nil {
			return errors.New("failed to hash password" + err.Error())
		}
		user.Password = string(hashPassword)
	}

	return nil
}

func (s *userService) Delete(params *repository.DeleteUserParams) error {
	return s.repo.Delete(params)
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

	hashPassword, err := hashPassword(input.Password)
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

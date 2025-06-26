package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/peruccii/roadmap-go-backend/internal/dtos"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	AuthUser(params dtos.AuthInputDTO) (dtos.AuthOutputDTO, error)
}

type authService struct {
	repo      repository.UserRepository
	secretKey []byte
}

func NewAuthService(repo repository.UserRepository) AuthService {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		secret = "default-secret" // Fallback for development
	}
	return &authService{
		repo:      repo,
		secretKey: []byte(secret),
	}
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *authService) createToken(userId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userId,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	return token.SignedString(s.secretKey)
}

func (s *authService) AuthUser(params dtos.AuthInputDTO) (dtos.AuthOutputDTO, error) {
	existingUser, err := s.repo.FindByEmail(params.Email)
	if err != nil {
		return dtos.AuthOutputDTO{}, errors.New("failed to find user: " + err.Error())
	}

	if existingUser == nil {
		return dtos.AuthOutputDTO{}, errors.New("user not found")
	}

	if !CheckPasswordHash(params.Password, existingUser.Password) {
		return dtos.AuthOutputDTO{}, errors.New("invalid password")
	}

	tokenString, err := s.createToken(existingUser.ID)
	if err != nil {
		return dtos.AuthOutputDTO{}, errors.New("failed to generate token")
	}

	response := dtos.AuthOutputDTO{
		AccessToken: tokenString,
	}

	return response, nil
}

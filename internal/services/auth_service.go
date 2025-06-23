package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/peruccii/roadmap-go-backend/internal/dtos"
	"github.com/peruccii/roadmap-go-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type authParams struct {
	Email    string
	Password string
}

type AuthService interface {
	AuthUser(params authParams) (dtos.AuthOutputDTO, error)
}

type uService struct {
	repo repository.UserRepository
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var secretKey = []byte("sadkfn72!")

func createToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userId,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *uService) AuthUser(params authParams) (dtos.AuthOutputDTO, error) {
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

	tokenString, _ := createToken(string(existingUser.ID))

	response := dtos.AuthOutputDTO{
		AccessToken: tokenString,
	}

	return response, nil
}

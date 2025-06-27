package repository

import (
	"errors"
	"fmt"

	"github.com/peruccii/roadmap-go-backend/internal/models"
	"gorm.io/gorm"
)

type userRepository struct{ db *gorm.DB }

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Delete(params *DeleteUserParams) error
}

type DeleteUserParams struct {
	Email string
	ID    string
}

func (r *userRepository) Delete(params *DeleteUserParams) error {
	if params.Email == "" && params.ID == "" {
		return fmt.Errorf("at least one of email or ID must be provided")
	}

	query := r.db

	if params.Email != "" {
		query = query.Where("email = ?", params.Email)
	} else if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	}

	if err := query.Delete(&models.User{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // verifica se o erro foi porque nenhum registro foi encontrado
			return nil, nil // sim
		}
		return nil, err // nao
	}
	return &user, nil
}

// First() -> retorna os dados se encontrados e popula a variavel user com os dados encontrados

// o Método associado ao strct ( r == this. )
func (r *userRepository) Create(user *models.User) error {
	tx := r.db.Begin() // starts a manual @Transactional
	result := tx.Create(user)

	err := result.Error
	if err != nil {
		tx.Rollback() // transactional cancelled
		return errors.New("failed to create user:" + err.Error())

	}

	return tx.Commit().Error // persisting in database ( if persist Error ( returned ) )
}

package repository

import (
	"errors"

	"github.com/peruccii/roadmap-go-backend/internal/models"
	"gorm.io/gorm"
)

type userRepository struct{ db *gorm.DB }

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
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

// Método associado ao struct ( r == this. )
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

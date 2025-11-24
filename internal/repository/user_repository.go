package repository

import (
	"time"

	"github.com/yourusername/pokemon-chatbot-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByTelegramID(telegramID string) (*models.User, error)
	UpdateLastActive(telegramID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByTelegramID(telegramID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateLastActive(telegramID string) error {
	return r.db.Model(&models.User{}).
		Where("telegram_id = ?", telegramID).
		Update("last_active", time.Now()).Error
}

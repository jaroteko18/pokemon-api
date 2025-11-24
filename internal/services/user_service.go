package services

import (
	"errors"
	"time"

	"github.com/yourusername/pokemon-chatbot-api/internal/models"
	"github.com/yourusername/pokemon-chatbot-api/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	Register(telegramID, firstName, lastName, username string) (*RegisterResponse, error)
	GetUserByTelegramID(telegramID string) (*models.User, error)
	IsUserRegistered(telegramID string) (bool, error)
}

type userService struct {
	repo repository.UserRepository
}

type RegisterResponse struct {
	Success bool         `json:"success"`
	Exists  bool         `json:"exists"`
	Message string       `json:"message,omitempty"`
	User    *models.User `json:"user,omitempty"`
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(telegramID, firstName, lastName, username string) (*RegisterResponse, error) {
	// Check if user exists
	existingUser, err := s.repo.FindByTelegramID(telegramID)
	if err == nil && existingUser != nil {
		return &RegisterResponse{
			Success: false,
			Exists:  true,
			Message: "User already registered",
			User:    existingUser,
		}, nil
	}

	// Create new user
	now := time.Now()
	user := &models.User{
		TelegramID:   telegramID,
		FirstName:    firstName,
		LastName:     lastName,
		Username:     username,
		RegisteredAt: now,
		LastActive:   now,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		Success: true,
		Exists:  false,
		User:    user,
	}, nil
}

func (s *userService) GetUserByTelegramID(telegramID string) (*models.User, error) {
	return s.repo.FindByTelegramID(telegramID)
}

func (s *userService) IsUserRegistered(telegramID string) (bool, error) {
	user, err := s.repo.FindByTelegramID(telegramID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return user != nil, nil
}

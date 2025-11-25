package repository

import (
	"encoding/json"
	"fmt"

	"github.com/yourusername/pokemon-chatbot-api/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByTelegramID(telegramID string) (*models.User, error)
	FindAll() ([]models.User, error)
	FindAllPaginated(page, limit int) ([]models.User, error)
	Count() (int, error)
	UpdateLastActive(telegramID string) error
}

type userRepository struct {
	client *SupabaseClient
}

func NewUserRepository(supabaseURL, supabaseKey string) UserRepository {
	return &userRepository{
		client: NewSupabaseClient(supabaseURL, supabaseKey),
	}
}

func (r *userRepository) Create(user *models.User) error {
	body, err := r.client.Insert("users", user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	var results []models.User
	if err := json.Unmarshal(body, &results); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(results) > 0 {
		*user = results[0]
	}

	return nil
}

func (r *userRepository) FindByTelegramID(telegramID string) (*models.User, error) {
	body, err := r.client.Select("users", "telegram_id", telegramID)
	if err != nil {
		return nil, err
	}

	var results []models.User
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &results[0], nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	body, err := r.client.SelectAll("users")
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return users, nil
}

func (r *userRepository) FindAllPaginated(page, limit int) ([]models.User, error) {
	body, err := r.client.SelectAllPaginated("users", page, limit)
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return users, nil
}

func (r *userRepository) Count() (int, error) {
	users, err := r.FindAll()
	if err != nil {
		return 0, err
	}
	return len(users), nil
}

func (r *userRepository) UpdateLastActive(telegramID string) error {
	updates := map[string]interface{}{
		"last_active": "now()",
	}

	_, err := r.client.Update("users", "telegram_id", telegramID, updates)
	return err
}

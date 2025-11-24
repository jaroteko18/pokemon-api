package config

import (
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

type Config struct {
	SupabaseURL string
	SupabaseKey string
	Port        string
}

func New() *Config {
	return &Config{
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		SupabaseKey: os.Getenv("SUPABASE_KEY"),
		Port:        os.Getenv("PORT"),
	}
}

func NewSupabaseClient(cfg *Config) (*supabase.Client, error) {
	client, err := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}
	return client, nil
}

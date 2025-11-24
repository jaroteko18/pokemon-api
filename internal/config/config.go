package config

import (
	"os"
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

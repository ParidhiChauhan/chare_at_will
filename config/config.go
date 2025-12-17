package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RzpKey    string
	RzpSecret string
	Port      string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		RzpKey:    os.Getenv("RZP_KEY"),
		RzpSecret: os.Getenv("RZP_SECRET"),
		Port:      os.Getenv("PORT"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.RzpKey == "" || cfg.RzpSecret == "" {
		log.Fatal("RZP_KEY or RZP_SECRET missing in environment")
	}

	return cfg
}

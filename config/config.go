package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	RzpKey    string
	RzpSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:      os.Getenv("PORT"),
		RzpKey:    os.Getenv("RAZORPAY_KEY_ID"),
		RzpSecret: os.Getenv("RAZORPAY_KEY_SECRET"),
	}
}

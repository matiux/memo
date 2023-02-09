package application

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() {
	localEnvFile := ".env.local"

	if _, err := os.Stat(localEnvFile); err == nil {
		_ = godotenv.Load(localEnvFile)
	}

	_ = godotenv.Load()
}

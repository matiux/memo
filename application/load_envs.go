package application

import (
	"github.com/joho/godotenv"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func LoadEnv() {

	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := path.Dir(currentFile)
	rootPath := filepath.Join(currentDir, "..")

	localEnvFile := rootPath + "/.env.local"

	if _, err := os.Stat(localEnvFile); err == nil {
		_ = godotenv.Load(localEnvFile)
	}

	_ = godotenv.Load()
}

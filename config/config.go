package config

import (
	"os"
)

type Config struct {
	StorageType string
	DSN         string
}

func LoadConfig() *Config {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "memory"
	}

	dsn := os.Getenv("POSTGRES_DSN")

	return &Config{
		StorageType: storageType,
		DSN:         dsn,
	}
}

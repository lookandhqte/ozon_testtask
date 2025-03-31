package config

import (
	"os"
)

type Config struct {
	StorageType string
	PostgresDSN string // Строка подключения к рostgreSQL
}

func LoadConfig() *Config {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "memory" // по умолчанию
	}

	postgresDSN := os.Getenv("POSTGRES_DSN")

	return &Config{
		StorageType: storageType,
		PostgresDSN: postgresDSN,
	}
}

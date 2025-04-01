package config

import "os"

type Config struct {
	StorageType string
	DSN         string
}

func LoadConfig() *Config {
	config := &Config{
		StorageType: getEnvOrDefault("STORAGE_TYPE", "memory"),
		DSN:         os.Getenv("POSTGRES_DSN"),
	}
	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

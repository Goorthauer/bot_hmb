package config

import (
	"os"
)

type Config struct {
	Db                 string
	RedisAddr          string
	Token              string
	TelegramURL        string
	MasterUserNickname string
	UserEncryptKey     string
	Debug              bool
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Debug:              getBoolEnv("DEBUG_MODE"),
		RedisAddr:          getEnv("REDIS_DSN", ""),
		MasterUserNickname: getEnv("MASTER_USER", ""),
		TelegramURL:        getEnv("URL", ""),
		Db:                 getEnv("DATABASE", ""),
		Token:              getEnv("TOKEN", ""),
		UserEncryptKey:     getEnv("USER_ENCRYPT_KEY", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment or return a default value
func getBoolEnv(key string) bool {
	const (
		positive = "TRUE"
		negative = "FALSE"
	)
	if value, exists := os.LookupEnv(key); exists {
		return value == positive
	}

	return false
}

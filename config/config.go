package config

import (
	"os"
	"strings"
)

type Config struct {
	Port        string
	Host        string
	Environment string
	SwaggerHost string
	IsDynamic   bool
}

func LoadConfig() *Config {
	config := &Config{
		Port:        getEnv("PORT", "52983"),
		Host:        getEnv("HOST", "localhost"),
		Environment: getEnv("GIN_MODE", "debug"),
		IsDynamic:   true, // Always use dynamic host detection
	}

	// For development, use localhost with port
	if config.Environment == "debug" {
		config.SwaggerHost = "localhost:" + config.Port
		config.IsDynamic = false
	} else {
		// In production, use dynamic detection (will be set by middleware)
		config.SwaggerHost = "localhost:" + config.Port // fallback
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}

package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                  string
	DatabaseURL              string
	JWTSecret                string
	LogLevel                 string
	GRPCPort                 string
	OTELEnabled              bool
	OTELServiceName          string
	OTELExporterOTLPEndpoint string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file found, relying on environment variables")
	}

	config := &Config{
		AppPort:                  os.Getenv("APP_PORT"),
		DatabaseURL:              os.Getenv("DATABASE_URL"),
		JWTSecret:                os.Getenv("JWT_SECRET"),
		LogLevel:                 os.Getenv("LOG_LEVEL"),
		GRPCPort:                 os.Getenv("GRPC_PORT"),
		OTELEnabled:              getEnvBool("OTEL_ENABLED", false),
		OTELServiceName:          os.Getenv("OTEL_SERVICE_NAME"),
		OTELExporterOTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
	}

	if config.AppPort == "" {
		config.AppPort = "8080"
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.GRPCPort == "" {
		config.GRPCPort = "9090"
	}

	if config.DatabaseURL == "" {
		return nil, ErrMissingDatabaseURL
	}

	if config.JWTSecret == "" {
		return nil, ErrMissingJWTSecret
	}

	return config, nil
}

func MustLoadConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	return config
}

func getEnvBool(key string, defaultValue bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return defaultValue
	}

	return value == "true" || value == "1" || value == "yes"
}

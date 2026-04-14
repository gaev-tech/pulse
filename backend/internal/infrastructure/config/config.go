package config

import (
	"os"
)

type Config struct {
	Port        string
	Env         string
	FrontendURL string
	Database    DatabaseConfig
	JWT         JWTConfig
	MinIO       MinIOConfig
	Email       EmailConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type JWTConfig struct {
	Secret string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type EmailConfig struct {
	APIKey    string
	FromEmail string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			Name:     getEnv("POSTGRES_DB", "pulse"),
			User:     getEnv("POSTGRES_USER", "pulse"),
			Password: getEnv("POSTGRES_PASSWORD", "pulse"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "changeme"),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "pulse"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		},
		Email: EmailConfig{
			APIKey:    getEnv("BREVO_API_KEY", ""),
			FromEmail: getEnv("EMAIL_FROM", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all required environment-based settings.
type Config struct {
	AppPort            string
	UploadSecretKey    string
	PublicBaseURL      string
	AWSRegion          string
	AWSBucket          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSSessionToken    string
	AWSEndpointURL     string
	AWSUsePathStyle    bool
}

// LoadConfig reads and validates config from .env / environment.
func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:            getEnv("APP_PORT", "3000"),
		UploadSecretKey:    strings.TrimSpace(os.Getenv("UPLOAD_SECRET_KEY")),
		PublicBaseURL:      strings.TrimSpace(os.Getenv("PUBLIC_BASE_URL")),
		AWSRegion:          strings.TrimSpace(os.Getenv("AWS_REGION")),
		AWSBucket:          strings.TrimSpace(os.Getenv("AWS_S3_BUCKET")),
		AWSAccessKeyID:     strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID")),
		AWSSecretAccessKey: strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY")),
		AWSSessionToken:    strings.TrimSpace(os.Getenv("AWS_SESSION_TOKEN")),
		AWSEndpointURL:     strings.TrimSpace(os.Getenv("AWS_ENDPOINT_URL")),
		AWSUsePathStyle:    getBoolEnv("AWS_S3_USE_PATH_STYLE", false),
	}

	missing := make([]string, 0)
	if cfg.UploadSecretKey == "" {
		missing = append(missing, "UPLOAD_SECRET_KEY")
	}
	if cfg.PublicBaseURL == "" {
		missing = append(missing, "PUBLIC_BASE_URL")
	}
	if cfg.AWSRegion == "" {
		missing = append(missing, "AWS_REGION")
	}
	if cfg.AWSBucket == "" {
		missing = append(missing, "AWS_S3_BUCKET")
	}
	if cfg.AWSAccessKeyID == "" {
		missing = append(missing, "AWS_ACCESS_KEY_ID")
	}
	if cfg.AWSSecretAccessKey == "" {
		missing = append(missing, "AWS_SECRET_ACCESS_KEY")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing env: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getBoolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

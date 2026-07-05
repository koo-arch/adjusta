package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const developmentEnv = "development"

type Config struct {
	DatabaseURL           string
	Port                  string
	SessionSecret         string
	CORSAllowOrigins      []string
	RedirectURLAfterLogin string
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleRedirectURI     string
	GoEnv                 string
	Domain                string
}

func New() Config {
	return Config{
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		Port:                  defaultString(os.Getenv("PORT"), "8080"),
		SessionSecret:         os.Getenv("SESSION_SECRET"),
		CORSAllowOrigins:      splitAndTrim(os.Getenv("CORS_ALLOW_ORIGINS")),
		RedirectURLAfterLogin: os.Getenv("REDIRECT_URL_AFTER_LOGIN"),
		GoogleClientID:        os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:    os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURI:     os.Getenv("GOOGLE_REDIRECT_URI"),
		GoEnv:                 os.Getenv("GO_ENV"),
		Domain:                os.Getenv("DOMAIN"),
	}
}

func NewServer() (Config, error) {
	cfg := New()
	if err := cfg.validateServer(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func NewDatabase() (Config, error) {
	cfg := New()
	if err := cfg.validateDatabase(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) IsDevelopment() bool {
	return c.GoEnv == developmentEnv
}

func defaultString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func splitAndTrim(value string) []string {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (c Config) validateServer() error {
	var errs []error
	required := map[string]string{
		"DATABASE_URL":         c.DatabaseURL,
		"SESSION_SECRET":       c.SessionSecret,
		"GOOGLE_CLIENT_ID":     c.GoogleClientID,
		"GOOGLE_CLIENT_SECRET": c.GoogleClientSecret,
		"GOOGLE_REDIRECT_URI":  c.GoogleRedirectURI,
	}

	for name, value := range required {
		if value == "" {
			errs = append(errs, fmt.Errorf("%s is not set", name))
		}
	}

	return errors.Join(errs...)
}

func (c Config) validateDatabase() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}
	return nil
}

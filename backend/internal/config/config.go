package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return Config{
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		Port:                  defaultString(os.Getenv("PORT"), "8080"),
		SessionSecret:         os.Getenv("SESSION_SECRET"),
		CORSAllowOrigins:      strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ","),
		RedirectURLAfterLogin: os.Getenv("REDIRECT_URL_AFTER_LOGIN"),
		GoogleClientID:        os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:    os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURI:     os.Getenv("GOOGLE_REDIRECT_URI"),
		GoEnv:                 os.Getenv("GO_ENV"),
		Domain:                os.Getenv("DOMAIN"),
	}
}

func defaultString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

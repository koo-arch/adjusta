package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/koo-arch/adjusta-backend/configs"
)

var GoogleOAuthConfig   oauth2.Config

func init() {
	configs.LoadEnv()

	GoogleOAuthConfig = oauth2.Config{
		ClientID:     configs.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: configs.GetEnv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  configs.GetEnv("GOOGLE_REDIRECT_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/calendar", "openid"},
		Endpoint:     google.Endpoint,
	}
}

func GetGoogleAuthConfig() *oauth2.Config {
	return &GoogleOAuthConfig
}

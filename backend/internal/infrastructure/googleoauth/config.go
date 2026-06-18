package googleoauth

import (
	infraConfigs "github.com/koo-arch/adjusta-backend/internal/infrastructure/configs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var config oauth2.Config

func init() {
	infraConfigs.LoadEnv()

	config = oauth2.Config{
		ClientID:     infraConfigs.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: infraConfigs.GetEnv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  infraConfigs.GetEnv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/calendar",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
}

func GetConfig() *oauth2.Config {
	return &config
}

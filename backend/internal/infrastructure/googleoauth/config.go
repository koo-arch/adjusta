package googleoauth

import (
	appconfig "github.com/koo-arch/adjusta-backend/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var config oauth2.Config

func init() {
	appConfig := appconfig.New()

	config = oauth2.Config{
		ClientID:     appConfig.GoogleClientID,
		ClientSecret: appConfig.GoogleClientSecret,
		RedirectURL:  appConfig.GoogleRedirectURI,
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

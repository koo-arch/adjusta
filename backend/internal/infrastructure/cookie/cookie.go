package cookie

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/configs"
)

const (
	SessionCookieName = "session"
	SessionTokenKey   = "session_token"
	OAuthStateKey     = "oauth_state"
)

func init() {
	configs.LoadEnv()
}

func DefaultCookieOptions() sessions.Options {
	secure := configs.GetEnv("GO_ENV") != "development"
	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}
	return sessions.Options{
		Domain:   configs.GetEnv("DOMAIN"),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	}
}

func ExpiredCookie(name string) *http.Cookie {
	domain := configs.GetEnv("DOMAIN")
	secure := configs.GetEnv("GO_ENV") != "development"

	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}
	return &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	}
}

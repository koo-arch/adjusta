package cookie

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func SetCookie(c *gin.Context, name, value string, maxAge int) {
	opt := DefaultCookieOptions()
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   opt.MaxAge,
		Path:     opt.Path,
		Domain:   opt.Domain,
		HttpOnly: opt.HttpOnly,
		Secure:   opt.Secure,
		SameSite: opt.SameSite,
	}
	http.SetCookie(c.Writer, cookie)
}

func DeleteCookie(c *gin.Context, name string) {
	domain := configs.GetEnv("DOMAIN")
	secure := configs.GetEnv("GO_ENV") != "development"

	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	}
	http.SetCookie(c.Writer, cookie)
}

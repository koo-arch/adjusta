package cookie

import (
	"net/http"

	"github.com/gin-contrib/sessions"
)

const (
	SessionCookieName = "session"
	SessionTokenKey   = "session_token"
	OAuthStateKey     = "oauth_state"

	defaultMaxAge = 60 * 60 * 24 * 7
)

type Manager struct {
	domain string
	secure bool
}

func NewManager(domain string, secure bool) *Manager {
	return &Manager{
		domain: domain,
		secure: secure,
	}
}

func (m *Manager) Options() sessions.Options {
	return sessions.Options{
		Domain:   m.domain,
		Path:     "/",
		MaxAge:   defaultMaxAge,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func (m *Manager) Expired(name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   m.domain,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	}
}

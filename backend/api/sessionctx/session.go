package sessionctx

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	infraCookie "github.com/koo-arch/adjusta-backend/internal/infrastructure/cookie"
)

func NewOAuthState(c *gin.Context) (string, error) {
	state := uuid.NewString()
	session := sessions.Default(c)
	session.Set(infraCookie.OAuthStateKey, state)
	if err := session.Save(); err != nil {
		return "", err
	}
	return state, nil
}

func OAuthState(c *gin.Context) (string, bool) {
	state, ok := sessions.Default(c).Get(infraCookie.OAuthStateKey).(string)
	return state, ok
}

func DeleteOAuthState(c *gin.Context) {
	sessions.Default(c).Delete(infraCookie.OAuthStateKey)
}

func SessionToken(c *gin.Context) (string, bool) {
	token, ok := sessions.Default(c).Get(infraCookie.SessionTokenKey).(string)
	return token, ok
}

func SetSessionToken(c *gin.Context, token string) error {
	session := sessions.Default(c)
	session.Set(infraCookie.SessionTokenKey, token)
	return session.Save()
}

func PutAuthenticatedSessionToken(c *gin.Context, token string) {
	c.Set(infraCookie.SessionTokenKey, token)
}

func Renew(c *gin.Context) error {
	session := sessions.Default(c)
	session.Options(infraCookie.DefaultCookieOptions())
	return session.Save()
}

func Clear(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		return err
	}
	http.SetCookie(c.Writer, infraCookie.ExpiredCookie(infraCookie.SessionCookieName))
	return nil
}

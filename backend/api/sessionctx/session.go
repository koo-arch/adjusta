package sessionctx

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/cookie"
)

func NewOAuthState(c *gin.Context) (string, error) {
	state := uuid.NewString()
	session := sessions.Default(c)
	session.Set(cookie.OAuthStateKey, state)
	if err := session.Save(); err != nil {
		return "", err
	}
	return state, nil
}

func OAuthState(c *gin.Context) (string, bool) {
	state, ok := sessions.Default(c).Get(cookie.OAuthStateKey).(string)
	return state, ok
}

func SessionToken(c *gin.Context) (string, bool) {
	token, ok := sessions.Default(c).Get(cookie.SessionTokenKey).(string)
	return token, ok
}

func CompleteOAuthSignIn(c *gin.Context, token string) error {
	session := sessions.Default(c)
	session.Delete(cookie.OAuthStateKey)
	session.Set(cookie.SessionTokenKey, token)
	return session.Save()
}

func PutAuthenticatedSessionToken(c *gin.Context, token string) {
	c.Set(cookie.SessionTokenKey, token)
}

func Renew(c *gin.Context, cookieManager *cookie.Manager) error {
	session := sessions.Default(c)
	session.Options(cookieManager.Options())
	return session.Save()
}

func Clear(c *gin.Context, cookieManager *cookie.Manager) error {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		return err
	}
	http.SetCookie(c.Writer, cookieManager.Expired(cookie.SessionCookieName))
	return nil
}

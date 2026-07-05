package sessionctx

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/cookie"
)

type CookieSessionStore struct {
	cookieManager *cookie.Manager
}

func NewCookieSessionStore(cookieManager *cookie.Manager) *CookieSessionStore {
	return &CookieSessionStore{cookieManager: cookieManager}
}

func (s *CookieSessionStore) IssueOAuthState(c *gin.Context) (string, error) {
	state := uuid.NewString()
	session := sessions.Default(c)
	session.Set(cookie.OAuthStateKey, state)
	if err := session.Save(); err != nil {
		return "", err
	}
	return state, nil
}

func (s *CookieSessionStore) CompleteOAuthSignIn(c *gin.Context, token string) error {
	session := sessions.Default(c)
	session.Delete(cookie.OAuthStateKey)
	session.Set(cookie.SessionTokenKey, token)
	return session.Save()
}

func (s *CookieSessionStore) ClearOAuthState(c *gin.Context) error {
	session := sessions.Default(c)
	session.Delete(cookie.OAuthStateKey)
	return session.Save()
}

func (s *CookieSessionStore) OAuthState(c *gin.Context) (string, bool) {
	state, ok := sessions.Default(c).Get(cookie.OAuthStateKey).(string)
	return state, ok
}

func (s *CookieSessionStore) SessionToken(c *gin.Context) (string, bool) {
	token, ok := sessions.Default(c).Get(cookie.SessionTokenKey).(string)
	return token, ok
}

func (s *CookieSessionStore) RenewSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Options(s.cookieManager.Options())
	return session.Save()
}

func (s *CookieSessionStore) ClearSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		return err
	}
	http.SetCookie(c.Writer, s.cookieManager.Expired(cookie.SessionCookieName))
	return nil
}

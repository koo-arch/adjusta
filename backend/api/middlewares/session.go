package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
)

type SessionMiddleware struct {
	cookieSessionStore *sessionctx.CookieSessionStore
}

func NewSessionMiddleware(cookieSessionStore *sessionctx.CookieSessionStore) *SessionMiddleware {
	return &SessionMiddleware{cookieSessionStore: cookieSessionStore}
}

func (sm *SessionMiddleware) SessionRenewal() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := sm.cookieSessionStore.RenewSession(c); err != nil {
			respond.Internal(c, "failed to save session")
			return
		}
		c.Next()
	}
}

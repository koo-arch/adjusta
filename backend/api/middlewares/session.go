package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/cookie"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
)

type SessionMiddleware struct {
	cookieManager *cookie.Manager
}

func NewSessionMiddleware(cookieManager *cookie.Manager) *SessionMiddleware {
	return &SessionMiddleware{cookieManager: cookieManager}
}

func (sm *SessionMiddleware) SessionRenewal() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := sessionctx.Renew(c, sm.cookieManager); err != nil {
			respond.Internal(c, "failed to save session")
			return
		}
		c.Next()
	}
}

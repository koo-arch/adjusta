package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
)

type SessionMiddleware struct {
	Middleware *Middleware
}

func NewSessionMiddleware(middleware *Middleware) *SessionMiddleware {
	return &SessionMiddleware{Middleware: middleware}
}

func (sm *SessionMiddleware) SessionRenewal() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := sessionctx.Renew(c); err != nil {
			respond.Internal(c, "failed to save session")
			return
		}
		c.Next()
	}
}

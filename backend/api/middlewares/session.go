package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/respond"
	infraCookie "github.com/koo-arch/adjusta-backend/internal/infrastructure/cookie"
)

type SessionMiddleware struct {
	Middleware *Middleware
}

func NewSessionMiddleware(middleware *Middleware) *SessionMiddleware {
	return &SessionMiddleware{Middleware: middleware}
}

func (sm *SessionMiddleware) SessionRenewal() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Options(infraCookie.DefaultCookieOptions())
		if err := session.Save(); err != nil {
			respond.Internal(c, "failed to save session")
			return
		}
		c.Next()
	}
}

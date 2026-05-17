package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/koo-arch/adjusta-backend/cookie"
)

type SessionMiddleware struct {
	Middleware *Middleware
}

func NewSessionMiddleware(middleware *Middleware) *SessionMiddleware {
	return &SessionMiddleware{Middleware: middleware}
}

func(sm *SessionMiddleware) SessionRenewal() gin.HandlerFunc {
	return func (c *gin.Context) {
		session := sessions.Default(c)
		session.Options(cookie.DefaultCookieOptions())
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
			c.Abort()
			return
		}
		c.Next()
	}
}
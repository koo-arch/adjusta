package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/cookie"
)

type AuthMiddleware struct {
	middleware *Middleware
}

func NewAuthMiddleware(middleware *Middleware) *AuthMiddleware {
	return &AuthMiddleware{middleware: middleware}
}

func (am *AuthMiddleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken, ok := session.Get("session_token").(string)
		if !ok || sessionToken == "" {
			am.clearSession(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報がありません"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		authenticator := am.middleware.Server.SessionAuthenticator

		authenticatedUser, err := authenticator.AuthenticateSession(ctx, sessionToken)
		if err != nil {
			log.Printf("failed to authenticate session: %v", err)
			am.clearSession(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証に失敗しました"})
			c.Abort()
			return
		}

		c.Set("user_id", authenticatedUser.ID)
		c.Set("email", authenticatedUser.Email)
		c.Set("session_token", sessionToken)
		c.Next()
	}
}

func (am *AuthMiddleware) clearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Printf("failed to clear session cookie: %v", err)
	}
	cookie.DeleteCookie(c, "session")
}

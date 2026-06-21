package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
)

type AuthMiddleware struct {
	middleware *Middleware
}

func NewAuthMiddleware(middleware *Middleware) *AuthMiddleware {
	return &AuthMiddleware{middleware: middleware}
}

func (am *AuthMiddleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, ok := sessionctx.SessionToken(c)
		if !ok || sessionToken == "" {
			am.clearSession(c)
			respond.Unauthorized(c, "認証情報がありません")
			return
		}

		ctx := c.Request.Context()
		authenticator := am.middleware.Server.SessionAuthenticator

		authenticatedUser, err := authenticator.AuthenticateSession(ctx, sessionToken)
		if err != nil {
			log.Printf("failed to authenticate session: %v", err)
			am.clearSession(c)
			respond.Unauthorized(c, "認証に失敗しました")
			return
		}

		c.Set("user_id", authenticatedUser.ID)
		c.Set("email", authenticatedUser.Email)
		sessionctx.PutAuthenticatedSessionToken(c, sessionToken)
		c.Next()
	}
}

func (am *AuthMiddleware) clearSession(c *gin.Context) {
	if err := sessionctx.Clear(c); err != nil {
		log.Printf("failed to clear session cookie: %v", err)
	}
}

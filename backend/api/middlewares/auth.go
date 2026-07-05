package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

type AuthMiddleware struct {
	sessionAuthenticator SessionAuthenticator
	cookieSessionStore   *sessionctx.CookieSessionStore
}

func NewAuthMiddleware(sessionAuthenticator SessionAuthenticator, cookieSessionStore *sessionctx.CookieSessionStore) *AuthMiddleware {
	return &AuthMiddleware{
		sessionAuthenticator: sessionAuthenticator,
		cookieSessionStore:   cookieSessionStore,
	}
}

func (am *AuthMiddleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, ok := am.cookieSessionStore.SessionToken(c)
		if !ok || sessionToken == "" {
			am.clearSession(c)
			respond.Unauthorized(c, "認証情報がありません")
			return
		}

		ctx := c.Request.Context()

		authenticatedUser, err := am.sessionAuthenticator.AuthenticateSession(ctx, sessionToken)
		if err != nil {
			log.Printf("failed to authenticate session: %v", err)
			if internalErrors.IsKind(err, internalErrors.KindUnauthorized) {
				am.clearSession(c)
			}
			respond.Error(c, err, "認証に失敗しました")
			return
		}

		requestctx.SetUser(c, authenticatedUser.ID, authenticatedUser.Email)
		c.Next()
	}
}

func (am *AuthMiddleware) clearSession(c *gin.Context) {
	if err := am.cookieSessionStore.ClearSession(c); err != nil {
		log.Printf("failed to clear session cookie: %v", err)
	}
}

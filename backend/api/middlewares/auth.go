package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/cookie"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
)

type AuthMiddleware struct {
	middleware *Middleware
}

func NewAuthMiddleware(middleware *Middleware) *AuthMiddleware {
	return &AuthMiddleware{middleware: middleware}
}

func (am *AuthMiddleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := am.middleware.Server.Client
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			log.Printf("failed to get access token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報がありません"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		jwtManager := am.middleware.Server.JWTManager

		// トークンの有効性を確認
		email, err := jwtManager.VerifyToken(ctx, client, accessToken, "access")
		if err != nil {
			// トークンの有効期限が切れている場合はリフレッシュトークンを利用してトークンを再発行
			if strings.Contains(err.Error(), "token is expired") {
				token, err := am.tokenRefresh(c)
				if err != nil {
					log.Printf("failed to refresh token: %v", err)
					c.JSON(http.StatusUnauthorized, gin.H{"error": "トークンの再発行に失敗しました"})
					c.Abort()
					return
				}

				accessToken = token.AccessToken
				email, err = jwtManager.VerifyToken(ctx, client, accessToken, "access")
				if err != nil {
					log.Printf("failed to verify token: %v", err)
					c.JSON(http.StatusUnauthorized, gin.H{"error": "トークン認証に失敗しました"})
					c.Abort()
					return
				}
			} else {
				log.Printf("failed to verify token: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "トークン認証に失敗しました"})
				c.Abort()
				return
			}
		}

		c.Set("email", email)
		c.Set("access_token", accessToken)
		c.Next()
	}
}

func (am *AuthMiddleware) tokenRefresh(c *gin.Context) (*models.JWTToken, error) {
	ctx := c.Request.Context()
	client := am.middleware.Server.Client
	session := sessions.Default(c)
	useridStr, ok := session.Get("userid").(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userid")
	}

	userid, err := uuid.Parse(useridStr)
	if err != nil {
		return nil, err
	}

	jwtManager := am.middleware.Server.JWTManager
	userRepo := am.middleware.Server.UserRepo

	// リフレッシュトークンの取得
	u, err := userRepo.Read(ctx, nil, userid, user.UserQueryOptions{})
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンの有効性を確認
	email, err := jwtManager.VerifyToken(ctx, client, u.RefreshToken, "refresh")
	if err != nil {
		return nil, err
	}

	// トークンの再発行
	token, err := jwtManager.GenerateTokens(ctx, client, email)
	if err != nil {
		return nil, err
	}

	maxAge := int(token.RefreshExpiration.Sub(time.Now()).Seconds())
	cookie.SetCookie(c, "access_token", token.AccessToken, maxAge)

	// リフレッシュトークンの更新
	_, err = userRepo.Update(ctx, nil, userid, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

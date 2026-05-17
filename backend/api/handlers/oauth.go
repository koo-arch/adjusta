package handlers

import (
	"net/http"
	"time"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/cookie"
	"github.com/koo-arch/adjusta-backend/configs"
	"github.com/koo-arch/adjusta-backend/internal/google/oauth"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"golang.org/x/oauth2"
	"github.com/koo-arch/adjusta-backend/utils"
)

type OauthHandler struct {
	handler *Handler
}

func NewOauthHandler(handler *Handler) *OauthHandler {
	return &OauthHandler{handler: handler}
}

func (oh *OauthHandler) GoogleLoginHandler(c *gin.Context) {
	url := oauth.GetGoogleAuthConfig().AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (oh *OauthHandler) LogoutHandler(c *gin.Context) {
	// セッションをクリア
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Printf("failed to save session for account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "セッションの保存に失敗しました"})
		return
	}

	// クッキーを削除
	cookie.DeleteCookie(c, "access_token")
	cookie.DeleteCookie(c, "refresh_token")
	cookie.DeleteCookie(c, "session")

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (oh *OauthHandler) GoogleCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := oh.handler.Server.Client
		session := sessions.Default(c)
		// クエリパラメータからcodeを取得
		code := c.Query("code")
		if code == "" {
			log.Printf("missing code")
			c.JSON(http.StatusBadRequest, gin.H{"error": "codeがありません"})
			return
		}

		ctx := c.Request.Context()

		// Googleからトークンを取得
		oauthToken, err := oauth.GetGoogleAuthConfig().Exchange(ctx, code)
		if err != nil {
			log.Printf("failed to exchange oauthToken: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuthトークンの取得に失敗しました"})
			return
		}

		// Googleからユーザー情報を取得
		userInfo, err := userinfo.FetchGoogleUserInfo(ctx, oauthToken)
		if err != nil {
			log.Printf("failed to fetch user info: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました"})
			return
		}

		jwtManager := oh.handler.Server.JWTManager
		authManager := oh.handler.Server.AuthManager
		// アプリ独自のJWTトークンを生成
		jwtToken, err := jwtManager.GenerateTokens(ctx, client, userInfo.Email)
		if err != nil {
			log.Printf("failed to generate jwtToken: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWTトークンの生成に失敗しました"})
			return
		}

		// ユーザー情報をデータベースに保存
		u, err := authManager.ProcessUserSignIn(ctx, userInfo, jwtToken, oauthToken)
		if err != nil {
			log.Printf("failed to create or update user: %v", err)
			utils.HandleAPIError(c, err, "ユーザーの作成または更新に失敗しました")
			return
		}

		
		// アカウントのoauthトークンを検証
		_, err = authManager.VerifyOAuthToken(ctx, u.ID)
		if err != nil {
			log.Printf("failed to reuse token source for account: %s, error: %v", u.Email, err)
			utils.HandleAPIError(c, err, "OAuthトークンの再利用に失敗しました")
			return
		}

		// クッキーにトークンをセット
		maxAge := int(jwtToken.RefreshExpiration.Sub(time.Now()).Seconds())
		cookie.SetCookie(c, "access_token", jwtToken.AccessToken, maxAge)

		// セッションにユーザー情報をセット
		session.Set("googleid", userInfo.GoogleID)
		session.Set("userid", u.ID.String())
		if err := session.Save(); err != nil {
			log.Printf("failed to save session for account: %s, error: %v", u.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "セッションの保存に失敗しました"})
			return
		}

		// リダイレクト
		c.Redirect(http.StatusTemporaryRedirect, configs.GetEnv("REDIRECT_URL_AFTER_LOGIN"))
	}
}

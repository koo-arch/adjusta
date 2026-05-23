package handlers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/configs"
	"github.com/koo-arch/adjusta-backend/cookie"
	"github.com/koo-arch/adjusta-backend/internal/google/oauth"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"github.com/koo-arch/adjusta-backend/utils"
	"golang.org/x/oauth2"
)

type OauthHandler struct {
	handler *Handler
}

func NewOauthHandler(handler *Handler) *OauthHandler {
	return &OauthHandler{handler: handler}
}

func (oh *OauthHandler) GoogleLoginHandler(c *gin.Context) {
	session := sessions.Default(c)
	state := uuid.NewString()
	session.Set("oauth_state", state)
	if err := session.Save(); err != nil {
		log.Printf("failed to save oauth state: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "認証状態の保存に失敗しました"})
		return
	}

	url := oauth.GetGoogleAuthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (oh *OauthHandler) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	sessionToken, _ := session.Get("session_token").(string)

	authManager := oh.handler.Server.AuthManager
	if err := authManager.DeleteSession(c.Request.Context(), sessionToken); err != nil {
		utils.HandleAPIError(c, err, "セッションの削除に失敗しました")
		return
	}

	session.Clear()
	if err := session.Save(); err != nil {
		log.Printf("failed to save cleared session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "セッションの保存に失敗しました"})
		return
	}

	cookie.DeleteCookie(c, "session")

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (oh *OauthHandler) GoogleCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		state := c.Query("state")
		code := c.Query("code")
		if code == "" {
			log.Printf("missing code")
			c.JSON(http.StatusBadRequest, gin.H{"error": "codeがありません"})
			return
		}

		expectedState, ok := session.Get("oauth_state").(string)
		if !ok || expectedState == "" || state != expectedState {
			log.Printf("invalid oauth state")
			c.JSON(http.StatusBadRequest, gin.H{"error": "stateが不正です"})
			return
		}

		ctx := c.Request.Context()

		oauthToken, err := oauth.GetGoogleAuthConfig().Exchange(ctx, code)
		if err != nil {
			log.Printf("failed to exchange oauth token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuthトークンの取得に失敗しました"})
			return
		}

		userInfo, err := userinfo.FetchGoogleUserInfo(ctx, oauthToken)
		if err != nil {
			log.Printf("failed to fetch user info: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました"})
			return
		}

		authManager := oh.handler.Server.AuthManager
		u, err := authManager.ProcessUserSignIn(ctx, userInfo, oauthToken)
		if err != nil {
			log.Printf("failed to create or update user: %v", err)
			utils.HandleAPIError(c, err, "ユーザーの作成または更新に失敗しました")
			return
		}

		entSession, err := authManager.CreateSession(ctx, u.ID)
		if err != nil {
			log.Printf("failed to create app session: %v", err)
			utils.HandleAPIError(c, err, "ログインセッションの作成に失敗しました")
			return
		}

		session.Delete("oauth_state")
		session.Set("session_token", entSession.SessionToken)
		if err := session.Save(); err != nil {
			log.Printf("failed to save session token for account: %s, error: %v", u.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "セッションの保存に失敗しました"})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, configs.GetEnv("REDIRECT_URL_AFTER_LOGIN"))
	}
}

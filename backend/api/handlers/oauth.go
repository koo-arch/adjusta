package handlers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/respond"
	infraConfigs "github.com/koo-arch/adjusta-backend/internal/infrastructure/configs"
	infraCookie "github.com/koo-arch/adjusta-backend/internal/infrastructure/cookie"
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
		respond.Internal(c, "認証状態の保存に失敗しました")
		return
	}

	url := oh.handler.Server.AuthSessionUsecase.GoogleLoginURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (oh *OauthHandler) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	sessionToken, _ := session.Get("session_token").(string)

	authSessionUsecase := oh.handler.Server.AuthSessionUsecase
	if err := authSessionUsecase.Logout(c.Request.Context(), sessionToken); err != nil {
		respond.Error(c, err, "セッションの削除に失敗しました")
		return
	}

	session.Clear()
	if err := session.Save(); err != nil {
		log.Printf("failed to save cleared session: %v", err)
		respond.Internal(c, "セッションの保存に失敗しました")
		return
	}

	infraCookie.DeleteCookie(c, "session")

	respond.OKMessage(c, "logged out")
}

func (oh *OauthHandler) GoogleCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		state := c.Query("state")
		code := c.Query("code")
		if code == "" {
			log.Printf("missing code")
			respond.BadRequest(c, "codeがありません")
			return
		}

		expectedState, ok := session.Get("oauth_state").(string)
		if !ok || expectedState == "" || state != expectedState {
			log.Printf("invalid oauth state")
			respond.BadRequest(c, "stateが不正です")
			return
		}

		ctx := c.Request.Context()
		authSessionUsecase := oh.handler.Server.AuthSessionUsecase
		signInResult, err := authSessionUsecase.CompleteGoogleSignIn(ctx, code)
		if err != nil {
			log.Printf("failed to complete google sign in: %v", err)
			respond.Error(c, err, "ログインに失敗しました")
			return
		}

		session.Delete("oauth_state")
		session.Set("session_token", signInResult.SessionToken)
		if err := session.Save(); err != nil {
			log.Printf("failed to save session token for account: %s, error: %v", signInResult.UserEmail, err)
			respond.Internal(c, "セッションの保存に失敗しました")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, infraConfigs.GetEnv("REDIRECT_URL_AFTER_LOGIN"))
	}
}

package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
	infraConfigs "github.com/koo-arch/adjusta-backend/internal/infrastructure/configs"
)

type OauthHandler struct {
	handler *Handler
}

func NewOauthHandler(handler *Handler) *OauthHandler {
	return &OauthHandler{handler: handler}
}

func (oh *OauthHandler) GoogleLoginHandler(c *gin.Context) {
	state, err := sessionctx.NewOAuthState(c)
	if err != nil {
		log.Printf("failed to save oauth state: %v", err)
		respond.Internal(c, "認証状態の保存に失敗しました")
		return
	}

	url := oh.handler.Server.AuthSessionUsecase.GoogleLoginURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (oh *OauthHandler) LogoutHandler(c *gin.Context) {
	sessionToken, _ := sessionctx.SessionToken(c)

	authSessionUsecase := oh.handler.Server.AuthSessionUsecase
	if err := authSessionUsecase.Logout(c.Request.Context(), sessionToken); err != nil {
		respond.Error(c, err, "セッションの削除に失敗しました")
		return
	}

	if err := sessionctx.Clear(c); err != nil {
		log.Printf("failed to save cleared session: %v", err)
		respond.Internal(c, "セッションの保存に失敗しました")
		return
	}

	respond.OKMessage(c, "logged out")
}

func (oh *OauthHandler) GoogleCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		state := c.Query("state")
		code := c.Query("code")
		if code == "" {
			log.Printf("missing code")
			respond.BadRequest(c, "codeがありません")
			return
		}

		expectedState, ok := sessionctx.OAuthState(c)
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

		sessionctx.DeleteOAuthState(c)
		if err := sessionctx.SetSessionToken(c, signInResult.SessionToken); err != nil {
			log.Printf("failed to save session token for account: %s, error: %v", signInResult.UserEmail, err)
			respond.Internal(c, "セッションの保存に失敗しました")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, infraConfigs.GetEnv("REDIRECT_URL_AFTER_LOGIN"))
	}
}

package oauth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/cookie"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
)

type Handler struct {
	oauthUsecase          OAuthUsecase
	redirectURLAfterLogin string
	cookieManager         *cookie.Manager
}

func NewHandler(oauthUsecase OAuthUsecase, redirectURLAfterLogin string, cookieManager *cookie.Manager) *Handler {
	return &Handler{
		oauthUsecase:          oauthUsecase,
		redirectURLAfterLogin: redirectURLAfterLogin,
		cookieManager:         cookieManager,
	}
}

func (oh *Handler) GoogleLoginHandler(c *gin.Context) {
	state, err := sessionctx.NewOAuthState(c)
	if err != nil {
		log.Printf("failed to save oauth state: %v", err)
		respond.Internal(c, "認証状態の保存に失敗しました")
		return
	}

	url := oh.oauthUsecase.GoogleLoginURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (oh *Handler) LogoutHandler(c *gin.Context) {
	sessionToken, _ := sessionctx.SessionToken(c)

	if err := oh.oauthUsecase.Logout(c.Request.Context(), sessionToken); err != nil {
		respond.Error(c, err, "セッションの削除に失敗しました")
		return
	}

	if err := sessionctx.Clear(c, oh.cookieManager); err != nil {
		log.Printf("failed to save cleared session: %v", err)
		respond.Internal(c, "セッションの保存に失敗しました")
		return
	}

	respond.OKMessage(c, "logged out")
}

func (oh *Handler) GoogleCallbackHandler() gin.HandlerFunc {
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
		signInResult, err := oh.oauthUsecase.CompleteGoogleSignIn(ctx, code)
		if err != nil {
			log.Printf("failed to complete google sign in: %v", err)
			respond.Error(c, err, "ログインに失敗しました")
			return
		}

		if err := sessionctx.CompleteOAuthSignIn(c, signInResult.SessionToken); err != nil {
			log.Printf("failed to save session token for account: %s, error: %v", signInResult.UserEmail, err)
			respond.Internal(c, "セッションの保存に失敗しました")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, oh.redirectURLAfterLogin)
	}
}

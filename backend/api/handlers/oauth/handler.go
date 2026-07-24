package oauth

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
)

type Handler struct {
	oauthUsecase          OAuthUsecase
	redirectURLAfterLogin string
	cookieSessionStore    *sessionctx.CookieSessionStore
}

func NewHandler(oauthUsecase OAuthUsecase, redirectURLAfterLogin string, cookieSessionStore *sessionctx.CookieSessionStore) *Handler {
	return &Handler{
		oauthUsecase:          oauthUsecase,
		redirectURLAfterLogin: redirectURLAfterLogin,
		cookieSessionStore:    cookieSessionStore,
	}
}

func (oh *Handler) GoogleLoginHandler(c *gin.Context) {
	state, err := oh.cookieSessionStore.IssueOAuthState(c)
	if err != nil {
		log.Printf("failed to save oauth state: %v", err)
		respond.Internal(c, "認証状態の保存に失敗しました")
		return
	}

	url := oh.oauthUsecase.GoogleLoginURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (oh *Handler) GoogleReauthorizationHandler(c *gin.Context) {
	returnTo := safeReturnPath(c.Query("return_to"), oh.redirectURLAfterLogin)
	state, err := oh.cookieSessionStore.IssueGoogleReauthorizationState(c, returnTo)
	if err != nil {
		log.Printf("failed to save google reauthorization state: %v", err)
		respond.Internal(c, "認証状態の保存に失敗しました")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, oh.oauthUsecase.GoogleLoginURL(state))
}

func (oh *Handler) LogoutHandler(c *gin.Context) {
	sessionToken, _ := oh.cookieSessionStore.SessionToken(c)

	logoutErr := oh.oauthUsecase.Logout(c.Request.Context(), sessionToken)

	if err := oh.cookieSessionStore.ClearSession(c); err != nil {
		log.Printf("failed to save cleared session: %v", err)
		respond.Internal(c, "セッションの保存に失敗しました")
		return
	}
	if logoutErr != nil {
		respond.Error(c, logoutErr, "セッションの削除に失敗しました")
		return
	}

	respond.OKMessage(c, "logged out")
}

func (oh *Handler) GoogleCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		state := c.Query("state")
		code := c.Query("code")
		oauthError := c.Query("error")
		if oauthError != "" {
			log.Printf("google oauth callback returned error: %s", oauthError)
			if err := oh.cookieSessionStore.ClearOAuthState(c); err != nil {
				log.Printf("failed to clear oauth state after callback error: %v", err)
				respond.Internal(c, "認証状態の保存に失敗しました")
				return
			}
			respond.BadRequest(c, "Google認証が完了しませんでした")
			return
		}

		if code == "" {
			log.Printf("missing code")
			if err := oh.cookieSessionStore.ClearOAuthState(c); err != nil {
				log.Printf("failed to clear oauth state after missing code: %v", err)
				respond.Internal(c, "認証状態の保存に失敗しました")
				return
			}
			respond.BadRequest(c, "codeがありません")
			return
		}

		expectedState, ok := oh.cookieSessionStore.OAuthState(c)
		if !ok || expectedState == "" || state != expectedState {
			log.Printf("invalid oauth state")
			if err := oh.cookieSessionStore.ClearOAuthState(c); err != nil {
				log.Printf("failed to clear invalid oauth state: %v", err)
				respond.Internal(c, "認証状態の保存に失敗しました")
				return
			}
			respond.BadRequest(c, "stateが不正です")
			return
		}

		ctx := c.Request.Context()
		if oh.cookieSessionStore.IsGoogleReauthorization(c) {
			sessionToken, ok := oh.cookieSessionStore.SessionToken(c)
			if !ok || sessionToken == "" {
				if clearErr := oh.cookieSessionStore.ClearOAuthState(c); clearErr != nil {
					log.Printf("failed to clear google reauthorization state: %v", clearErr)
				}
				respond.Unauthorized(c, "認証情報がありません")
				return
			}

			returnTo := safeReturnPath(oh.cookieSessionStore.OAuthReturnTo(c), oh.redirectURLAfterLogin)
			if err := oh.oauthUsecase.CompleteGoogleReauthorization(ctx, code, sessionToken); err != nil {
				log.Printf("failed to complete google reauthorization: %v", err)
				if clearErr := oh.cookieSessionStore.ClearOAuthState(c); clearErr != nil {
					log.Printf("failed to clear google reauthorization state after error: %v", clearErr)
				}
				respond.Error(c, err, "Google再認可に失敗しました")
				return
			}
			if err := oh.cookieSessionStore.ClearOAuthState(c); err != nil {
				log.Printf("failed to clear google reauthorization state: %v", err)
				respond.Internal(c, "認証状態の保存に失敗しました")
				return
			}

			c.Redirect(http.StatusTemporaryRedirect, returnTo)
			return
		}

		signInResult, err := oh.oauthUsecase.CompleteGoogleSignIn(ctx, code)
		if err != nil {
			log.Printf("failed to complete google sign in: %v", err)
			respond.Error(c, err, "ログインに失敗しました")
			return
		}

		if err := oh.cookieSessionStore.CompleteOAuthSignIn(c, signInResult.SessionToken); err != nil {
			log.Printf("failed to save session token for account: %s, error: %v", signInResult.UserEmail, err)
			respond.Internal(c, "セッションの保存に失敗しました")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, oh.redirectURLAfterLogin)
	}
}

func safeReturnPath(raw, fallback string) string {
	parsed, err := url.Parse(raw)
	if raw == "" || err != nil || parsed.IsAbs() || parsed.Host != "" || !strings.HasPrefix(parsed.Path, "/") || strings.HasPrefix(raw, "//") {
		return fallback
	}
	return parsed.RequestURI()
}

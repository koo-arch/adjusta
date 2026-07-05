package oauth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	apiCookie "github.com/koo-arch/adjusta-backend/api/cookie"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
)

type fakeOAuthUsecase struct {
	completeGoogleSignInCalled bool
	logoutErr                  error
}

func (f *fakeOAuthUsecase) GoogleLoginURL(state string) string {
	return "https://example.com/oauth?state=" + state
}

func (f *fakeOAuthUsecase) CompleteGoogleSignIn(ctx context.Context, code string) (*usecaseAuth.GoogleSignInResult, error) {
	f.completeGoogleSignInCalled = true
	return &usecaseAuth.GoogleSignInResult{SessionToken: "session-token", UserEmail: "user@example.com"}, nil
}

func (f *fakeOAuthUsecase) Logout(ctx context.Context, sessionToken string) error {
	return f.logoutErr
}

func TestGoogleCallbackHandlerReturnsBadRequestOnGoogleOAuthError(t *testing.T) {
	usecase := &fakeOAuthUsecase{}
	router := newOAuthTestRouter(usecase)

	response := performOAuthRequest(router, "/auth/google/callback?error=access_denied&state=oauth-state")

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if usecase.completeGoogleSignInCalled {
		t.Fatal("CompleteGoogleSignIn should not be called")
	}
}

func TestGoogleCallbackHandlerReturnsBadRequestOnInvalidState(t *testing.T) {
	usecase := &fakeOAuthUsecase{}
	router := newOAuthTestRouter(usecase)

	response := performOAuthRequest(router, "/auth/google/callback?code=oauth-code&state=invalid-state")

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if usecase.completeGoogleSignInCalled {
		t.Fatal("CompleteGoogleSignIn should not be called")
	}
}

func TestLogoutHandlerClearsCookieWhenSessionDeleteFails(t *testing.T) {
	usecase := &fakeOAuthUsecase{logoutErr: errors.New("delete failed")}
	router := newOAuthTestRouter(usecase)

	response := performOAuthRequest(router, "/auth/logout")

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if len(response.Result().Cookies()) == 0 {
		t.Fatal("expected cleared session cookie")
	}
}

func newOAuthTestRouter(usecase OAuthUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	store := cookie.NewStore([]byte("test-session-secret"))
	router.Use(sessions.Sessions(apiCookie.SessionCookieName, store))
	sessionStore := sessionctx.NewCookieSessionStore(apiCookie.NewManager("", false))
	handler := NewHandler(usecase, "/", sessionStore)
	router.GET("/auth/google/callback", handler.GoogleCallbackHandler())
	router.GET("/auth/logout", handler.LogoutHandler)
	return router
}

func performOAuthRequest(router http.Handler, path string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

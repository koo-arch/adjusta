package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	apiCookie "github.com/koo-arch/adjusta-backend/api/cookie"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
)

type fakeOAuthUsecase struct {
	completeGoogleSignInCalled bool
	completeCode               string
	completeResult             *usecaseAuth.GoogleSignInResult
	completeErr                error
	loginState                 string
	logoutCalled               bool
	logoutToken                string
	logoutErr                  error
}

func (f *fakeOAuthUsecase) GoogleLoginURL(state string) string {
	f.loginState = state
	return "https://example.com/oauth?state=" + state
}

func (f *fakeOAuthUsecase) CompleteGoogleSignIn(ctx context.Context, code string) (*usecaseAuth.GoogleSignInResult, error) {
	f.completeGoogleSignInCalled = true
	f.completeCode = code
	if f.completeResult != nil || f.completeErr != nil {
		return f.completeResult, f.completeErr
	}
	return &usecaseAuth.GoogleSignInResult{SessionToken: "session-token", UserEmail: "user@example.com"}, nil
}

func (f *fakeOAuthUsecase) Logout(ctx context.Context, sessionToken string) error {
	f.logoutCalled = true
	f.logoutToken = sessionToken
	return f.logoutErr
}

func TestGoogleLoginHandlerStoresStateAndRedirects(t *testing.T) {
	usecase := &fakeOAuthUsecase{}
	router := newOAuthTestRouter(usecase)

	response := performOAuthRequest(router, "/auth/google/login")

	if response.Code != http.StatusTemporaryRedirect {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if usecase.loginState == "" {
		t.Fatal("expected OAuth state")
	}
	location, err := url.Parse(response.Header().Get("Location"))
	if err != nil {
		t.Fatalf("failed to parse redirect URL: %v", err)
	}
	if location.Query().Get("state") != usecase.loginState {
		t.Fatalf("unexpected redirect state: %s", location.Query().Get("state"))
	}
	if len(response.Result().Cookies()) == 0 {
		t.Fatal("expected session cookie")
	}
}

func TestGoogleCallbackHandlerStoresSessionAndRedirects(t *testing.T) {
	usecase := &fakeOAuthUsecase{}
	router := newOAuthTestRouter(usecase)
	loginResponse := performOAuthRequest(router, "/auth/google/login")
	loginLocation, err := url.Parse(loginResponse.Header().Get("Location"))
	if err != nil {
		t.Fatalf("failed to parse login redirect: %v", err)
	}
	state := loginLocation.Query().Get("state")

	callbackResponse := performOAuthRequestWithCookies(
		router,
		"/auth/google/callback?code=oauth-code&state="+url.QueryEscape(state),
		loginResponse.Result().Cookies(),
	)

	if callbackResponse.Code != http.StatusTemporaryRedirect {
		t.Fatalf("unexpected status: %d", callbackResponse.Code)
	}
	if callbackResponse.Header().Get("Location") != "/" {
		t.Fatalf("unexpected redirect: %s", callbackResponse.Header().Get("Location"))
	}
	if !usecase.completeGoogleSignInCalled || usecase.completeCode != "oauth-code" {
		t.Fatalf("unexpected sign-in call: %+v", usecase)
	}
	sessionResponse := performOAuthRequestWithCookies(router, "/test/session", callbackResponse.Result().Cookies())
	var sessionBody struct {
		Token string `json:"token"`
	}
	decodeOAuthResponse(t, sessionResponse, &sessionBody)
	if sessionBody.Token != "session-token" {
		t.Fatalf("unexpected session token: %s", sessionBody.Token)
	}
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

func TestGoogleCallbackHandlerReturnsBadRequestWhenCodeIsMissing(t *testing.T) {
	usecase := &fakeOAuthUsecase{}
	router := newOAuthTestRouter(usecase)

	response := performOAuthRequest(router, "/auth/google/callback?state=oauth-state")

	assertOAuthError(t, response, http.StatusBadRequest, internalErrors.KindBadRequest)
	if usecase.completeGoogleSignInCalled {
		t.Fatal("CompleteGoogleSignIn should not be called")
	}
}

func TestGoogleCallbackHandlerReturnsUsecaseErrorCode(t *testing.T) {
	usecase := &fakeOAuthUsecase{completeErr: internalErrors.NewBadGatewayError("Google認証に失敗しました")}
	router := newOAuthTestRouter(usecase)
	loginResponse := performOAuthRequest(router, "/auth/google/login")
	loginLocation, err := url.Parse(loginResponse.Header().Get("Location"))
	if err != nil {
		t.Fatalf("failed to parse login redirect: %v", err)
	}
	state := loginLocation.Query().Get("state")

	response := performOAuthRequestWithCookies(
		router,
		"/auth/google/callback?code=oauth-code&state="+url.QueryEscape(state),
		loginResponse.Result().Cookies(),
	)

	assertOAuthError(t, response, http.StatusBadGateway, internalErrors.KindBadGateway)
}

func TestLogoutHandlerDeletesSessionAndClearsCookie(t *testing.T) {
	usecase := &fakeOAuthUsecase{}
	router := newOAuthTestRouter(usecase)
	seedResponse := performOAuthRequest(router, "/test/session?token=session-token")

	response := performOAuthRequestWithCookies(router, "/auth/logout", seedResponse.Result().Cookies())

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.logoutCalled || usecase.logoutToken != "session-token" {
		t.Fatalf("unexpected logout call: %+v", usecase)
	}
	expired := false
	for _, cookie := range response.Result().Cookies() {
		if cookie.Name == apiCookie.SessionCookieName && cookie.MaxAge < 0 {
			expired = true
			break
		}
	}
	if !expired {
		t.Fatal("expected expired session cookie")
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
	router.GET("/auth/google/login", handler.GoogleLoginHandler)
	router.GET("/auth/google/callback", handler.GoogleCallbackHandler())
	router.GET("/auth/logout", handler.LogoutHandler)
	router.GET("/test/session", func(c *gin.Context) {
		session := sessions.Default(c)
		if token := c.Query("token"); token != "" {
			session.Set(apiCookie.SessionTokenKey, token)
			if err := session.Save(); err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
		}
		token, _ := session.Get(apiCookie.SessionTokenKey).(string)
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
	return router
}

func performOAuthRequest(router http.Handler, path string) *httptest.ResponseRecorder {
	return performOAuthRequestWithCookies(router, path, nil)
}

func performOAuthRequestWithCookies(router http.Handler, path string, cookies []*http.Cookie) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func decodeOAuthResponse(t *testing.T, response *httptest.ResponseRecorder, destination any) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), destination); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func assertOAuthError(t *testing.T, response *httptest.ResponseRecorder, status int, code internalErrors.Kind) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body struct {
		Code internalErrors.Kind `json:"code"`
	}
	decodeOAuthResponse(t, response, &body)
	if body.Code != code {
		t.Fatalf("unexpected error code: %s", body.Code)
	}
}

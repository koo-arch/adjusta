package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	ginSessions "github.com/gin-contrib/sessions"
	sessionCookie "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gorillaSessions "github.com/gorilla/sessions"
	apiCookie "github.com/koo-arch/adjusta-backend/api/cookie"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
)

type fakeSessionAuthenticator struct {
	user   *repoUser.User
	err    error
	called bool
	token  string
}

func (f *fakeSessionAuthenticator) AuthenticateSession(_ context.Context, token string) (*repoUser.User, error) {
	f.called = true
	f.token = token
	return f.user, f.err
}

func TestAuthUserAuthenticatesSessionAndSetsUserContext(t *testing.T) {
	authenticator := &fakeSessionAuthenticator{
		user: &repoUser.User{ID: middlewareTestUserID, Email: middlewareTestEmail},
	}
	router := newAuthMiddlewareTestRouter(authenticator)
	seedResponse := performMiddlewareRequest(router, "/test/session?token=session-token", nil)

	response := performMiddlewareRequest(router, "/protected", seedResponse.Result().Cookies())

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !authenticator.called || authenticator.token != "session-token" {
		t.Fatalf("unexpected authentication call: %+v", authenticator)
	}
	var body struct {
		UserID uuid.UUID `json:"user_id"`
		Email  string    `json:"email"`
	}
	decodeMiddlewareResponse(t, response, &body)
	if body.UserID != middlewareTestUserID || body.Email != middlewareTestEmail {
		t.Fatalf("unexpected user context: %+v", body)
	}
}

func TestAuthUserRejectsMissingSessionAndClearsCookie(t *testing.T) {
	authenticator := &fakeSessionAuthenticator{}
	router := newAuthMiddlewareTestRouter(authenticator)

	response := performMiddlewareRequest(router, "/protected", nil)

	assertMiddlewareError(t, response, http.StatusUnauthorized, internalErrors.KindUnauthorized)
	if authenticator.called {
		t.Fatal("AuthenticateSession should not be called")
	}
	assertExpiredSessionCookie(t, response)
}

func TestAuthUserClearsExpiredSession(t *testing.T) {
	authenticator := &fakeSessionAuthenticator{
		err: internalErrors.NewUnauthorizedError("セッションの有効期限が切れています"),
	}
	router := newAuthMiddlewareTestRouter(authenticator)
	seedResponse := performMiddlewareRequest(router, "/test/session?token=expired-token", nil)

	response := performMiddlewareRequest(router, "/protected", seedResponse.Result().Cookies())

	assertMiddlewareError(t, response, http.StatusUnauthorized, internalErrors.KindUnauthorized)
	if !authenticator.called || authenticator.token != "expired-token" {
		t.Fatalf("unexpected authentication call: %+v", authenticator)
	}
	assertExpiredSessionCookie(t, response)
}

func TestSessionRenewalRefreshesCookie(t *testing.T) {
	store := sessionCookie.NewStore([]byte("test-session-secret"))
	router := newSessionMiddlewareTestRouter(store)

	response := performMiddlewareRequest(router, "/renew", nil)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	cookies := response.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected renewed session cookie")
	}
	if !cookies[0].HttpOnly || cookies[0].MaxAge <= 0 || cookies[0].SameSite != http.SameSiteLaxMode {
		t.Fatalf("unexpected cookie options: %+v", cookies[0])
	}
}

func TestSessionRenewalStopsRequestWhenSaveFails(t *testing.T) {
	baseStore := sessionCookie.NewStore([]byte("test-session-secret"))
	store := &failingSessionStore{Store: baseStore}
	router := newSessionMiddlewareTestRouter(store)

	response := performMiddlewareRequest(router, "/renew", nil)

	assertMiddlewareError(t, response, http.StatusInternalServerError, internalErrors.KindInternal)
}

type fakeCalendarSyncUsecase struct {
	output []*usecaseCalendar.ExternalCalendar
	err    error
	called bool
	userID uuid.UUID
	email  string
}

func (f *fakeCalendarSyncUsecase) SyncGoogleCalendars(_ context.Context, userID uuid.UUID, email string) ([]*usecaseCalendar.ExternalCalendar, error) {
	f.called = true
	f.userID = userID
	f.email = email
	return f.output, f.err
}

func TestSyncGoogleCalendarsStoresCalendarListAndContinues(t *testing.T) {
	usecase := &fakeCalendarSyncUsecase{
		output: []*usecaseCalendar.ExternalCalendar{{CalendarID: "primary", Summary: "メイン", Primary: true}},
	}
	router := newCalendarMiddlewareTestRouter(usecase, true)

	response := performMiddlewareRequest(router, "/calendar", nil)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.called || usecase.userID != middlewareTestUserID || usecase.email != middlewareTestEmail {
		t.Fatalf("unexpected sync call: %+v", usecase)
	}
	var body struct {
		Count int `json:"count"`
	}
	decodeMiddlewareResponse(t, response, &body)
	if body.Count != 1 {
		t.Fatalf("unexpected calendar count: %d", body.Count)
	}
}

func TestSyncGoogleCalendarsRequiresUserContext(t *testing.T) {
	usecase := &fakeCalendarSyncUsecase{}
	router := newCalendarMiddlewareTestRouter(usecase, false)

	response := performMiddlewareRequest(router, "/calendar", nil)

	assertMiddlewareError(t, response, http.StatusUnauthorized, internalErrors.KindUnauthorized)
	if usecase.called {
		t.Fatal("SyncGoogleCalendars should not be called")
	}
}

func TestSyncGoogleCalendarsReturnsReauthorizationError(t *testing.T) {
	usecase := &fakeCalendarSyncUsecase{
		err: internalErrors.NewGoogleReauthorizationRequiredError("Googleの再認可が必要です"),
	}
	router := newCalendarMiddlewareTestRouter(usecase, true)

	response := performMiddlewareRequest(router, "/calendar", nil)

	assertMiddlewareError(t, response, http.StatusConflict, internalErrors.KindGoogleReauth)
}

const middlewareTestEmail = "user@example.com"

var middlewareTestUserID = uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")

func newAuthMiddlewareTestRouter(authenticator SessionAuthenticator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	store := sessionCookie.NewStore([]byte("test-session-secret"))
	router.Use(ginSessions.Sessions(apiCookie.SessionCookieName, store))
	cookieSessionStore := sessionctx.NewCookieSessionStore(apiCookie.NewManager("", false))
	authMiddleware := NewAuthMiddleware(authenticator, cookieSessionStore)
	router.GET("/test/session", seedMiddlewareSession)
	router.GET("/protected", authMiddleware.AuthUser(), func(c *gin.Context) {
		userID, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID, "email": email})
	})
	return router
}

func newSessionMiddlewareTestRouter(store ginSessions.Store) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ginSessions.Sessions(apiCookie.SessionCookieName, store))
	cookieSessionStore := sessionctx.NewCookieSessionStore(apiCookie.NewManager("", false))
	sessionMiddleware := NewSessionMiddleware(cookieSessionStore)
	router.GET("/renew", sessionMiddleware.SessionRenewal(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return router
}

func newCalendarMiddlewareTestRouter(usecase CalendarSyncUsecase, authenticated bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if authenticated {
		router.Use(func(c *gin.Context) {
			requestctx.SetUser(c, middlewareTestUserID, middlewareTestEmail)
			c.Next()
		})
	}
	calendarMiddleware := NewCalendarMiddleware(usecase)
	router.GET("/calendar", calendarMiddleware.SyncGoogleCalendars(), func(c *gin.Context) {
		calendarList, _ := c.Get("calendarList")
		calendars, _ := calendarList.([]*usecaseCalendar.ExternalCalendar)
		c.JSON(http.StatusOK, gin.H{"count": len(calendars)})
	})
	return router
}

func seedMiddlewareSession(c *gin.Context) {
	session := ginSessions.Default(c)
	session.Set(apiCookie.SessionTokenKey, c.Query("token"))
	if err := session.Save(); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func performMiddlewareRequest(router http.Handler, path string, cookies []*http.Cookie) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func decodeMiddlewareResponse(t *testing.T, response *httptest.ResponseRecorder, destination any) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), destination); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func assertMiddlewareError(t *testing.T, response *httptest.ResponseRecorder, status int, code internalErrors.Kind) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body struct {
		Code internalErrors.Kind `json:"code"`
	}
	decodeMiddlewareResponse(t, response, &body)
	if body.Code != code {
		t.Fatalf("unexpected error code: %s", body.Code)
	}
}

func assertExpiredSessionCookie(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	for _, cookie := range response.Result().Cookies() {
		if cookie.Name == apiCookie.SessionCookieName && cookie.MaxAge < 0 {
			return
		}
	}
	t.Fatal("expected expired session cookie")
}

type failingSessionStore struct {
	ginSessions.Store
}

func (s *failingSessionStore) Get(request *http.Request, name string) (*gorillaSessions.Session, error) {
	return gorillaSessions.GetRegistry(request).Get(s, name)
}

func (s *failingSessionStore) New(_ *http.Request, name string) (*gorillaSessions.Session, error) {
	session := gorillaSessions.NewSession(s, name)
	session.IsNew = true
	return session, nil
}

func (s *failingSessionStore) Save(*http.Request, http.ResponseWriter, *gorillaSessions.Session) error {
	return errors.New("save failed")
}

package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/google"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
)

type fakeProfileUsecase struct {
	output *usecaseAccount.GoogleProfile
	err    error
	called bool
	userID uuid.UUID
}

func (f *fakeProfileUsecase) FetchGoogleProfile(_ context.Context, userID uuid.UUID) (*usecaseAccount.GoogleProfile, error) {
	f.called = true
	f.userID = userID
	return f.output, f.err
}

func TestGetCurrentUserHandlerReturnsProfile(t *testing.T) {
	usecase := &fakeProfileUsecase{
		output: &google.UserProfile{
			GoogleID: "google-user",
			Email:    userHandlerTestEmail,
			Name:     "テストユーザー",
			Picture:  "https://example.com/avatar.png",
		},
	}
	router := newUserHandlerTestRouter(usecase, true)

	response := performUserHandlerRequest(router)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.called || usecase.userID != userHandlerTestUserID {
		t.Fatalf("unexpected usecase call: %+v", usecase)
	}
	var body google.UserProfile
	decodeUserHandlerResponse(t, response, &body)
	if body.GoogleID != "google-user" || body.Email != userHandlerTestEmail || body.Name != "テストユーザー" {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestGetCurrentUserHandlerRequiresUserContext(t *testing.T) {
	usecase := &fakeProfileUsecase{}
	router := newUserHandlerTestRouter(usecase, false)

	response := performUserHandlerRequest(router)

	assertUserHandlerError(t, response, http.StatusUnauthorized, internalErrors.KindUnauthorized)
	if usecase.called {
		t.Fatal("FetchGoogleProfile should not be called")
	}
}

func TestGetCurrentUserHandlerReturnsUsecaseErrorCode(t *testing.T) {
	usecase := &fakeProfileUsecase{
		err: internalErrors.NewGoogleReauthorizationRequiredError("Googleの再認可が必要です"),
	}
	router := newUserHandlerTestRouter(usecase, true)

	response := performUserHandlerRequest(router)

	assertUserHandlerError(t, response, http.StatusConflict, internalErrors.KindGoogleReauth)
}

const userHandlerTestEmail = "user@example.com"

var userHandlerTestUserID = uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

func newUserHandlerTestRouter(usecase ProfileUsecase, authenticated bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if authenticated {
		router.Use(func(c *gin.Context) {
			requestctx.SetUser(c, userHandlerTestUserID, userHandlerTestEmail)
			c.Next()
		})
	}
	handler := NewHandler(usecase)
	router.GET("/api/users/me", handler.GetCurrentUserHandler())
	return router
}

func performUserHandlerRequest(router http.Handler) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func decodeUserHandlerResponse(t *testing.T, response *httptest.ResponseRecorder, destination any) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), destination); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func assertUserHandlerError(t *testing.T, response *httptest.ResponseRecorder, status int, code internalErrors.Kind) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body struct {
		Code internalErrors.Kind `json:"code"`
	}
	decodeUserHandlerResponse(t, response, &body)
	if body.Code != code {
		t.Fatalf("unexpected error code: %s", body.Code)
	}
}

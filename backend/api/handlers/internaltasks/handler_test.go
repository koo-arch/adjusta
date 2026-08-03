package internaltasks

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fakeProcessor struct {
	err    error
	called bool
}

func (f *fakeProcessor) ProcessMessage(context.Context, uuid.UUID) error {
	f.called = true
	return f.err
}

type fakeDispatcher struct{}

func (fakeDispatcher) DispatchPending(context.Context, int) error { return nil }

type fakeVerifier struct{ err error }

func (f fakeVerifier) VerifyAuthorization(context.Context, string) error { return f.err }

func TestProcessGoogleCalendarSyncRejectsUnauthorizedRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	processor := &fakeProcessor{}
	handler := NewHandler(processor, fakeDispatcher{}, fakeVerifier{err: errors.New("invalid token")})
	router := gin.New()
	router.POST("/internal/tasks/google-calendar-sync", handler.ProcessGoogleCalendarSync)

	request := httptest.NewRequest(http.MethodPost, "/internal/tasks/google-calendar-sync", strings.NewReader(`{"outbox_message_id":"`+uuid.NewString()+`"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if processor.called {
		t.Fatal("processor must not be called for an unauthorized request")
	}
}

func TestProcessGoogleCalendarSyncReturnsRetryableStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	processor := &fakeProcessor{err: errors.New("temporary failure")}
	handler := NewHandler(processor, fakeDispatcher{}, fakeVerifier{})
	router := gin.New()
	router.POST("/internal/tasks/google-calendar-sync", handler.ProcessGoogleCalendarSync)

	request := httptest.NewRequest(http.MethodPost, "/internal/tasks/google-calendar-sync", strings.NewReader(`{"outbox_message_id":"`+uuid.NewString()+`"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("unexpected status: %d", response.Code)
	}
}

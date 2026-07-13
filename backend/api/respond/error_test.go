package respond

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

type errorResponse struct {
	Code    string              `json:"code"`
	Error   string              `json:"error"`
	Details map[string][]string `json:"details,omitempty"`
}

func TestStatusCodeForKindGoogleReauthorization(t *testing.T) {
	t.Parallel()

	if got := statusCodeForKind(internalErrors.KindGoogleReauth); got != http.StatusConflict {
		t.Fatalf("unexpected status code: %d", got)
	}
}

func TestErrorIncludesMachineReadableCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	apiErr := internalErrors.NewAPIErrorWithDetails(
		internalErrors.KindValidation,
		"入力内容を確認してください",
		map[string][]string{"title": {"必須です"}},
	)

	Error(ctx, apiErr, "fallback")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", recorder.Code)
	}
	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Code != string(internalErrors.KindValidation) {
		t.Fatalf("unexpected error code: %s", response.Code)
	}
	if response.Error != apiErr.Message {
		t.Fatalf("unexpected error message: %s", response.Error)
	}
	if len(response.Details["title"]) != 1 {
		t.Fatalf("unexpected error details: %+v", response.Details)
	}
}

func TestErrorUsesInternalCodeForUnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	Error(ctx, errors.New("database unavailable"), "サーバーでエラーが発生しました")

	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Code != string(internalErrors.KindInternal) {
		t.Fatalf("unexpected error code: %s", response.Code)
	}
}

func TestHTTPLocalErrorsIncludeMachineReadableCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		respond    func(*gin.Context, string)
		statusCode int
		code       internalErrors.Kind
	}{
		{name: "bad request", respond: BadRequest, statusCode: http.StatusBadRequest, code: internalErrors.KindBadRequest},
		{name: "unauthorized", respond: Unauthorized, statusCode: http.StatusUnauthorized, code: internalErrors.KindUnauthorized},
		{name: "internal", respond: Internal, statusCode: http.StatusInternalServerError, code: internalErrors.KindInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			tt.respond(ctx, "message")

			if recorder.Code != tt.statusCode {
				t.Fatalf("unexpected status code: %d", recorder.Code)
			}
			var response errorResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if response.Code != string(tt.code) {
				t.Fatalf("unexpected error code: %s", response.Code)
			}
		})
	}
}

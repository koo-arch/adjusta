package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/dto"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/usecase/account/calendarsetting"
)

type fakeCalendarSettingsUsecase struct {
	listOutput   []calendarsetting.CalendarSettingOutput
	listErr      error
	updateOutput *calendarsetting.CalendarSettingOutput
	updateErr    error
	getOutput    *calendarsetting.CandidateSyncSettingOutput
	getErr       error
	setOutput    *calendarsetting.CandidateSyncSettingOutput
	setErr       error

	listCalled       bool
	updateCalled     bool
	updateUserID     uuid.UUID
	updateCalendarID uuid.UUID
	updateEmail      string
	updateRequest    calendarsetting.CalendarSettingUpdateRequest
	getCalled        bool
	setCalled        bool
	setUserID        uuid.UUID
	setEmail         string
	setEnabled       bool
}

func (f *fakeCalendarSettingsUsecase) ListCalendarSettings(context.Context, uuid.UUID, string) ([]calendarsetting.CalendarSettingOutput, error) {
	f.listCalled = true
	return f.listOutput, f.listErr
}

func (f *fakeCalendarSettingsUsecase) UpdateCalendarSetting(_ context.Context, userID, userCalendarID uuid.UUID, email string, request calendarsetting.CalendarSettingUpdateRequest) (*calendarsetting.CalendarSettingOutput, error) {
	f.updateCalled = true
	f.updateUserID = userID
	f.updateCalendarID = userCalendarID
	f.updateEmail = email
	f.updateRequest = request
	return f.updateOutput, f.updateErr
}

func (f *fakeCalendarSettingsUsecase) GetCandidateSyncSetting(context.Context, uuid.UUID) (*calendarsetting.CandidateSyncSettingOutput, error) {
	f.getCalled = true
	return f.getOutput, f.getErr
}

func (f *fakeCalendarSettingsUsecase) SetCandidateSyncSetting(_ context.Context, userID uuid.UUID, email string, enabled bool) (*calendarsetting.CandidateSyncSettingOutput, error) {
	f.setCalled = true
	f.setUserID = userID
	f.setEmail = email
	f.setEnabled = enabled
	return f.setOutput, f.setErr
}

func TestListCalendarSettingsHandlerReturnsSettings(t *testing.T) {
	calendarID := uuid.New()
	usecase := &fakeCalendarSettingsUsecase{
		listOutput: []calendarsetting.CalendarSettingOutput{
			{
				ID:               uuid.New(),
				CalendarID:       calendarID,
				GoogleCalendarID: "primary",
				Summary:          "メインカレンダー",
				Role:             value.UserCalendarRolePrimary,
				IsVisible:        true,
			},
		},
	}
	router := newCalendarSettingsTestRouter(usecase, true)

	response := performCalendarSettingsRequest(router, http.MethodGet, "/api/user-calendars", "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body []dto.CalendarSetting
	decodeJSONResponse(t, response, &body)
	if len(body) != 1 {
		t.Fatalf("unexpected settings count: %d", len(body))
	}
	if body[0].CalendarID != calendarID || body[0].Role != value.UserCalendarRolePrimary {
		t.Fatalf("unexpected setting: %+v", body[0])
	}
}

func TestListCalendarSettingsHandlerReturnsUsecaseErrorCode(t *testing.T) {
	usecase := &fakeCalendarSettingsUsecase{
		listErr: internalErrors.NewGoogleReauthorizationRequiredError("Googleの再認可が必要です"),
	}
	router := newCalendarSettingsTestRouter(usecase, true)

	response := performCalendarSettingsRequest(router, http.MethodGet, "/api/user-calendars", "")

	assertErrorResponse(t, response, http.StatusConflict, internalErrors.KindGoogleReauth)
}

func TestUpdateCalendarSettingHandlerUpdatesSetting(t *testing.T) {
	userCalendarID := uuid.New()
	role := value.UserCalendarRoleReference
	visible := false
	usecase := &fakeCalendarSettingsUsecase{
		updateOutput: &calendarsetting.CalendarSettingOutput{
			ID:        userCalendarID,
			Role:      role,
			IsVisible: visible,
		},
	}
	router := newCalendarSettingsTestRouter(usecase, true)

	response := performCalendarSettingsRequest(
		router,
		http.MethodPatch,
		"/api/user-calendars/"+userCalendarID.String(),
		`{"role":"reference","is_visible":false}`,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.updateCalled {
		t.Fatal("UpdateCalendarSetting should be called")
	}
	if usecase.updateUserID != accountHandlerTestUserID || usecase.updateCalendarID != userCalendarID || usecase.updateEmail != accountHandlerTestEmail {
		t.Fatalf("unexpected update arguments: %s, %s, %s", usecase.updateUserID, usecase.updateCalendarID, usecase.updateEmail)
	}
	if usecase.updateRequest.Role == nil || *usecase.updateRequest.Role != role {
		t.Fatalf("unexpected role: %+v", usecase.updateRequest.Role)
	}
	if usecase.updateRequest.IsVisible == nil || *usecase.updateRequest.IsVisible {
		t.Fatalf("unexpected visibility: %+v", usecase.updateRequest.IsVisible)
	}
}

func TestUpdateCalendarSettingHandlerRejectsInvalidRequest(t *testing.T) {
	tests := []struct {
		name string
		id   string
		body string
	}{
		{name: "invalid ID", id: "invalid", body: `{}`},
		{name: "malformed JSON", id: uuid.NewString(), body: `{"role":`},
		{name: "invalid role", id: uuid.NewString(), body: `{"role":"invalid"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &fakeCalendarSettingsUsecase{}
			router := newCalendarSettingsTestRouter(usecase, true)

			response := performCalendarSettingsRequest(router, http.MethodPatch, "/api/user-calendars/"+tt.id, tt.body)

			assertErrorResponse(t, response, http.StatusBadRequest, internalErrors.KindBadRequest)
			if usecase.updateCalled {
				t.Fatal("UpdateCalendarSetting should not be called")
			}
		})
	}
}

func TestUpdateCalendarSettingHandlerReturnsNotFound(t *testing.T) {
	usecase := &fakeCalendarSettingsUsecase{
		updateErr: internalErrors.NewNotFoundError("カレンダー設定が見つかりません"),
	}
	router := newCalendarSettingsTestRouter(usecase, true)

	response := performCalendarSettingsRequest(router, http.MethodPatch, "/api/user-calendars/"+uuid.NewString(), `{}`)

	assertErrorResponse(t, response, http.StatusNotFound, internalErrors.KindNotFound)
}

func TestGetCandidateSyncSettingHandlerReturnsDisabledBeforeCalendarCreation(t *testing.T) {
	usecase := &fakeCalendarSettingsUsecase{
		getOutput: &calendarsetting.CandidateSyncSettingOutput{Enabled: false},
	}
	router := newCalendarSettingsTestRouter(usecase, true)

	response := performCalendarSettingsRequest(router, http.MethodGet, "/api/calendar-settings/candidate-sync", "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body dto.CandidateSyncSetting
	decodeJSONResponse(t, response, &body)
	if body.Enabled {
		t.Fatal("expected candidate sync to be disabled")
	}
	if body.Calendar != nil {
		t.Fatal("expected calendar to be nil before creation")
	}
}

func TestGetCandidateSyncSettingHandlerReturnsCalendar(t *testing.T) {
	calendarID := uuid.New()
	userCalendarID := uuid.New()
	timezone := "Asia/Tokyo"
	usecase := &fakeCalendarSettingsUsecase{
		getOutput: &calendarsetting.CandidateSyncSettingOutput{
			Enabled: true,
			Calendar: &calendarsetting.CalendarSettingOutput{
				ID:                userCalendarID,
				CalendarID:        calendarID,
				GoogleCalendarID:  "adjusta@example.com",
				Summary:           "Adjusta 調整用",
				Timezone:          &timezone,
				Role:              value.UserCalendarRoleAdjustaCandidate,
				IsVisible:         true,
				SyncProposedDates: true,
			},
		},
	}
	router := newCalendarSettingsTestRouter(usecase, true)

	response := performCalendarSettingsRequest(router, http.MethodGet, "/api/calendar-settings/candidate-sync", "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body dto.CandidateSyncSetting
	decodeJSONResponse(t, response, &body)
	if !body.Enabled || body.Calendar == nil {
		t.Fatalf("unexpected response: %+v", body)
	}
	if body.Calendar.ID != userCalendarID || body.Calendar.CalendarID != calendarID {
		t.Fatalf("unexpected calendar IDs: %+v", body.Calendar)
	}
	if body.Calendar.Role != value.UserCalendarRoleAdjustaCandidate || !body.Calendar.SyncProposedDates {
		t.Fatalf("unexpected calendar settings: %+v", body.Calendar)
	}
}

func TestGetCandidateSyncSettingHandlerReturnsUsecaseErrorCode(t *testing.T) {
	usecase := &fakeCalendarSettingsUsecase{
		getErr: internalErrors.NewGoogleReauthorizationRequiredError("Googleの再認可が必要です"),
	}
	router := newCalendarSettingsTestRouter(usecase, true)

	response := performCalendarSettingsRequest(router, http.MethodGet, "/api/calendar-settings/candidate-sync", "")

	assertErrorResponse(t, response, http.StatusConflict, internalErrors.KindGoogleReauth)
}

func TestGetCandidateSyncSettingHandlerRequiresUserContext(t *testing.T) {
	usecase := &fakeCalendarSettingsUsecase{}
	router := newCalendarSettingsTestRouter(usecase, false)

	response := performCalendarSettingsRequest(router, http.MethodGet, "/api/calendar-settings/candidate-sync", "")

	assertErrorResponse(t, response, http.StatusUnauthorized, internalErrors.KindUnauthorized)
	if usecase.getCalled {
		t.Fatal("GetCandidateSyncSetting should not be called")
	}
}

func TestSetCandidateSyncSettingHandlerUpdatesEnabledValue(t *testing.T) {
	tests := []struct {
		name    string
		request string
		enabled bool
	}{
		{name: "enable", request: `{"enabled":true}`, enabled: true},
		{name: "disable", request: `{"enabled":false}`, enabled: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &fakeCalendarSettingsUsecase{
				setOutput: &calendarsetting.CandidateSyncSettingOutput{Enabled: tt.enabled},
			}
			router := newCalendarSettingsTestRouter(usecase, true)

			response := performCalendarSettingsRequest(router, http.MethodPut, "/api/calendar-settings/candidate-sync", tt.request)

			if response.Code != http.StatusOK {
				t.Fatalf("unexpected status: %d", response.Code)
			}
			if !usecase.setCalled || usecase.setEnabled != tt.enabled {
				t.Fatalf("unexpected enabled argument: %t", usecase.setEnabled)
			}
			if usecase.setUserID != accountHandlerTestUserID || usecase.setEmail != accountHandlerTestEmail {
				t.Fatalf("unexpected user: %s, %s", usecase.setUserID, usecase.setEmail)
			}
			var body dto.CandidateSyncSetting
			decodeJSONResponse(t, response, &body)
			if body.Enabled != tt.enabled {
				t.Fatalf("unexpected response: %+v", body)
			}
		})
	}
}

func TestSetCandidateSyncSettingHandlerRejectsInvalidRequest(t *testing.T) {
	tests := []struct {
		name    string
		request string
	}{
		{name: "missing enabled", request: `{}`},
		{name: "malformed JSON", request: `{"enabled":`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &fakeCalendarSettingsUsecase{}
			router := newCalendarSettingsTestRouter(usecase, true)

			response := performCalendarSettingsRequest(router, http.MethodPut, "/api/calendar-settings/candidate-sync", tt.request)

			assertErrorResponse(t, response, http.StatusBadRequest, internalErrors.KindBadRequest)
			if usecase.setCalled {
				t.Fatal("SetCandidateSyncSetting should not be called")
			}
		})
	}
}

func TestSetCandidateSyncSettingHandlerReturnsUnexpectedErrorAsInternal(t *testing.T) {
	usecase := &fakeCalendarSettingsUsecase{setErr: errors.New("unexpected error")}
	router := newCalendarSettingsTestRouter(usecase, true)

	response := performCalendarSettingsRequest(router, http.MethodPut, "/api/calendar-settings/candidate-sync", `{"enabled":true}`)

	assertErrorResponse(t, response, http.StatusInternalServerError, internalErrors.KindInternal)
}

const accountHandlerTestEmail = "user@example.com"

var accountHandlerTestUserID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

func newCalendarSettingsTestRouter(usecase CalendarSettingsUsecase, authenticated bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if authenticated {
		router.Use(func(c *gin.Context) {
			requestctx.SetUser(c, accountHandlerTestUserID, accountHandlerTestEmail)
			c.Next()
		})
	}
	handler := NewHandler(nil, usecase)
	router.GET("/api/user-calendars", handler.ListCalendarSettingsHandler())
	router.PATCH("/api/user-calendars/:id", handler.UpdateCalendarSettingHandler())
	router.GET("/api/calendar-settings/candidate-sync", handler.GetCandidateSyncSettingHandler())
	router.PUT("/api/calendar-settings/candidate-sync", handler.SetCandidateSyncSettingHandler())
	return router
}

func performCalendarSettingsRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func decodeJSONResponse(t *testing.T, response *httptest.ResponseRecorder, destination any) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), destination); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func assertErrorResponse(t *testing.T, response *httptest.ResponseRecorder, status int, code internalErrors.Kind) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body struct {
		Code internalErrors.Kind `json:"code"`
	}
	decodeJSONResponse(t, response, &body)
	if body.Code != code {
		t.Fatalf("unexpected error code: %s", body.Code)
	}
}

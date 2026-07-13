package events

import (
	"context"
	"encoding/json"
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
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type fakeDraftUsecase struct {
	createOutput  *usecaseEvents.EventDraftDetailOutput
	createErr     error
	createCalled  bool
	createUserID  uuid.UUID
	createEmail   string
	createRequest usecaseEvents.DraftCreationRequest

	updateErr     error
	updateCalled  bool
	updateUserID  uuid.UUID
	updateEventID uuid.UUID
	updateEmail   string
	updateRequest usecaseEvents.DraftUpdateRequest

	deleteErr     error
	deleteCalled  bool
	deleteUserID  uuid.UUID
	deleteEmail   string
	deleteEventID uuid.UUID
}

func (f *fakeDraftUsecase) FetchAllDraftedEvents(context.Context, uuid.UUID, string) ([]*usecaseEvents.EventDraftDetailOutput, error) {
	return nil, nil
}

func (f *fakeDraftUsecase) FetchDraftedEventsPage(context.Context, uuid.UUID, string, usecaseEvents.SearchDraftQuery) (*usecaseEvents.EventDraftListOutput, error) {
	return nil, nil
}

func (f *fakeDraftUsecase) SearchDraftedEvents(context.Context, uuid.UUID, string, usecaseEvents.SearchDraftQuery) ([]*usecaseEvents.EventDraftDetailOutput, error) {
	return nil, nil
}

func (f *fakeDraftUsecase) SearchDraftedEventsPage(context.Context, uuid.UUID, string, usecaseEvents.SearchDraftQuery) (*usecaseEvents.EventDraftListOutput, error) {
	return nil, nil
}

func (f *fakeDraftUsecase) CreateDraftedEvents(_ context.Context, userID uuid.UUID, email string, request usecaseEvents.DraftCreationRequest) (*usecaseEvents.EventDraftDetailOutput, error) {
	f.createCalled = true
	f.createUserID = userID
	f.createEmail = email
	f.createRequest = request
	return f.createOutput, f.createErr
}

func (f *fakeDraftUsecase) UpdateDraftedEvents(_ context.Context, userID, eventID uuid.UUID, email string, request usecaseEvents.DraftUpdateRequest) error {
	f.updateCalled = true
	f.updateUserID = userID
	f.updateEventID = eventID
	f.updateEmail = email
	f.updateRequest = request
	return f.updateErr
}

func (f *fakeDraftUsecase) DeleteDraftedEvents(_ context.Context, userID uuid.UUID, email string, eventID uuid.UUID) error {
	f.deleteCalled = true
	f.deleteUserID = userID
	f.deleteEmail = email
	f.deleteEventID = eventID
	return f.deleteErr
}

type fakeDetailUsecase struct {
	output  *usecaseEvents.EventDraftDetailOutput
	err     error
	called  bool
	userID  uuid.UUID
	email   string
	eventID uuid.UUID
}

func (f *fakeDetailUsecase) FetchDraftedEventDetail(_ context.Context, userID uuid.UUID, email string, eventID uuid.UUID) (*usecaseEvents.EventDraftDetailOutput, error) {
	f.called = true
	f.userID = userID
	f.email = email
	f.eventID = eventID
	return f.output, f.err
}

func TestCreateEventDraftHandlerCreatesDraft(t *testing.T) {
	eventID := uuid.New()
	usecase := &fakeDraftUsecase{
		createOutput: &usecaseEvents.EventDraftDetailOutput{
			ID:         eventID,
			Title:      "定例ミーティング",
			Status:     value.StatusDraft,
			SyncStatus: value.SyncStatusNotSynced,
		},
	}
	router := newDraftHandlerTestRouter(usecase, nil, true)

	response := performDraftHandlerRequest(router, http.MethodPost, "/api/calendar/event/draft", validCreateDraftBody)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.createCalled {
		t.Fatal("CreateDraftedEvents should be called")
	}
	if usecase.createUserID != eventHandlerTestUserID || usecase.createEmail != eventHandlerTestEmail {
		t.Fatalf("unexpected user: %s, %s", usecase.createUserID, usecase.createEmail)
	}
	if usecase.createRequest.Title != "定例ミーティング" || len(usecase.createRequest.SelectedDates) != 1 {
		t.Fatalf("unexpected creation request: %+v", usecase.createRequest)
	}
	var body dto.EventDraftDetail
	decodeEventHandlerResponse(t, response, &body)
	if body.ID != eventID || body.Status != value.StatusDraft {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestCreateEventDraftHandlerRejectsInvalidRequest(t *testing.T) {
	tests := []struct {
		name string
		body string
		code internalErrors.Kind
	}{
		{name: "malformed JSON", body: `{"title":`, code: internalErrors.KindBadRequest},
		{name: "null JSON", body: `null`, code: internalErrors.KindBadRequest},
		{name: "missing title and dates", body: `{}`, code: internalErrors.KindBadRequest},
		{name: "invalid date range", body: `{"title":"予定","selected_dates":[{"start":"2026-07-13T11:00:00Z","end":"2026-07-13T10:00:00Z","priority":1}]}`, code: internalErrors.KindValidation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &fakeDraftUsecase{}
			router := newDraftHandlerTestRouter(usecase, nil, true)

			response := performDraftHandlerRequest(router, http.MethodPost, "/api/calendar/event/draft", tt.body)

			assertEventHandlerError(t, response, http.StatusBadRequest, tt.code)
			if usecase.createCalled {
				t.Fatal("CreateDraftedEvents should not be called")
			}
		})
	}
}

func TestCreateEventDraftHandlerRequiresUserContext(t *testing.T) {
	usecase := &fakeDraftUsecase{}
	router := newDraftHandlerTestRouter(usecase, nil, false)

	response := performDraftHandlerRequest(router, http.MethodPost, "/api/calendar/event/draft", validCreateDraftBody)

	assertEventHandlerError(t, response, http.StatusUnauthorized, internalErrors.KindUnauthorized)
	if usecase.createCalled {
		t.Fatal("CreateDraftedEvents should not be called")
	}
}

func TestFetchEventDraftDetailHandlerReturnsDraft(t *testing.T) {
	eventID := uuid.New()
	usecase := &fakeDetailUsecase{
		output: &usecaseEvents.EventDraftDetailOutput{
			ID:         eventID,
			Title:      "定例ミーティング",
			Status:     value.StatusDraft,
			SyncStatus: value.SyncStatusSynced,
		},
	}
	router := newDraftHandlerTestRouter(nil, usecase, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/calendar/event/draft/"+eventID.String(), "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.called || usecase.userID != eventHandlerTestUserID || usecase.eventID != eventID || usecase.email != eventHandlerTestEmail {
		t.Fatalf("unexpected usecase arguments: %+v", usecase)
	}
	var body dto.EventDraftDetail
	decodeEventHandlerResponse(t, response, &body)
	if body.ID != eventID || body.SyncStatus != value.SyncStatusSynced {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestFetchEventDraftDetailHandlerRejectsInvalidID(t *testing.T) {
	usecase := &fakeDetailUsecase{}
	router := newDraftHandlerTestRouter(nil, usecase, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/calendar/event/draft/invalid", "")

	assertEventHandlerError(t, response, http.StatusBadRequest, internalErrors.KindBadRequest)
	if usecase.called {
		t.Fatal("FetchDraftedEventDetail should not be called")
	}
}

func TestFetchEventDraftDetailHandlerReturnsNotFound(t *testing.T) {
	usecase := &fakeDetailUsecase{err: internalErrors.NewNotFoundError("イベントが見つかりません")}
	router := newDraftHandlerTestRouter(nil, usecase, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/calendar/event/draft/"+uuid.NewString(), "")

	assertEventHandlerError(t, response, http.StatusNotFound, internalErrors.KindNotFound)
}

func TestUpdateEventDraftHandlerUpdatesDraft(t *testing.T) {
	eventID := uuid.New()
	usecase := &fakeDraftUsecase{}
	router := newDraftHandlerTestRouter(usecase, nil, true)

	response := performDraftHandlerRequest(router, http.MethodPut, "/api/calendar/event/draft/"+eventID.String(), validUpdateDraftBody)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.updateCalled || usecase.updateUserID != eventHandlerTestUserID || usecase.updateEventID != eventID || usecase.updateEmail != eventHandlerTestEmail {
		t.Fatalf("unexpected update arguments: %+v", usecase)
	}
	if usecase.updateRequest.Title != "更新後の予定" || len(usecase.updateRequest.ProposedDates) != 1 {
		t.Fatalf("unexpected update request: %+v", usecase.updateRequest)
	}
}

func TestUpdateEventDraftHandlerRejectsInvalidRequest(t *testing.T) {
	tests := []struct {
		name string
		id   string
		body string
		code internalErrors.Kind
	}{
		{name: "invalid ID", id: "invalid", body: validUpdateDraftBody, code: internalErrors.KindBadRequest},
		{name: "malformed JSON", id: uuid.NewString(), body: `{"title":`, code: internalErrors.KindBadRequest},
		{name: "null JSON", id: uuid.NewString(), body: `null`, code: internalErrors.KindValidation},
		{name: "missing dates", id: uuid.NewString(), body: `{"title":"予定"}`, code: internalErrors.KindValidation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &fakeDraftUsecase{}
			router := newDraftHandlerTestRouter(usecase, nil, true)

			response := performDraftHandlerRequest(router, http.MethodPut, "/api/calendar/event/draft/"+tt.id, tt.body)

			assertEventHandlerError(t, response, http.StatusBadRequest, tt.code)
			if usecase.updateCalled {
				t.Fatal("UpdateDraftedEvents should not be called")
			}
		})
	}
}

func TestDeleteEventDraftHandlerDeletesDraft(t *testing.T) {
	eventID := uuid.New()
	usecase := &fakeDraftUsecase{}
	router := newDraftHandlerTestRouter(usecase, nil, true)

	response := performDraftHandlerRequest(router, http.MethodDelete, "/api/calendar/event/draft/"+eventID.String(), "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.deleteCalled || usecase.deleteUserID != eventHandlerTestUserID || usecase.deleteEmail != eventHandlerTestEmail || usecase.deleteEventID != eventID {
		t.Fatalf("unexpected delete arguments: %+v", usecase)
	}
}

func TestDeleteEventDraftHandlerRejectsInvalidID(t *testing.T) {
	usecase := &fakeDraftUsecase{}
	router := newDraftHandlerTestRouter(usecase, nil, true)

	response := performDraftHandlerRequest(router, http.MethodDelete, "/api/calendar/event/draft/invalid", "")

	assertEventHandlerError(t, response, http.StatusBadRequest, internalErrors.KindBadRequest)
	if usecase.deleteCalled {
		t.Fatal("DeleteDraftedEvents should not be called")
	}
}

func TestDeleteEventDraftHandlerReturnsNotFound(t *testing.T) {
	usecase := &fakeDraftUsecase{deleteErr: internalErrors.NewNotFoundError("イベントが見つかりません")}
	router := newDraftHandlerTestRouter(usecase, nil, true)

	response := performDraftHandlerRequest(router, http.MethodDelete, "/api/calendar/event/draft/"+uuid.NewString(), "")

	assertEventHandlerError(t, response, http.StatusNotFound, internalErrors.KindNotFound)
}

const (
	eventHandlerTestEmail = "user@example.com"
	validCreateDraftBody  = `{"title":"定例ミーティング","location":"会議室A","description":"週次確認","selected_dates":[{"start":"2026-07-13T10:00:00Z","end":"2026-07-13T11:00:00Z","priority":1}]}`
	validUpdateDraftBody  = `{"title":"更新後の予定","status":"draft","proposed_dates":[{"start":"2026-07-14T10:00:00Z","end":"2026-07-14T11:00:00Z","priority":1}]}`
)

var eventHandlerTestUserID = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

func newDraftHandlerTestRouter(draftUsecase DraftUsecase, detailUsecase DetailUsecase, authenticated bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if authenticated {
		router.Use(func(c *gin.Context) {
			requestctx.SetUser(c, eventHandlerTestUserID, eventHandlerTestEmail)
			c.Next()
		})
	}
	handler := NewHandler(nil, draftUsecase, detailUsecase, nil, nil)
	router.POST("/api/calendar/event/draft", handler.CreateEventDraftHandler())
	router.GET("/api/calendar/event/draft/:id", handler.FetchEventDraftDetailHandler())
	router.PUT("/api/calendar/event/draft/:id", handler.UpdateEventDraftHandler())
	router.DELETE("/api/calendar/event/draft/:id", handler.DeleteEventDraftHandler())
	return router
}

func performDraftHandlerRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func decodeEventHandlerResponse(t *testing.T, response *httptest.ResponseRecorder, destination any) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), destination); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func assertEventHandlerError(t *testing.T, response *httptest.ResponseRecorder, status int, code internalErrors.Kind) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body struct {
		Code internalErrors.Kind `json:"code"`
	}
	decodeEventHandlerResponse(t, response, &body)
	if body.Code != code {
		t.Fatalf("unexpected error code: %s", body.Code)
	}
}

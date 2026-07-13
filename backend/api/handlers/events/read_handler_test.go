package events

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/dto"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type fakeDraftQueryUsecase struct {
	listOutput *usecaseEvents.EventDraftListOutput
	listErr    error
	listCalled bool
	listQuery  usecaseEvents.SearchDraftQuery

	searchOutput *usecaseEvents.EventDraftListOutput
	searchErr    error
	searchCalled bool
	searchQuery  usecaseEvents.SearchDraftQuery
}

func (f *fakeDraftQueryUsecase) FetchAllDraftedEvents(context.Context, uuid.UUID, string) ([]*usecaseEvents.EventDraftDetailOutput, error) {
	return nil, nil
}

func (f *fakeDraftQueryUsecase) FetchDraftedEventsPage(_ context.Context, _ uuid.UUID, _ string, query usecaseEvents.SearchDraftQuery) (*usecaseEvents.EventDraftListOutput, error) {
	f.listCalled = true
	f.listQuery = query
	return f.listOutput, f.listErr
}

func (f *fakeDraftQueryUsecase) SearchDraftedEvents(context.Context, uuid.UUID, string, usecaseEvents.SearchDraftQuery) ([]*usecaseEvents.EventDraftDetailOutput, error) {
	return nil, nil
}

func (f *fakeDraftQueryUsecase) SearchDraftedEventsPage(_ context.Context, _ uuid.UUID, _ string, query usecaseEvents.SearchDraftQuery) (*usecaseEvents.EventDraftListOutput, error) {
	f.searchCalled = true
	f.searchQuery = query
	return f.searchOutput, f.searchErr
}

func (f *fakeDraftQueryUsecase) CreateDraftedEvents(context.Context, uuid.UUID, string, usecaseEvents.DraftCreationRequest) (*usecaseEvents.EventDraftDetailOutput, error) {
	return nil, nil
}

func (f *fakeDraftQueryUsecase) UpdateDraftedEvents(context.Context, uuid.UUID, uuid.UUID, string, usecaseEvents.DraftUpdateRequest) error {
	return nil
}

func (f *fakeDraftQueryUsecase) DeleteDraftedEvents(context.Context, uuid.UUID, string, uuid.UUID) error {
	return nil
}

type fakeAgendaUsecase struct {
	upcomingOutput []usecaseEvents.UpcomingEventOutput
	upcomingErr    error
	upcomingCalled bool
	upcomingDays   int

	needsActionOutput []usecaseEvents.NeedsActionDraftOutput
	needsActionErr    error
	needsActionCalled bool
	needsActionDays   int
}

func (f *fakeAgendaUsecase) FetchUpcomingEvents(_ context.Context, _ uuid.UUID, _ string, daysBefore int) ([]usecaseEvents.UpcomingEventOutput, error) {
	f.upcomingCalled = true
	f.upcomingDays = daysBefore
	return f.upcomingOutput, f.upcomingErr
}

func (f *fakeAgendaUsecase) FetchNeedsActionDrafts(_ context.Context, _ uuid.UUID, _ string, daysBefore int) ([]usecaseEvents.NeedsActionDraftOutput, error) {
	f.needsActionCalled = true
	f.needsActionDays = daysBefore
	return f.needsActionOutput, f.needsActionErr
}

type fakeGoogleCalendarUsecase struct {
	output []*usecaseEvents.FetchedGoogleEvent
	err    error
	called bool
}

func (f *fakeGoogleCalendarUsecase) FetchAllGoogleEvents(context.Context, uuid.UUID, string) ([]*usecaseEvents.FetchedGoogleEvent, error) {
	f.called = true
	return f.output, f.err
}

func TestFetchAllEventDraftListHandlerReturnsPaginatedDrafts(t *testing.T) {
	eventID := uuid.New()
	usecase := &fakeDraftQueryUsecase{
		listOutput: &usecaseEvents.EventDraftListOutput{
			Items:      []*usecaseEvents.EventDraftDetailOutput{{ID: eventID, Title: "予定", Status: value.StatusDraft}},
			Pagination: usecaseEvents.PaginationOutput{Page: 2, PerPage: 10, TotalItems: 11, TotalPages: 2},
		},
	}
	router := newReadHandlerTestRouter(usecase, nil, nil, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/calendar/event/draft/list?page=2&per_page=10&sort_by=title&sort_order=asc", "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.listCalled || usecase.listQuery.Page != 2 || usecase.listQuery.PerPage != 10 || usecase.listQuery.SortBy != "title" || usecase.listQuery.SortOrder != "asc" {
		t.Fatalf("unexpected list query: %+v", usecase.listQuery)
	}
	var body dto.EventDraftList
	decodeEventHandlerResponse(t, response, &body)
	if len(body.Items) != 1 || body.Items[0].ID != eventID || body.Pagination.TotalPages != 2 {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestFetchAllEventDraftListHandlerRejectsInvalidQuery(t *testing.T) {
	usecase := &fakeDraftQueryUsecase{}
	router := newReadHandlerTestRouter(usecase, nil, nil, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/calendar/event/draft/list?page=0", "")

	assertEventHandlerError(t, response, http.StatusBadRequest, internalErrors.KindBadRequest)
	if usecase.listCalled {
		t.Fatal("FetchDraftedEventsPage should not be called")
	}
}

func TestSearchEventDraftHandlerPassesFilters(t *testing.T) {
	usecase := &fakeDraftQueryUsecase{
		searchOutput: &usecaseEvents.EventDraftListOutput{
			Items:      []*usecaseEvents.EventDraftDetailOutput{},
			Pagination: usecaseEvents.PaginationOutput{Page: 1, PerPage: 20},
		},
	}
	router := newReadHandlerTestRouter(usecase, nil, nil, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/event/draft/search?title=定例&status=pending&start_time_gte=2026-07-13T10:00:00Z", "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.searchCalled || usecase.searchQuery.Title == nil || *usecase.searchQuery.Title != "定例" {
		t.Fatalf("unexpected search query: %+v", usecase.searchQuery)
	}
	if usecase.searchQuery.Status == nil || *usecase.searchQuery.Status != value.StatusActive || usecase.searchQuery.StartTimeGTE == nil {
		t.Fatalf("unexpected search filters: %+v", usecase.searchQuery)
	}
	if usecase.searchQuery.Page != 1 || usecase.searchQuery.PerPage != 20 || usecase.searchQuery.SortBy != "created_at" || usecase.searchQuery.SortOrder != "desc" {
		t.Fatalf("unexpected search defaults: %+v", usecase.searchQuery)
	}
}

func TestSearchEventDraftHandlerRejectsInvalidStatus(t *testing.T) {
	usecase := &fakeDraftQueryUsecase{}
	router := newReadHandlerTestRouter(usecase, nil, nil, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/event/draft/search?status=invalid", "")

	assertEventHandlerError(t, response, http.StatusBadRequest, internalErrors.KindBadRequest)
	if usecase.searchCalled {
		t.Fatal("SearchDraftedEventsPage should not be called")
	}
}

func TestFetchUpcomingEventsHandlerReturnsEvents(t *testing.T) {
	eventID := uuid.New()
	confirmedDateID := uuid.New()
	start := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	usecase := &fakeAgendaUsecase{
		upcomingOutput: []usecaseEvents.UpcomingEventOutput{{
			ID: eventID, Title: "確定した予定", Status: value.StatusConfirmed, ConfirmedDateID: confirmedDateID, Start: start, End: start.Add(time.Hour),
		}},
	}
	router := newReadHandlerTestRouter(nil, usecase, nil, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/event/confirmed/upcoming", "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.upcomingCalled || usecase.upcomingDays != 3 {
		t.Fatalf("unexpected upcoming request: %+v", usecase)
	}
	var body []dto.UpcomingEvent
	decodeEventHandlerResponse(t, response, &body)
	if len(body) != 1 || body[0].ID != eventID || body[0].ConfirmedDateID != confirmedDateID {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestFetchNeedsActionDraftsHandlerReturnsDrafts(t *testing.T) {
	eventID := uuid.New()
	usecase := &fakeAgendaUsecase{
		needsActionOutput: []usecaseEvents.NeedsActionDraftOutput{{ID: eventID, Title: "要対応", NeedsAttention: true}},
	}
	router := newReadHandlerTestRouter(nil, usecase, nil, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/event/draft/needs-action", "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.needsActionCalled || usecase.needsActionDays != 3 {
		t.Fatalf("unexpected needs-action request: %+v", usecase)
	}
	var body []dto.NeedsActionDraft
	decodeEventHandlerResponse(t, response, &body)
	if len(body) != 1 || body[0].ID != eventID || !body[0].NeedsAttention {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestFetchUpcomingEventsHandlerRequiresUserContext(t *testing.T) {
	usecase := &fakeAgendaUsecase{}
	router := newReadHandlerTestRouter(nil, usecase, nil, false)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/event/confirmed/upcoming", "")

	assertEventHandlerError(t, response, http.StatusUnauthorized, internalErrors.KindUnauthorized)
	if usecase.upcomingCalled {
		t.Fatal("FetchUpcomingEvents should not be called")
	}
}

func TestFetchEventListHandlerReturnsGoogleEvents(t *testing.T) {
	usecase := &fakeGoogleCalendarUsecase{
		output: []*usecaseEvents.FetchedGoogleEvent{{ID: "google-event", Summary: "Google予定", Start: "2026-07-13T10:00:00Z", End: "2026-07-13T11:00:00Z"}},
	}
	router := newReadHandlerTestRouter(nil, nil, usecase, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/calendar/list", "")

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body struct {
		Events []*dto.GoogleEvent `json:"events"`
	}
	decodeEventHandlerResponse(t, response, &body)
	if len(body.Events) != 1 || body.Events[0].ID != "google-event" {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestFetchEventListHandlerReturnsPartialResults(t *testing.T) {
	usecase := &fakeGoogleCalendarUsecase{
		output: []*usecaseEvents.FetchedGoogleEvent{{ID: "available-event"}},
		err: internalErrors.NewPartialContentError("一部のカレンダー取得に失敗しました", map[string][]string{
			"failed_calendars": {"calendar@example.com"},
		}),
	}
	router := newReadHandlerTestRouter(nil, nil, usecase, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/calendar/list", "")

	if response.Code != http.StatusPartialContent {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body struct {
		Events  []*dto.GoogleEvent  `json:"events"`
		Warning map[string][]string `json:"warning"`
	}
	decodeEventHandlerResponse(t, response, &body)
	if len(body.Events) != 1 || len(body.Warning["failed_calendars"]) != 1 {
		t.Fatalf("unexpected partial response: %+v", body)
	}
}

func TestFetchEventListHandlerReturnsReauthorizationError(t *testing.T) {
	usecase := &fakeGoogleCalendarUsecase{err: internalErrors.NewGoogleReauthorizationRequiredError("Googleの再認可が必要です")}
	router := newReadHandlerTestRouter(nil, nil, usecase, true)

	response := performDraftHandlerRequest(router, http.MethodGet, "/api/calendar/list", "")

	assertEventHandlerError(t, response, http.StatusConflict, internalErrors.KindGoogleReauth)
}

func newReadHandlerTestRouter(draftUsecase DraftUsecase, agendaUsecase AgendaUsecase, googleCalendarUsecase GoogleCalendarUsecase, authenticated bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if authenticated {
		router.Use(func(c *gin.Context) {
			requestctx.SetUser(c, eventHandlerTestUserID, eventHandlerTestEmail)
			c.Next()
		})
	}
	handler := NewHandler(googleCalendarUsecase, draftUsecase, nil, agendaUsecase, nil)
	router.GET("/api/calendar/list", handler.FetchEventListHandler())
	router.GET("/api/calendar/event/draft/list", handler.FetchAllEventDraftListHandler())
	router.GET("/api/event/draft/search", handler.SearchEventDraftHandler())
	router.GET("/api/event/confirmed/upcoming", handler.FetchUpcomingEventsHandler())
	router.GET("/api/event/draft/needs-action", handler.FetchNeedsActionDraftsHandler())
	return router
}

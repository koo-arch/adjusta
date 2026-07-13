package events

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type fakeConfirmationUsecase struct {
	err          error
	called       bool
	userID       uuid.UUID
	eventID      uuid.UUID
	email        string
	confirmation usecaseEvents.ConfirmationRequest
}

func (f *fakeConfirmationUsecase) FinalizeProposedDate(_ context.Context, userID, eventID uuid.UUID, email string, confirmation usecaseEvents.ConfirmationRequest) error {
	f.called = true
	f.userID = userID
	f.eventID = eventID
	f.email = email
	f.confirmation = confirmation
	return f.err
}

func TestEventFinalizeHandlerFinalizesProposedDate(t *testing.T) {
	eventID := uuid.New()
	proposedDateID := uuid.New()
	usecase := &fakeConfirmationUsecase{}
	router := newConfirmationHandlerTestRouter(usecase, true)
	body := `{"confirm_date":{"id":"` + proposedDateID.String() + `","google_event_id":"google-event","start":"2026-07-13T10:00:00Z","end":"2026-07-13T11:00:00Z","priority":1}}`

	response := performDraftHandlerRequest(router, http.MethodPatch, "/api/calendar/event/confirm/"+eventID.String(), body)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !usecase.called || usecase.userID != eventHandlerTestUserID || usecase.eventID != eventID || usecase.email != eventHandlerTestEmail {
		t.Fatalf("unexpected usecase arguments: %+v", usecase)
	}
	if usecase.confirmation.ID == nil || *usecase.confirmation.ID != proposedDateID {
		t.Fatalf("unexpected proposed date ID: %+v", usecase.confirmation.ID)
	}
	if usecase.confirmation.Start == nil || usecase.confirmation.End == nil || usecase.confirmation.GoogleEventID != "google-event" {
		t.Fatalf("unexpected confirmation: %+v", usecase.confirmation)
	}
}

func TestEventFinalizeHandlerRejectsInvalidRequest(t *testing.T) {
	tests := []struct {
		name string
		id   string
		body string
		code internalErrors.Kind
	}{
		{name: "invalid ID", id: "invalid", body: validConfirmationBody, code: internalErrors.KindBadRequest},
		{name: "malformed JSON", id: uuid.NewString(), body: `{"confirm_date":`, code: internalErrors.KindBadRequest},
		{name: "null JSON", id: uuid.NewString(), body: `null`, code: internalErrors.KindValidation},
		{name: "missing dates", id: uuid.NewString(), body: `{"confirm_date":{}}`, code: internalErrors.KindValidation},
		{name: "invalid date range", id: uuid.NewString(), body: `{"confirm_date":{"start":"2026-07-13T11:00:00Z","end":"2026-07-13T10:00:00Z"}}`, code: internalErrors.KindValidation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &fakeConfirmationUsecase{}
			router := newConfirmationHandlerTestRouter(usecase, true)

			response := performDraftHandlerRequest(router, http.MethodPatch, "/api/calendar/event/confirm/"+tt.id, tt.body)

			assertEventHandlerError(t, response, http.StatusBadRequest, tt.code)
			if usecase.called {
				t.Fatal("FinalizeProposedDate should not be called")
			}
		})
	}
}

func TestEventFinalizeHandlerRequiresUserContext(t *testing.T) {
	usecase := &fakeConfirmationUsecase{}
	router := newConfirmationHandlerTestRouter(usecase, false)

	response := performDraftHandlerRequest(router, http.MethodPatch, "/api/calendar/event/confirm/"+uuid.NewString(), validConfirmationBody)

	assertEventHandlerError(t, response, http.StatusUnauthorized, internalErrors.KindUnauthorized)
	if usecase.called {
		t.Fatal("FinalizeProposedDate should not be called")
	}
}

func TestEventFinalizeHandlerReturnsForbidden(t *testing.T) {
	usecase := &fakeConfirmationUsecase{err: internalErrors.NewForbiddenError("イベントを確定できません")}
	router := newConfirmationHandlerTestRouter(usecase, true)

	response := performDraftHandlerRequest(router, http.MethodPatch, "/api/calendar/event/confirm/"+uuid.NewString(), validConfirmationBody)

	assertEventHandlerError(t, response, http.StatusForbidden, internalErrors.KindForbidden)
}

const validConfirmationBody = `{"confirm_date":{"start":"2026-07-13T10:00:00Z","end":"2026-07-13T11:00:00Z","priority":1}}`

func newConfirmationHandlerTestRouter(usecase ConfirmationUsecase, authenticated bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if authenticated {
		router.Use(func(c *gin.Context) {
			requestctx.SetUser(c, eventHandlerTestUserID, eventHandlerTestEmail)
			c.Next()
		})
	}
	handler := NewHandler(nil, nil, nil, nil, usecase)
	router.PATCH("/api/calendar/event/confirm/:id", handler.EventFinalizeHandler())
	return router
}

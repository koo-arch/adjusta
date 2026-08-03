package events

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	domainProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func TestFinalizeProposedDateStoresPendingConfirmationWithoutCallingGoogle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()
	confirmedDateID := uuid.New()
	start := time.Now().UTC().Add(3 * time.Hour)
	end := start.Add(time.Hour)
	var eventMutation EventMutation

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findEventByIDFn: func(context.Context, uuid.UUID, uuid.UUID, bool) (*domainEvent.Event, error) {
					return &domainEvent.Event{ID: eventID, UserID: userID, Status: value.StatusActive}, nil
				},
				listProposedDatesByEventFn: func(context.Context, uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
					return nil, nil
				},
				createProposedDateFn: func(context.Context, ProposedDateMutation, uuid.UUID) (*domainProposedDate.ProposedDate, error) {
					return &domainProposedDate.ProposedDate{ID: confirmedDateID, EventID: eventID}, nil
				},
				updateEventFn: func(_ context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					eventMutation = opt
					return &domainEvent.Event{ID: id}, nil
				},
			},
		},
		&fakeGoogleCalendarGateway{
			upsertEventFn: func(context.Context, uuid.UUID, string, *string, string, string, string, time.Time, time.Time) (string, error) {
				t.Fatal("Google Calendar must not be called by the confirmation request")
				return "", nil
			},
		},
	)

	err := uc.FinalizeProposedDate(ctx, userID, eventID, "user@example.com", ConfirmationRequest{
		Start: &start,
		End:   &end,
	})
	if err != nil {
		t.Fatalf("FinalizeProposedDate returned error: %v", err)
	}
	if eventMutation.Status == nil || *eventMutation.Status != value.StatusConfirmed {
		t.Fatalf("unexpected event status: %#v", eventMutation.Status)
	}
	if eventMutation.SyncStatus == nil || *eventMutation.SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected sync status: %#v", eventMutation.SyncStatus)
	}
	if eventMutation.ConfirmedDateID == nil || *eventMutation.ConfirmedDateID != confirmedDateID {
		t.Fatalf("unexpected confirmed date id: %#v", eventMutation.ConfirmedDateID)
	}
}

package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func TestFinalizeProposedDateMarksSyncFailedOnGoogleError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	eventID := uuid.New()
	start := time.Now().UTC().Add(3 * time.Hour)
	end := start.Add(time.Hour)

	var failureMutation EventMutation

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
					if gotUserID != userID || gotEventID != eventID {
						t.Fatalf("unexpected find event args: %s %s", gotUserID, gotEventID)
					}
					return &domainEvent.Event{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Title:             "Finalize",
						Status:            value.StatusActive,
					}, nil
				},
				readCalendarFn: func(ctx context.Context, gotCalendarID uuid.UUID) (*EventCalendar, error) {
					if gotCalendarID != calendarID {
						t.Fatalf("unexpected calendar id: %s", gotCalendarID)
					}
					return &EventCalendar{
						ID:               calendarID,
						GoogleCalendarID: "primary-calendar",
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					failureMutation = opt
					return &domainEvent.Event{ID: eventID}, nil
				},
			},
		},
		&fakeGoogleCalendarGateway{
			fetchEventsFn: func(ctx context.Context, userID uuid.UUID, calendars []*EventCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
				t.Fatalf("FetchEvents should not be called")
				return nil, nil
			},
			upsertEventFn: func(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
				return "", errors.New("google unavailable")
			},
		},
	)

	err := uc.FinalizeProposedDate(ctx, userID, eventID, "user@example.com", ConfirmationRequest{
		Start: &start,
		End:   &end,
	})
	if !internalErrors.IsKind(err, internalErrors.KindInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
	if failureMutation.SyncStatus == nil || *failureMutation.SyncStatus != value.SyncStatusFailed {
		t.Fatalf("unexpected failure sync status: %#v", failureMutation.SyncStatus)
	}
	if failureMutation.LastSyncError == nil || *failureMutation.LastSyncError == "" {
		t.Fatalf("expected last sync error to be recorded, got %#v", failureMutation.LastSyncError)
	}
}

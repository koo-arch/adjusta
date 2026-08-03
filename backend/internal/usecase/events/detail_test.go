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

func TestFetchDraftedEventDetailDoesNotCallGoogleCalendar(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	eventID := uuid.New()
	dateID := uuid.New()
	start := time.Now().UTC().Add(time.Hour)
	storedEvent := &domainEvent.Event{
		ID:         eventID,
		UserID:     userID,
		Title:      "Draft",
		Status:     value.StatusActive,
		SyncStatus: value.SyncStatusPending,
		ProposedDates: []*domainProposedDate.ProposedDate{{
			ID:         dateID,
			EventID:    eventID,
			StartTime:  start,
			EndTime:    start.Add(time.Hour),
			Priority:   1024,
			SyncStatus: value.SyncStatusPending,
		}},
	}

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{store: &fakeEventTxStore{
			t: t,
			findEventByIDFn: func(_ context.Context, gotUserID, gotEventID uuid.UUID, withDates bool) (*domainEvent.Event, error) {
				if gotUserID != userID || gotEventID != eventID || !withDates {
					t.Fatalf("unexpected event lookup: %s %s %t", gotUserID, gotEventID, withDates)
				}
				return storedEvent, nil
			},
		}},
		&fakeGoogleCalendarGateway{
			upsertEventFn: func(context.Context, uuid.UUID, string, *string, string, string, string, time.Time, time.Time) (string, error) {
				t.Fatal("detail retrieval must not call Google Calendar")
				return "", nil
			},
		},
	)

	detail, err := uc.FetchDraftedEventDetail(context.Background(), userID, "user@example.com", eventID)
	if err != nil {
		t.Fatalf("FetchDraftedEventDetail returned error: %v", err)
	}
	if detail.SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected sync status: %s", detail.SyncStatus)
	}
	if len(detail.ProposedDates) != 1 || detail.ProposedDates[0].SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected proposed dates: %#v", detail.ProposedDates)
	}
}

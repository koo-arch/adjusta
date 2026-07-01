package events

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func TestFetchAllGoogleEventsReturnsPartialContent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendars := []*EventCalendar{
		{ID: uuid.New(), GoogleCalendarID: "cal-1", Summary: "Primary"},
	}
	events := []*FetchedGoogleEvent{
		{ID: "event-1", Summary: "Google event"},
	}

	uc := NewUsecase(
		fakeReposFromReader(&fakeEventReader{
			findPrimaryCalendarFn: func(ctx context.Context, userID uuid.UUID) (*EventCalendar, error) {
				t.Fatalf("find primary calendar should not be called")
				return nil, nil
			},
			listCalendarsByUserFn: func(ctx context.Context, gotUserID uuid.UUID) ([]*EventCalendar, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				return calendars, nil
			},
			searchEventsFn: func(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*domainEvent.Event, error) {
				t.Fatalf("search events should not be called")
				return nil, nil
			},
			findEventByIDFn: func(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
				t.Fatalf("find by id should not be called")
				return nil, nil
			},
		}),
		nil,
		&fakeGoogleCalendarGateway{
			fetchEventsFn: func(ctx context.Context, gotUserID uuid.UUID, gotCalendars []*EventCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				if len(gotCalendars) != 1 || gotCalendars[0].GoogleCalendarID != "cal-1" {
					t.Fatalf("unexpected calendars: %#v", gotCalendars)
				}
				return &GoogleEventFetchResult{
					Events:          events,
					FailedCalendars: []string{"Primary"},
				}, nil
			},
			upsertEventFn: func(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
				t.Fatalf("upsert should not be called")
				return "", nil
			},
		},
	)

	gotEvents, err := uc.FetchAllGoogleEvents(ctx, userID, "user@example.com")
	if len(gotEvents) != 1 || gotEvents[0].ID != "event-1" {
		t.Fatalf("unexpected events: %#v", gotEvents)
	}
	apiErr, ok := err.(*internalErrors.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Kind != internalErrors.KindPartial {
		t.Fatalf("unexpected error kind: %s", apiErr.Kind)
	}
	if apiErr.Details["failed_calendars"][0] != "Primary" {
		t.Fatalf("unexpected error details: %#v", apiErr.Details)
	}
}

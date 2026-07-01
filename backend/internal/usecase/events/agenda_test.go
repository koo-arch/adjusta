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

func TestFetchUpcomingEventsSortsConfirmedDates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	now := time.Now()
	confirmed := value.StatusConfirmed

	eventID1 := uuid.New()
	dateID1 := uuid.New()
	eventID2 := uuid.New()
	dateID2 := uuid.New()

	var receivedOptions EventSearchOptions

	uc := NewUsecase(
		fakeReposFromReader(&fakeEventReader{
			findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				return &EventCalendar{ID: calendarID}, nil
			},
			listCalendarsByUserFn: func(ctx context.Context, userID uuid.UUID) ([]*EventCalendar, error) {
				t.Fatalf("list calendars should not be called")
				return nil, nil
			},
			searchEventsFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, opt EventSearchOptions) ([]*domainEvent.Event, error) {
				receivedOptions = opt
				return []*domainEvent.Event{
					{
						ID:              eventID2,
						Title:           "Later",
						Status:          value.StatusConfirmed,
						ConfirmedDateID: dateID2,
						ProposedDates: []*domainProposedDate.ProposedDate{
							{ID: dateID2, StartTime: now.Add(3 * time.Hour), EndTime: now.Add(4 * time.Hour)},
						},
					},
					{
						ID:              eventID1,
						Title:           "Sooner",
						Status:          value.StatusConfirmed,
						ConfirmedDateID: dateID1,
						ProposedDates: []*domainProposedDate.ProposedDate{
							{ID: dateID1, StartTime: now.Add(1 * time.Hour), EndTime: now.Add(2 * time.Hour)},
						},
					},
				}, nil
			},
			findEventByIDFn: func(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
				t.Fatalf("find by id should not be called")
				return nil, nil
			},
		}),
		nil,
		nil,
	)

	upcoming, err := uc.FetchUpcomingEvents(ctx, userID, "user@example.com", 3)
	if err != nil {
		t.Fatalf("FetchUpcomingEvents returned error: %v", err)
	}
	if receivedOptions.Status == nil || *receivedOptions.Status != confirmed {
		t.Fatalf("expected confirmed status filter, got %#v", receivedOptions.Status)
	}
	if !receivedOptions.WithProposedDates {
		t.Fatalf("expected proposed dates to be loaded")
	}
	if len(upcoming) != 2 {
		t.Fatalf("unexpected upcoming events length: %d", len(upcoming))
	}
	if upcoming[0].ID != eventID1 || upcoming[1].ID != eventID2 {
		t.Fatalf("events are not sorted by start time: %#v", upcoming)
	}
}

func TestFetchNeedsActionDraftsFiltersActiveEvents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	now := time.Now()
	active := value.StatusActive

	var receivedOptions EventSearchOptions

	uc := NewUsecase(
		fakeReposFromReader(&fakeEventReader{
			findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				return &EventCalendar{ID: calendarID}, nil
			},
			listCalendarsByUserFn: func(ctx context.Context, userID uuid.UUID) ([]*EventCalendar, error) {
				t.Fatalf("list calendars should not be called")
				return nil, nil
			},
			searchEventsFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, opt EventSearchOptions) ([]*domainEvent.Event, error) {
				receivedOptions = opt
				return []*domainEvent.Event{
					{
						ID:     uuid.New(),
						Title:  "Needs action",
						Status: value.StatusActive,
						ProposedDates: []*domainProposedDate.ProposedDate{
							{ID: uuid.New(), StartTime: now.Add(-1 * time.Hour), EndTime: now.Add(1 * time.Hour)},
						},
					},
				}, nil
			},
			findEventByIDFn: func(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
				t.Fatalf("find by id should not be called")
				return nil, nil
			},
		}),
		nil,
		nil,
	)

	drafts, err := uc.FetchNeedsActionDrafts(ctx, userID, "user@example.com", 3)
	if err != nil {
		t.Fatalf("FetchNeedsActionDrafts returned error: %v", err)
	}
	if receivedOptions.Status == nil || *receivedOptions.Status != active {
		t.Fatalf("expected active status filter, got %#v", receivedOptions.Status)
	}
	if receivedOptions.SortBy != "ProposedDatePriority" {
		t.Fatalf("expected sort by proposed date priority, got %s", receivedOptions.SortBy)
	}
	if receivedOptions.SortOrder != "desc" {
		t.Fatalf("expected descending priority sort, got %s", receivedOptions.SortOrder)
	}
	if len(drafts) != 1 {
		t.Fatalf("unexpected needs action drafts length: %d", len(drafts))
	}
	if drafts[0].Status != value.StatusActive {
		t.Fatalf("unexpected draft status: %s", drafts[0].Status)
	}
}

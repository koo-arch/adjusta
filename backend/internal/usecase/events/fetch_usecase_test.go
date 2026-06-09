package events

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repositorymodel"
)

type fakeEventReader struct {
	findPrimaryCalendarFn func(ctx context.Context, userID uuid.UUID) (*repositorymodel.StoredCalendar, error)
	listCalendarsByUserFn func(ctx context.Context, userID uuid.UUID) ([]*repositorymodel.StoredCalendar, error)
	searchEventsFn        func(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*repositorymodel.StoredEvent, error)
	findEventBySlugFn     func(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error)
}

func (f *fakeEventReader) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
	return f.findPrimaryCalendarFn(ctx, userID)
}

func (f *fakeEventReader) ListCalendarsByUser(ctx context.Context, userID uuid.UUID) ([]*repositorymodel.StoredCalendar, error) {
	return f.listCalendarsByUserFn(ctx, userID)
}

func (f *fakeEventReader) SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*repositorymodel.StoredEvent, error) {
	return f.searchEventsFn(ctx, userID, calendarID, opt)
}

func (f *fakeEventReader) FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
	return f.findEventBySlugFn(ctx, userID, slug, withProposedDates)
}

type fakeGoogleCalendarGateway struct {
	fetchEventsFn func(ctx context.Context, userID uuid.UUID, calendars []*repositorymodel.StoredCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error)
	upsertEventFn func(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
}

func (f *fakeGoogleCalendarGateway) FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*repositorymodel.StoredCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
	return f.fetchEventsFn(ctx, userID, calendars, startTime, endTime)
}

func (f *fakeGoogleCalendarGateway) UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
	return f.upsertEventFn(ctx, userID, calendarID, existingGoogleEventID, title, location, description, start, end)
}

func TestFetchAllGoogleEventsReturnsPartialContent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendars := []*repositorymodel.StoredCalendar{
		{ID: uuid.New(), GoogleCalendarID: "cal-1", Summary: "Primary"},
	}
	events := []*appmodel.GoogleEvent{
		{ID: "event-1", Summary: "Google event"},
	}

	uc := NewUsecase(
		&fakeEventReader{
			findPrimaryCalendarFn: func(ctx context.Context, userID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
				t.Fatalf("find primary calendar should not be called")
				return nil, nil
			},
			listCalendarsByUserFn: func(ctx context.Context, gotUserID uuid.UUID) ([]*repositorymodel.StoredCalendar, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				return calendars, nil
			},
			searchEventsFn: func(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*repositorymodel.StoredEvent, error) {
				t.Fatalf("search events should not be called")
				return nil, nil
			},
			findEventBySlugFn: func(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
				t.Fatalf("find by slug should not be called")
				return nil, nil
			},
		},
		nil,
		&fakeGoogleCalendarGateway{
			fetchEventsFn: func(ctx context.Context, gotUserID uuid.UUID, gotCalendars []*repositorymodel.StoredCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
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

func TestFetchUpcomingEventsSortsConfirmedDates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	now := time.Now()
	confirmed := domainvalue.StatusConfirmed

	eventID1 := uuid.New()
	dateID1 := uuid.New()
	eventID2 := uuid.New()
	dateID2 := uuid.New()

	var receivedOptions EventSearchOptions

	uc := NewUsecase(
		&fakeEventReader{
			findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				return &repositorymodel.StoredCalendar{ID: calendarID}, nil
			},
			listCalendarsByUserFn: func(ctx context.Context, userID uuid.UUID) ([]*repositorymodel.StoredCalendar, error) {
				t.Fatalf("list calendars should not be called")
				return nil, nil
			},
			searchEventsFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, opt EventSearchOptions) ([]*repositorymodel.StoredEvent, error) {
				receivedOptions = opt
				return []*repositorymodel.StoredEvent{
					{
						ID:              eventID2,
						Title:           "Later",
						Status:          domainvalue.StatusConfirmed,
						ConfirmedDateID: dateID2,
						ProposedDates: []*repositorymodel.StoredProposedDate{
							{ID: dateID2, StartTime: now.Add(3 * time.Hour), EndTime: now.Add(4 * time.Hour)},
						},
					},
					{
						ID:              eventID1,
						Title:           "Sooner",
						Status:          domainvalue.StatusConfirmed,
						ConfirmedDateID: dateID1,
						ProposedDates: []*repositorymodel.StoredProposedDate{
							{ID: dateID1, StartTime: now.Add(1 * time.Hour), EndTime: now.Add(2 * time.Hour)},
						},
					},
				}, nil
			},
			findEventBySlugFn: func(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
				t.Fatalf("find by slug should not be called")
				return nil, nil
			},
		},
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
	active := domainvalue.StatusActive

	var receivedOptions EventSearchOptions

	uc := NewUsecase(
		&fakeEventReader{
			findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				return &repositorymodel.StoredCalendar{ID: calendarID}, nil
			},
			listCalendarsByUserFn: func(ctx context.Context, userID uuid.UUID) ([]*repositorymodel.StoredCalendar, error) {
				t.Fatalf("list calendars should not be called")
				return nil, nil
			},
			searchEventsFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, opt EventSearchOptions) ([]*repositorymodel.StoredEvent, error) {
				receivedOptions = opt
				return []*repositorymodel.StoredEvent{
					{
						ID:     uuid.New(),
						Title:  "Needs action",
						Status: domainvalue.StatusActive,
						Slug:   "needs-action",
						ProposedDates: []*repositorymodel.StoredProposedDate{
							{ID: uuid.New(), StartTime: now.Add(-1 * time.Hour), EndTime: now.Add(1 * time.Hour)},
						},
					},
				}, nil
			},
			findEventBySlugFn: func(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
				t.Fatalf("find by slug should not be called")
				return nil, nil
			},
		},
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
	if len(drafts) != 1 {
		t.Fatalf("unexpected needs action drafts length: %d", len(drafts))
	}
	if drafts[0].Status != domainvalue.StatusActive {
		t.Fatalf("unexpected draft status: %s", drafts[0].Status)
	}
}

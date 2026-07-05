package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	domainProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func TestFetchDraftedEventDetailResyncsProposedDates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()
	dateID1 := uuid.New()
	dateID2 := uuid.New()
	start1 := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	end1 := start1.Add(time.Hour)
	start2 := start1.Add(24 * time.Hour)
	end2 := start2.Add(time.Hour)
	existingGoogleEventID := "google-existing-2"
	upsertedEventIDs := []string{"google-created-1", "google-updated-2"}
	upsertCallCount := 0

	storedEvent := &domainEvent.Event{
		ID:          eventID,
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "Discuss roadmap",
		Status:      value.StatusActive,
		SyncStatus:  value.SyncStatusPending,
		ProposedDates: []*domainProposedDate.ProposedDate{
			{
				ID:         dateID1,
				StartTime:  start1,
				EndTime:    end1,
				Priority:   20,
				Status:     value.ProposedDateStatusActive,
				SyncStatus: value.SyncStatusPending,
			},
			{
				ID:            dateID2,
				GoogleEventID: &existingGoogleEventID,
				StartTime:     start2,
				EndTime:       end2,
				Priority:      10,
				Status:        value.ProposedDateStatusActive,
				SyncStatus:    value.SyncStatusPending,
			},
		},
	}

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					if !withProposedDates {
						t.Fatal("expected proposed dates to be loaded")
					}
					return storedEvent, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{
						ID:                uuid.New(),
						GoogleCalendarID:  "adjusta-candidate",
						SyncProposedDates: true,
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
					for _, proposedDate := range storedEvent.ProposedDates {
						if proposedDate.ID != id {
							continue
						}
						applyProposedDateMutation(proposedDate, opt)
						return proposedDate, nil
					}
					t.Fatalf("unexpected proposed date id: %s", id)
					return nil, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					applyEventMutation(storedEvent, opt)
					return storedEvent, nil
				},
			},
		},
		&fakeGoogleCalendarGateway{
			upsertEventFn: func(ctx context.Context, gotUserID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				if calendarID != "adjusta-candidate" {
					t.Fatalf("unexpected calendar id: %s", calendarID)
				}
				if title != "Draft" || location != "Tokyo" || description != "Discuss roadmap" {
					t.Fatalf("unexpected event payload: %s %s %s", title, location, description)
				}
				if upsertCallCount == 0 && existingGoogleEventID != nil {
					t.Fatalf("expected first upsert to create new event, got %#v", existingGoogleEventID)
				}
				if upsertCallCount == 1 && (existingGoogleEventID == nil || *existingGoogleEventID != "google-existing-2") {
					t.Fatalf("unexpected existing google event id: %#v", existingGoogleEventID)
				}
				id := upsertedEventIDs[upsertCallCount]
				upsertCallCount++
				return id, nil
			},
			fetchEventsFn: func(ctx context.Context, userID uuid.UUID, calendars []*EventCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
				t.Fatalf("fetch events should not be called")
				return nil, nil
			},
		},
	)

	detail, err := uc.FetchDraftedEventDetail(ctx, userID, "user@example.com", eventID)
	if err != nil {
		t.Fatalf("FetchDraftedEventDetail returned error: %v", err)
	}
	if upsertCallCount != 2 {
		t.Fatalf("expected two upsert calls, got %d", upsertCallCount)
	}
	if detail.SyncStatus != value.SyncStatusSynced {
		t.Fatalf("unexpected event sync status: %s", detail.SyncStatus)
	}
	if detail.LastSyncedAt == nil {
		t.Fatal("expected event last synced at to be set")
	}
	if detail.LastSyncError != nil {
		t.Fatalf("expected event last sync error to be cleared, got %#v", detail.LastSyncError)
	}

	datesByID := make(map[uuid.UUID]ProposedDateOutput, len(detail.ProposedDates))
	for _, proposedDate := range detail.ProposedDates {
		if proposedDate.ID == nil {
			t.Fatalf("expected proposed date id, got %#v", proposedDate)
		}
		datesByID[*proposedDate.ID] = proposedDate
	}

	for id, expectedGoogleEventID := range map[uuid.UUID]string{
		dateID1: upsertedEventIDs[0],
		dateID2: upsertedEventIDs[1],
	} {
		proposedDate, ok := datesByID[id]
		if !ok {
			t.Fatalf("missing proposed date: %s", id)
		}
		if proposedDate.GoogleEventID == nil || *proposedDate.GoogleEventID != expectedGoogleEventID {
			t.Fatalf("unexpected google event id for %s: %#v", id, proposedDate.GoogleEventID)
		}
		if proposedDate.SyncStatus != value.SyncStatusSynced {
			t.Fatalf("unexpected sync status for %s: %s", id, proposedDate.SyncStatus)
		}
		if proposedDate.LastSyncedAt == nil {
			t.Fatalf("expected last synced at for %s", id)
		}
		if proposedDate.LastSyncError != nil {
			t.Fatalf("expected last sync error to be cleared for %s, got %#v", id, proposedDate.LastSyncError)
		}
	}
}

func TestFetchDraftedEventDetailMarksSyncFailureButReturnsDetail(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()
	dateID := uuid.New()
	start := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	upsertCallCount := 0

	storedEvent := &domainEvent.Event{
		ID:          eventID,
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "Discuss roadmap",
		Status:      value.StatusActive,
		SyncStatus:  value.SyncStatusPending,
		ProposedDates: []*domainProposedDate.ProposedDate{
			{
				ID:         dateID,
				StartTime:  start,
				EndTime:    end,
				Priority:   10,
				Status:     value.ProposedDateStatusActive,
				SyncStatus: value.SyncStatusPending,
			},
		},
	}

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return storedEvent, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					return &EventCalendar{
						ID:                uuid.New(),
						GoogleCalendarID:  "adjusta-candidate",
						SyncProposedDates: true,
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
					if id != dateID {
						t.Fatalf("unexpected proposed date id: %s", id)
					}
					applyProposedDateMutation(storedEvent.ProposedDates[0], opt)
					return storedEvent.ProposedDates[0], nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					applyEventMutation(storedEvent, opt)
					return storedEvent, nil
				},
			},
		},
		&fakeGoogleCalendarGateway{
			upsertEventFn: func(ctx context.Context, gotUserID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
				upsertCallCount++
				return "", errors.New("google unavailable")
			},
			fetchEventsFn: func(ctx context.Context, userID uuid.UUID, calendars []*EventCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
				t.Fatalf("fetch events should not be called")
				return nil, nil
			},
		},
	)

	detail, err := uc.FetchDraftedEventDetail(ctx, userID, "user@example.com", eventID)
	if err != nil {
		t.Fatalf("FetchDraftedEventDetail returned error: %v", err)
	}
	if upsertCallCount != 1 {
		t.Fatalf("expected one upsert call, got %d", upsertCallCount)
	}
	if detail.SyncStatus != value.SyncStatusFailed {
		t.Fatalf("unexpected event sync status: %s", detail.SyncStatus)
	}
	if detail.LastSyncError == nil || *detail.LastSyncError != "google unavailable" {
		t.Fatalf("unexpected event last sync error: %#v", detail.LastSyncError)
	}
	if len(detail.ProposedDates) != 1 {
		t.Fatalf("unexpected proposed dates: %#v", detail.ProposedDates)
	}
	if detail.ProposedDates[0].SyncStatus != value.SyncStatusFailed {
		t.Fatalf("unexpected proposed date sync status: %s", detail.ProposedDates[0].SyncStatus)
	}
	if detail.ProposedDates[0].LastSyncError == nil || *detail.ProposedDates[0].LastSyncError != "google unavailable" {
		t.Fatalf("unexpected proposed date last sync error: %#v", detail.ProposedDates[0].LastSyncError)
	}
}

func TestFetchDraftedEventDetailSkipsResyncWhenCandidateSyncDisabled(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()
	dateID := uuid.New()
	start := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	googleEventID := "google-event-id"

	storedEvent := &domainEvent.Event{
		ID:          eventID,
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "Discuss roadmap",
		Status:      value.StatusActive,
		SyncStatus:  value.SyncStatusNotSynced,
		ProposedDates: []*domainProposedDate.ProposedDate{
			{
				ID:            dateID,
				GoogleEventID: &googleEventID,
				StartTime:     start,
				EndTime:       end,
				Priority:      10,
				Status:        value.ProposedDateStatusActive,
				SyncStatus:    value.SyncStatusNotSynced,
			},
		},
	}

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return storedEvent, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					return &EventCalendar{
						ID:                uuid.New(),
						GoogleCalendarID:  "adjusta-candidate",
						SyncProposedDates: false,
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
					t.Fatal("UpdateProposedDate should not be called when candidate sync is disabled")
					return nil, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					t.Fatal("UpdateEvent should not be called when candidate sync is disabled")
					return nil, nil
				},
			},
		},
		&fakeGoogleCalendarGateway{
			upsertEventFn: func(ctx context.Context, gotUserID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
				t.Fatal("UpsertEvent should not be called when candidate sync is disabled")
				return "", nil
			},
			fetchEventsFn: func(ctx context.Context, userID uuid.UUID, calendars []*EventCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
				t.Fatalf("fetch events should not be called")
				return nil, nil
			},
		},
	)

	detail, err := uc.FetchDraftedEventDetail(ctx, userID, "user@example.com", eventID)
	if err != nil {
		t.Fatalf("FetchDraftedEventDetail returned error: %v", err)
	}
	if detail.SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected event sync status: %s", detail.SyncStatus)
	}
	if len(detail.ProposedDates) != 1 || detail.ProposedDates[0].SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected proposed dates: %#v", detail.ProposedDates)
	}
}

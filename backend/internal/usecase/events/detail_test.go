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

func TestFetchDraftedEventDetailSkipsSyncedProposedDatesUntilEventChanges(t *testing.T) {
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
	upsertCallCount := 0
	expectedTitle := "Draft"

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
				if title != expectedTitle || location != "Tokyo" || description != "Discuss roadmap" {
					t.Fatalf("unexpected event payload: %s %s %s", title, location, description)
				}
				var expectedExistingID *string
				var returnedID string
				switch upsertCallCount {
				case 0:
					returnedID = "google-created-1"
				case 1:
					expectedExistingID = stringPointer("google-existing-2")
					returnedID = "google-updated-2"
				case 2:
					expectedExistingID = stringPointer("google-created-1")
					returnedID = "google-created-1"
				case 3:
					expectedExistingID = stringPointer("google-updated-2")
					returnedID = "google-updated-2"
				default:
					t.Fatalf("unexpected upsert call: %d", upsertCallCount+1)
				}

				if expectedExistingID == nil && existingGoogleEventID != nil {
					t.Fatalf("expected a new event, got existing id %#v", existingGoogleEventID)
				}
				if expectedExistingID != nil && (existingGoogleEventID == nil || *existingGoogleEventID != *expectedExistingID) {
					t.Fatalf("expected existing google event id %q, got %#v", *expectedExistingID, existingGoogleEventID)
				}

				upsertCallCount++
				return returnedID, nil
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
		t.Fatalf("expected two initial upsert calls, got %d", upsertCallCount)
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
		dateID1: "google-created-1",
		dateID2: "google-updated-2",
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

	// 同期済みの詳細を再取得しても、Google Calendar API は呼ばない。
	secondDetail, err := uc.FetchDraftedEventDetail(ctx, userID, "user@example.com", eventID)
	if err != nil {
		t.Fatalf("second FetchDraftedEventDetail returned error: %v", err)
	}
	if upsertCallCount != 2 {
		t.Fatalf("expected no additional upsert calls, got %d total calls", upsertCallCount)
	}
	for _, proposedDate := range secondDetail.ProposedDates {
		if proposedDate.GoogleEventID == nil {
			t.Fatalf("expected google event id after repeated sync: %#v", proposedDate)
		}
	}

	// Adjusta 側でイベント基本情報が変更された場合は、同期済み候補も再同期する。
	expectedTitle = "Updated Draft"
	storedEvent.Title = expectedTitle
	storedEvent.SyncStatus = value.SyncStatusPending

	updatedDetail, err := uc.FetchDraftedEventDetail(ctx, userID, "user@example.com", eventID)
	if err != nil {
		t.Fatalf("FetchDraftedEventDetail after event update returned error: %v", err)
	}
	if upsertCallCount != 4 {
		t.Fatalf("expected synced proposed dates to be resynced after event update, got %d total calls", upsertCallCount)
	}
	if updatedDetail.SyncStatus != value.SyncStatusSynced {
		t.Fatalf("expected updated event to be synced, got %s", updatedDetail.SyncStatus)
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
	googleAvailable := false

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
				if !googleAvailable {
					return "", errors.New("google unavailable")
				}
				return "google-created-after-retry", nil
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

	googleAvailable = true
	retriedDetail, err := uc.FetchDraftedEventDetail(ctx, userID, "user@example.com", eventID)
	if err != nil {
		t.Fatalf("retry FetchDraftedEventDetail returned error: %v", err)
	}
	if upsertCallCount != 2 {
		t.Fatalf("expected one failed and one successful upsert, got %d", upsertCallCount)
	}
	if retriedDetail.SyncStatus != value.SyncStatusSynced || retriedDetail.LastSyncError != nil {
		t.Fatalf("expected event sync error to be cleared after retry, got %+v", retriedDetail)
	}
	if len(retriedDetail.ProposedDates) != 1 || retriedDetail.ProposedDates[0].SyncStatus != value.SyncStatusSynced {
		t.Fatalf("expected proposed date to be synced after retry, got %#v", retriedDetail.ProposedDates)
	}
	if retriedDetail.ProposedDates[0].GoogleEventID == nil || *retriedDetail.ProposedDates[0].GoogleEventID != "google-created-after-retry" {
		t.Fatalf("unexpected google event id after retry: %#v", retriedDetail.ProposedDates[0].GoogleEventID)
	}
	if retriedDetail.ProposedDates[0].LastSyncError != nil {
		t.Fatalf("expected proposed date sync error to be cleared after retry, got %#v", retriedDetail.ProposedDates[0].LastSyncError)
	}

	if _, err := uc.FetchDraftedEventDetail(ctx, userID, "user@example.com", eventID); err != nil {
		t.Fatalf("FetchDraftedEventDetail after successful retry returned error: %v", err)
	}
	if upsertCallCount != 2 {
		t.Fatalf("expected synced proposed date to skip another retry, got %d upsert calls", upsertCallCount)
	}
}

func stringPointer(value string) *string {
	return &value
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

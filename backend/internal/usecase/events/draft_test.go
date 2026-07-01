package events

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	domainProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func TestCreateDraftedEventsMarksSyncPending(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	eventID := uuid.New()
	dateID := uuid.New()
	start := time.Now().UTC().Add(time.Hour)
	end := start.Add(time.Hour)

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: calendarID}, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: uuid.New(), SyncProposedDates: true}, nil
				},
				createEventFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, title, location, description string, gotStart, gotEnd time.Time) (*domainEvent.Event, error) {
					if gotUserID != userID || gotCalendarID != calendarID {
						t.Fatalf("unexpected create event args: %s %s", gotUserID, gotCalendarID)
					}
					return &domainEvent.Event{
						ID:          eventID,
						Title:       title,
						Location:    location,
						Description: description,
						Status:      value.StatusActive,
						SyncStatus:  value.SyncStatusNotSynced,
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					if opt.SyncStatus == nil || *opt.SyncStatus != value.SyncStatusPending {
						t.Fatalf("expected pending sync mutation, got %#v", opt.SyncStatus)
					}
					return &domainEvent.Event{
						ID:          eventID,
						Title:       "Draft",
						Location:    "Tokyo",
						Description: "desc",
						Status:      value.StatusActive,
						SyncStatus:  value.SyncStatusPending,
					}, nil
				},
				createProposedDatesFn: func(ctx context.Context, selectedDates []SelectedDate, gotEventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return []*domainProposedDate.ProposedDate{
						{
							ID:         dateID,
							StartTime:  selectedDates[0].Start,
							EndTime:    selectedDates[0].End,
							Priority:   selectedDates[0].Priority,
							Status:     value.ProposedDateStatusActive,
							SyncStatus: value.SyncStatusNotSynced,
						},
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
					if id != dateID {
						t.Fatalf("unexpected date id: %s", id)
					}
					if opt.SyncStatus == nil || *opt.SyncStatus != value.SyncStatusPending {
						t.Fatalf("expected pending proposed date sync, got %#v", opt.SyncStatus)
					}
					return &domainProposedDate.ProposedDate{
						ID:         dateID,
						StartTime:  start,
						EndTime:    end,
						Priority:   10,
						Status:     value.ProposedDateStatusActive,
						SyncStatus: value.SyncStatusPending,
					}, nil
				},
			},
		},
		nil,
	)

	response, err := uc.CreateDraftedEvents(ctx, userID, "user@example.com", DraftCreationRequest{
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "desc",
		SelectedDates: []SelectedDate{
			{Start: start, End: end, Priority: 10},
		},
	})
	if err != nil {
		t.Fatalf("CreateDraftedEvents returned error: %v", err)
	}
	if response.SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected event sync status: %s", response.SyncStatus)
	}
	if len(response.ProposedDates) != 1 || response.ProposedDates[0].SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected proposed dates: %#v", response.ProposedDates)
	}
}

func TestCreateDraftedEventsKeepsNotSyncedWhenCandidateSyncDisabled(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	eventID := uuid.New()
	dateID := uuid.New()
	start := time.Now().UTC().Add(time.Hour)
	end := start.Add(time.Hour)

	var updateEventCalled bool
	var updateProposedDateCalled bool

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: calendarID}, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: uuid.New(), SyncProposedDates: false}, nil
				},
				createEventFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, title, location, description string, gotStart, gotEnd time.Time) (*domainEvent.Event, error) {
					if gotUserID != userID || gotCalendarID != calendarID {
						t.Fatalf("unexpected create event args: %s %s", gotUserID, gotCalendarID)
					}
					return &domainEvent.Event{
						ID:          eventID,
						Title:       title,
						Location:    location,
						Description: description,
						Status:      value.StatusActive,
						SyncStatus:  value.SyncStatusNotSynced,
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					updateEventCalled = true
					t.Fatalf("UpdateEvent should not be called when candidate sync is disabled")
					return nil, nil
				},
				createProposedDatesFn: func(ctx context.Context, selectedDates []SelectedDate, gotEventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return []*domainProposedDate.ProposedDate{
						{
							ID:         dateID,
							StartTime:  selectedDates[0].Start,
							EndTime:    selectedDates[0].End,
							Priority:   selectedDates[0].Priority,
							Status:     value.ProposedDateStatusActive,
							SyncStatus: value.SyncStatusNotSynced,
						},
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
					updateProposedDateCalled = true
					t.Fatalf("UpdateProposedDate should not be called when candidate sync is disabled")
					return nil, nil
				},
			},
		},
		nil,
	)

	response, err := uc.CreateDraftedEvents(ctx, userID, "user@example.com", DraftCreationRequest{
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "desc",
		SelectedDates: []SelectedDate{
			{Start: start, End: end, Priority: 10},
		},
	})
	if err != nil {
		t.Fatalf("CreateDraftedEvents returned error: %v", err)
	}
	if updateEventCalled || updateProposedDateCalled {
		t.Fatalf("unexpected sync mutation calls: updateEvent=%v updateProposedDate=%v", updateEventCalled, updateProposedDateCalled)
	}
	if response.SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected event sync status: %s", response.SyncStatus)
	}
	if len(response.ProposedDates) != 1 || response.ProposedDates[0].SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected proposed dates: %#v", response.ProposedDates)
	}
}

func TestCreateDraftedEventsKeepsNotSyncedWhenCandidateCalendarMissing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	eventID := uuid.New()
	dateID := uuid.New()
	start := time.Now().UTC().Add(time.Hour)
	end := start.Add(time.Hour)

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: calendarID}, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return nil, repoerr.ErrNotFound
				},
				createEventFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, title, location, description string, gotStart, gotEnd time.Time) (*domainEvent.Event, error) {
					return &domainEvent.Event{
						ID:          eventID,
						Title:       title,
						Location:    location,
						Description: description,
						Status:      value.StatusActive,
						SyncStatus:  value.SyncStatusNotSynced,
					}, nil
				},
				createProposedDatesFn: func(ctx context.Context, selectedDates []SelectedDate, gotEventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
					return []*domainProposedDate.ProposedDate{
						{
							ID:         dateID,
							StartTime:  selectedDates[0].Start,
							EndTime:    selectedDates[0].End,
							Priority:   selectedDates[0].Priority,
							Status:     value.ProposedDateStatusActive,
							SyncStatus: value.SyncStatusNotSynced,
						},
					}, nil
				},
			},
		},
		nil,
	)

	response, err := uc.CreateDraftedEvents(ctx, userID, "user@example.com", DraftCreationRequest{
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "desc",
		SelectedDates: []SelectedDate{
			{Start: start, End: end, Priority: 10},
		},
	})
	if err != nil {
		t.Fatalf("CreateDraftedEvents returned error: %v", err)
	}
	if response.SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected event sync status: %s", response.SyncStatus)
	}
	if len(response.ProposedDates) != 1 || response.ProposedDates[0].SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected proposed dates: %#v", response.ProposedDates)
	}
}

func TestUpdateDraftedEventsMarksPendingSyncForDraftEdits(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	eventID := uuid.New()
	dateID := uuid.New()
	start := time.Now().UTC().Add(2 * time.Hour)
	end := start.Add(time.Hour)

	var eventMutation EventMutation
	var dateMutation ProposedDateMutation

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: calendarID}, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: uuid.New(), SyncProposedDates: true}, nil
				},
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
					if gotUserID != userID || gotEventID != eventID {
						t.Fatalf("unexpected find event args: %s %s", gotUserID, gotEventID)
					}
					return &domainEvent.Event{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            value.StatusActive,
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					eventMutation = opt
					return &domainEvent.Event{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            value.StatusActive,
					}, nil
				},
				listProposedDatesByEventFn: func(ctx context.Context, gotEventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return []*domainProposedDate.ProposedDate{
						{
							ID:         dateID,
							StartTime:  start.Add(-time.Hour),
							EndTime:    end.Add(-time.Hour),
							Priority:   1,
							Status:     value.ProposedDateStatusActive,
							SyncStatus: value.SyncStatusSynced,
						},
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
					if id != dateID {
						t.Fatalf("unexpected date id: %s", id)
					}
					dateMutation = opt
					return &domainProposedDate.ProposedDate{ID: dateID}, nil
				},
				deleteProposedDateFn: func(ctx context.Context, id uuid.UUID) error {
					t.Fatalf("DeleteProposedDate should not be called")
					return nil
				},
				createProposedDateFn: func(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*domainProposedDate.ProposedDate, error) {
					t.Fatalf("CreateProposedDate should not be called")
					return nil, nil
				},
			},
		},
		nil,
	)

	err := uc.UpdateDraftedEvents(ctx, userID, eventID, "user@example.com", DraftUpdateRequest{
		Title:       "Updated title",
		Location:    "Osaka",
		Description: "updated desc",
		Status:      value.StatusActive,
		ProposedDates: []ProposedDateRequest{
			{
				ID:       &dateID,
				Start:    &start,
				End:      &end,
				Priority: 2,
			},
		},
	})
	if err != nil {
		t.Fatalf("UpdateDraftedEvents returned error: %v", err)
	}
	if eventMutation.Status == nil || *eventMutation.Status != value.StatusActive {
		t.Fatalf("unexpected event status mutation: %#v", eventMutation.Status)
	}
	if eventMutation.SyncStatus == nil || *eventMutation.SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected event sync mutation: %#v", eventMutation.SyncStatus)
	}
	if dateMutation.SyncStatus == nil || *dateMutation.SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected proposed date sync mutation: %#v", dateMutation.SyncStatus)
	}
}

func TestUpdateDraftedEventsKeepsNotSyncedWhenCandidateSyncDisabled(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	eventID := uuid.New()
	keepDateID := uuid.New()
	deleteDateID := uuid.New()
	createdDateID := uuid.New()
	start := time.Now().UTC().Add(2 * time.Hour)
	end := start.Add(time.Hour)
	newStart := start.Add(2 * time.Hour)
	newEnd := newStart.Add(time.Hour)

	var eventMutation EventMutation
	mutations := make(map[uuid.UUID]ProposedDateMutation)
	var createdMutation ProposedDateMutation
	var deletedID uuid.UUID

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: calendarID}, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: uuid.New(), SyncProposedDates: false}, nil
				},
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
					if gotUserID != userID || gotEventID != eventID {
						t.Fatalf("unexpected find event args: %s %s", gotUserID, gotEventID)
					}
					return &domainEvent.Event{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            value.StatusActive,
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					eventMutation = opt
					return &domainEvent.Event{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            value.StatusActive,
						SyncStatus:        value.SyncStatusNotSynced,
					}, nil
				},
				listProposedDatesByEventFn: func(ctx context.Context, gotEventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return []*domainProposedDate.ProposedDate{
						{
							ID:         keepDateID,
							StartTime:  start.Add(-time.Hour),
							EndTime:    end.Add(-time.Hour),
							Priority:   10,
							Status:     value.ProposedDateStatusActive,
							SyncStatus: value.SyncStatusSynced,
						},
						{
							ID:         deleteDateID,
							StartTime:  start.Add(-2 * time.Hour),
							EndTime:    end.Add(-2 * time.Hour),
							Priority:   5,
							Status:     value.ProposedDateStatusActive,
							SyncStatus: value.SyncStatusSynced,
						},
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
					mutations[id] = opt
					return &domainProposedDate.ProposedDate{ID: id}, nil
				},
				deleteProposedDateFn: func(ctx context.Context, id uuid.UUID) error {
					deletedID = id
					return nil
				},
				createProposedDateFn: func(ctx context.Context, opt ProposedDateMutation, gotEventID uuid.UUID) (*domainProposedDate.ProposedDate, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					createdMutation = opt
					return &domainProposedDate.ProposedDate{ID: createdDateID}, nil
				},
			},
		},
		nil,
	)

	err := uc.UpdateDraftedEvents(ctx, userID, eventID, "user@example.com", DraftUpdateRequest{
		Title:       "Updated title",
		Location:    "Osaka",
		Description: "updated desc",
		Status:      value.StatusActive,
		ProposedDates: []ProposedDateRequest{
			{
				ID:       &keepDateID,
				Start:    &start,
				End:      &end,
				Priority: 20,
			},
			{
				Start:    &newStart,
				End:      &newEnd,
				Priority: 15,
			},
		},
	})
	if err != nil {
		t.Fatalf("UpdateDraftedEvents returned error: %v", err)
	}
	if eventMutation.Status == nil || *eventMutation.Status != value.StatusActive {
		t.Fatalf("unexpected event status mutation: %#v", eventMutation.Status)
	}
	if eventMutation.SyncStatus == nil || *eventMutation.SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected event sync mutation: %#v", eventMutation.SyncStatus)
	}
	if !eventMutation.ClearLastSyncError {
		t.Fatalf("expected event last sync error to be cleared")
	}
	if deletedID != deleteDateID {
		t.Fatalf("unexpected deleted id: %s", deletedID)
	}
	if mutations[keepDateID].SyncStatus == nil || *mutations[keepDateID].SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected kept proposed date sync mutation: %#v", mutations[keepDateID].SyncStatus)
	}
	if !mutations[keepDateID].ClearLastSyncError {
		t.Fatalf("expected kept proposed date last sync error to be cleared")
	}
	if mutations[deleteDateID].SyncStatus == nil || *mutations[deleteDateID].SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected deleted proposed date sync mutation: %#v", mutations[deleteDateID].SyncStatus)
	}
	if createdMutation.SyncStatus == nil || *createdMutation.SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected created proposed date sync mutation: %#v", createdMutation.SyncStatus)
	}
}

func TestUpdateDraftedEventsKeepsNotSyncedWhenCandidateCalendarMissing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	eventID := uuid.New()
	dateID := uuid.New()
	start := time.Now().UTC().Add(2 * time.Hour)
	end := start.Add(time.Hour)

	var eventMutation EventMutation
	var dateMutation ProposedDateMutation

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: calendarID}, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return nil, repoerr.ErrNotFound
				},
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return &domainEvent.Event{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            value.StatusActive,
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					eventMutation = opt
					return &domainEvent.Event{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            value.StatusActive,
						SyncStatus:        value.SyncStatusNotSynced,
					}, nil
				},
				listProposedDatesByEventFn: func(ctx context.Context, gotEventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
					return []*domainProposedDate.ProposedDate{
						{
							ID:         dateID,
							StartTime:  start.Add(-time.Hour),
							EndTime:    end.Add(-time.Hour),
							Priority:   10,
							Status:     value.ProposedDateStatusActive,
							SyncStatus: value.SyncStatusSynced,
						},
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
					dateMutation = opt
					return &domainProposedDate.ProposedDate{ID: id}, nil
				},
				deleteProposedDateFn: func(ctx context.Context, id uuid.UUID) error {
					t.Fatalf("DeleteProposedDate should not be called")
					return nil
				},
				createProposedDateFn: func(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*domainProposedDate.ProposedDate, error) {
					t.Fatalf("CreateProposedDate should not be called")
					return nil, nil
				},
			},
		},
		nil,
	)

	err := uc.UpdateDraftedEvents(ctx, userID, eventID, "user@example.com", DraftUpdateRequest{
		Title:       "Updated title",
		Location:    "Osaka",
		Description: "updated desc",
		Status:      value.StatusActive,
		ProposedDates: []ProposedDateRequest{
			{
				ID:       &dateID,
				Start:    &start,
				End:      &end,
				Priority: 20,
			},
		},
	})
	if err != nil {
		t.Fatalf("UpdateDraftedEvents returned error: %v", err)
	}
	if eventMutation.SyncStatus == nil || *eventMutation.SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected event sync mutation: %#v", eventMutation.SyncStatus)
	}
	if dateMutation.SyncStatus == nil || *dateMutation.SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected proposed date sync mutation: %#v", dateMutation.SyncStatus)
	}
}

func TestDeleteDraftedEventsMarksPendingBeforeSoftDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	eventID := uuid.New()

	var updateCalled bool
	var deleteCalled bool

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: calendarID}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					updateCalled = true
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					if opt.SyncStatus == nil || *opt.SyncStatus != value.SyncStatusPending {
						t.Fatalf("unexpected sync status mutation: %#v", opt.SyncStatus)
					}
					return &domainEvent.Event{ID: eventID}, nil
				},
				softDeleteEventFn: func(ctx context.Context, id uuid.UUID) error {
					deleteCalled = true
					if !updateCalled {
						t.Fatalf("soft delete called before sync status update")
					}
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					return nil
				},
			},
		},
		nil,
	)

	err := uc.DeleteDraftedEvents(ctx, userID, "user@example.com", eventID)
	if err != nil {
		t.Fatalf("DeleteDraftedEvents returned error: %v", err)
	}
	if !updateCalled || !deleteCalled {
		t.Fatalf("expected update and soft delete to be called")
	}
}

func TestUpdateDraftedEventsMarksDeletedProposedDatesPendingBeforeSoftDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	eventID := uuid.New()
	keepDateID := uuid.New()
	deleteDateID := uuid.New()
	start := time.Now().UTC().Add(4 * time.Hour)
	end := start.Add(time.Hour)

	mutations := make(map[uuid.UUID]ProposedDateMutation)
	var deletedID uuid.UUID

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: calendarID}, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*EventCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &EventCalendar{ID: uuid.New(), SyncProposedDates: true}, nil
				},
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return &domainEvent.Event{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            value.StatusActive,
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
					return &domainEvent.Event{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            value.StatusActive,
					}, nil
				},
				listProposedDatesByEventFn: func(ctx context.Context, gotEventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
					return []*domainProposedDate.ProposedDate{
						{ID: keepDateID, StartTime: start.Add(-time.Hour), EndTime: end.Add(-time.Hour), Priority: 1},
						{ID: deleteDateID, StartTime: start.Add(-2 * time.Hour), EndTime: end.Add(-2 * time.Hour), Priority: 2},
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
					mutations[id] = opt
					return &domainProposedDate.ProposedDate{ID: id}, nil
				},
				deleteProposedDateFn: func(ctx context.Context, id uuid.UUID) error {
					deletedID = id
					return nil
				},
				createProposedDateFn: func(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*domainProposedDate.ProposedDate, error) {
					t.Fatalf("CreateProposedDate should not be called")
					return nil, nil
				},
			},
		},
		nil,
	)

	err := uc.UpdateDraftedEvents(ctx, userID, eventID, "user@example.com", DraftUpdateRequest{
		Title:       "Updated title",
		Location:    "Osaka",
		Description: "updated desc",
		Status:      value.StatusActive,
		ProposedDates: []ProposedDateRequest{
			{
				ID:       &keepDateID,
				Start:    &start,
				End:      &end,
				Priority: 1,
			},
		},
	})
	if err != nil {
		t.Fatalf("UpdateDraftedEvents returned error: %v", err)
	}
	if deletedID != deleteDateID {
		t.Fatalf("unexpected deleted id: %s", deletedID)
	}
	deleteMutation, ok := mutations[deleteDateID]
	if !ok {
		t.Fatalf("expected pending mutation for deleted proposed date")
	}
	if deleteMutation.SyncStatus == nil || *deleteMutation.SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected deleted proposed date sync mutation: %#v", deleteMutation.SyncStatus)
	}
}

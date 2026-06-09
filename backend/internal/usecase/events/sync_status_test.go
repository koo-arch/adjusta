package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repositorymodel"
)

type fakeEventTxStore struct {
	t                          *testing.T
	findPrimaryCalendarFn      func(ctx context.Context, userID uuid.UUID) (*repositorymodel.StoredCalendar, error)
	findEventBySlugFn          func(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error)
	readCalendarFn             func(ctx context.Context, calendarID uuid.UUID) (*repositorymodel.StoredCalendar, error)
	createEventFn              func(ctx context.Context, userID, primaryCalendarID uuid.UUID, title, location, description string, start, end time.Time) (*repositorymodel.StoredEvent, error)
	updateEventFn              func(ctx context.Context, id uuid.UUID, opt EventMutation) (*repositorymodel.StoredEvent, error)
	softDeleteEventFn          func(ctx context.Context, id uuid.UUID) error
	listProposedDatesByEventFn func(ctx context.Context, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error)
	createProposedDatesFn      func(ctx context.Context, selectedDates []appmodel.SelectedDate, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error)
	updateProposedDateFn       func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*repositorymodel.StoredProposedDate, error)
	deleteProposedDateFn       func(ctx context.Context, id uuid.UUID) error
	createProposedDateFn       func(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*repositorymodel.StoredProposedDate, error)
	decrementPriorityFn        func(ctx context.Context, eventID, excludeID uuid.UUID) error
	reorderPriorityFn          func(ctx context.Context, eventID uuid.UUID) error
}

func (f *fakeEventTxStore) unexpected(method string) {
	f.t.Fatalf("%s should not be called", method)
}

func (f *fakeEventTxStore) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
	if f.findPrimaryCalendarFn == nil {
		f.unexpected("FindPrimaryCalendar")
	}
	return f.findPrimaryCalendarFn(ctx, userID)
}

func (f *fakeEventTxStore) FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
	if f.findEventBySlugFn == nil {
		f.unexpected("FindEventBySlug")
	}
	return f.findEventBySlugFn(ctx, userID, slug, withProposedDates)
}

func (f *fakeEventTxStore) ReadCalendar(ctx context.Context, calendarID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
	if f.readCalendarFn == nil {
		f.unexpected("ReadCalendar")
	}
	return f.readCalendarFn(ctx, calendarID)
}

func (f *fakeEventTxStore) CreateEvent(ctx context.Context, userID, primaryCalendarID uuid.UUID, title, location, description string, start, end time.Time) (*repositorymodel.StoredEvent, error) {
	if f.createEventFn == nil {
		f.unexpected("CreateEvent")
	}
	return f.createEventFn(ctx, userID, primaryCalendarID, title, location, description, start, end)
}

func (f *fakeEventTxStore) UpdateEvent(ctx context.Context, id uuid.UUID, opt EventMutation) (*repositorymodel.StoredEvent, error) {
	if f.updateEventFn == nil {
		f.unexpected("UpdateEvent")
	}
	return f.updateEventFn(ctx, id, opt)
}

func (f *fakeEventTxStore) SoftDeleteEvent(ctx context.Context, id uuid.UUID) error {
	if f.softDeleteEventFn == nil {
		f.unexpected("SoftDeleteEvent")
	}
	return f.softDeleteEventFn(ctx, id)
}

func (f *fakeEventTxStore) ListProposedDatesByEvent(ctx context.Context, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
	if f.listProposedDatesByEventFn == nil {
		f.unexpected("ListProposedDatesByEvent")
	}
	return f.listProposedDatesByEventFn(ctx, eventID)
}

func (f *fakeEventTxStore) CreateProposedDates(ctx context.Context, selectedDates []appmodel.SelectedDate, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
	if f.createProposedDatesFn == nil {
		f.unexpected("CreateProposedDates")
	}
	return f.createProposedDatesFn(ctx, selectedDates, eventID)
}

func (f *fakeEventTxStore) UpdateProposedDate(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*repositorymodel.StoredProposedDate, error) {
	if f.updateProposedDateFn == nil {
		f.unexpected("UpdateProposedDate")
	}
	return f.updateProposedDateFn(ctx, id, opt)
}

func (f *fakeEventTxStore) DeleteProposedDate(ctx context.Context, id uuid.UUID) error {
	if f.deleteProposedDateFn == nil {
		f.unexpected("DeleteProposedDate")
	}
	return f.deleteProposedDateFn(ctx, id)
}

func (f *fakeEventTxStore) CreateProposedDate(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*repositorymodel.StoredProposedDate, error) {
	if f.createProposedDateFn == nil {
		f.unexpected("CreateProposedDate")
	}
	return f.createProposedDateFn(ctx, opt, eventID)
}

func (f *fakeEventTxStore) DecrementPriorityExceptID(ctx context.Context, eventID, excludeID uuid.UUID) error {
	if f.decrementPriorityFn == nil {
		f.unexpected("DecrementPriorityExceptID")
	}
	return f.decrementPriorityFn(ctx, eventID, excludeID)
}

func (f *fakeEventTxStore) ReorderPriority(ctx context.Context, eventID uuid.UUID) error {
	if f.reorderPriorityFn == nil {
		f.unexpected("ReorderPriority")
	}
	return f.reorderPriorityFn(ctx, eventID)
}

type fakeEventTransaction struct {
	store EventTxStore
}

func (f *fakeEventTransaction) Do(ctx context.Context, fn func(store EventTxStore) error) error {
	return fn(f.store)
}

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
		nil,
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &repositorymodel.StoredCalendar{ID: calendarID}, nil
				},
				createEventFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, title, location, description string, gotStart, gotEnd time.Time) (*repositorymodel.StoredEvent, error) {
					if gotUserID != userID || gotCalendarID != calendarID {
						t.Fatalf("unexpected create event args: %s %s", gotUserID, gotCalendarID)
					}
					return &repositorymodel.StoredEvent{
						ID:          eventID,
						Title:       title,
						Location:    location,
						Description: description,
						Status:      domainvalue.StatusActive,
						SyncStatus:  domainvalue.SyncStatusNotSynced,
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*repositorymodel.StoredEvent, error) {
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					if opt.SyncStatus == nil || *opt.SyncStatus != domainvalue.SyncStatusPending {
						t.Fatalf("expected pending sync mutation, got %#v", opt.SyncStatus)
					}
					return &repositorymodel.StoredEvent{
						ID:          eventID,
						Title:       "Draft",
						Location:    "Tokyo",
						Description: "desc",
						Status:      domainvalue.StatusActive,
						SyncStatus:  domainvalue.SyncStatusPending,
					}, nil
				},
				createProposedDatesFn: func(ctx context.Context, selectedDates []appmodel.SelectedDate, gotEventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return []*repositorymodel.StoredProposedDate{
						{
							ID:         dateID,
							StartTime:  selectedDates[0].Start,
							EndTime:    selectedDates[0].End,
							Priority:   selectedDates[0].Priority,
							Status:     domainvalue.ProposedDateStatusActive,
							SyncStatus: domainvalue.SyncStatusNotSynced,
						},
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*repositorymodel.StoredProposedDate, error) {
					if id != dateID {
						t.Fatalf("unexpected date id: %s", id)
					}
					if opt.SyncStatus == nil || *opt.SyncStatus != domainvalue.SyncStatusPending {
						t.Fatalf("expected pending proposed date sync, got %#v", opt.SyncStatus)
					}
					return &repositorymodel.StoredProposedDate{
						ID:         dateID,
						StartTime:  start,
						EndTime:    end,
						Priority:   10,
						Status:     domainvalue.ProposedDateStatusActive,
						SyncStatus: domainvalue.SyncStatusPending,
					}, nil
				},
			},
		},
		nil,
	)

	response, err := uc.CreateDraftedEvents(ctx, userID, "user@example.com", &appmodel.EventDraftCreation{
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "desc",
		SelectedDates: []appmodel.SelectedDate{
			{Start: start, End: end, Priority: 10},
		},
	})
	if err != nil {
		t.Fatalf("CreateDraftedEvents returned error: %v", err)
	}
	if response.SyncStatus != domainvalue.SyncStatusPending {
		t.Fatalf("unexpected event sync status: %s", response.SyncStatus)
	}
	if len(response.ProposedDates) != 1 || response.ProposedDates[0].SyncStatus != domainvalue.SyncStatusPending {
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
		nil,
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &repositorymodel.StoredCalendar{ID: calendarID}, nil
				},
				findEventBySlugFn: func(ctx context.Context, gotUserID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
					if gotUserID != userID || slug != "draft-event" {
						t.Fatalf("unexpected find event args: %s %s", gotUserID, slug)
					}
					return &repositorymodel.StoredEvent{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            domainvalue.StatusActive,
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*repositorymodel.StoredEvent, error) {
					eventMutation = opt
					return &repositorymodel.StoredEvent{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            domainvalue.StatusActive,
					}, nil
				},
				listProposedDatesByEventFn: func(ctx context.Context, gotEventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return []*repositorymodel.StoredProposedDate{
						{
							ID:         dateID,
							StartTime:  start.Add(-time.Hour),
							EndTime:    end.Add(-time.Hour),
							Priority:   1,
							Status:     domainvalue.ProposedDateStatusActive,
							SyncStatus: domainvalue.SyncStatusSynced,
						},
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*repositorymodel.StoredProposedDate, error) {
					if id != dateID {
						t.Fatalf("unexpected date id: %s", id)
					}
					dateMutation = opt
					return &repositorymodel.StoredProposedDate{ID: dateID}, nil
				},
				deleteProposedDateFn: func(ctx context.Context, id uuid.UUID) error {
					t.Fatalf("DeleteProposedDate should not be called")
					return nil
				},
				createProposedDateFn: func(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*repositorymodel.StoredProposedDate, error) {
					t.Fatalf("CreateProposedDate should not be called")
					return nil, nil
				},
			},
		},
		nil,
	)

	err := uc.UpdateDraftedEvents(ctx, userID, "draft-event", "user@example.com", &appmodel.EventDraftUpdate{
		Title:       "Updated title",
		Location:    "Osaka",
		Description: "updated desc",
		Status:      domainvalue.StatusActive,
		ProposedDates: []appmodel.ProposedDate{
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
	if eventMutation.Status == nil || *eventMutation.Status != domainvalue.StatusActive {
		t.Fatalf("unexpected event status mutation: %#v", eventMutation.Status)
	}
	if eventMutation.SyncStatus == nil || *eventMutation.SyncStatus != domainvalue.SyncStatusPending {
		t.Fatalf("unexpected event sync mutation: %#v", eventMutation.SyncStatus)
	}
	if dateMutation.SyncStatus == nil || *dateMutation.SyncStatus != domainvalue.SyncStatusPending {
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
		nil,
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &repositorymodel.StoredCalendar{ID: calendarID}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*repositorymodel.StoredEvent, error) {
					updateCalled = true
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					if opt.SyncStatus == nil || *opt.SyncStatus != domainvalue.SyncStatusPending {
						t.Fatalf("unexpected sync status mutation: %#v", opt.SyncStatus)
					}
					return &repositorymodel.StoredEvent{ID: eventID}, nil
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

	err := uc.DeleteDraftedEvents(ctx, userID, "user@example.com", &appmodel.EventDraftDetail{ID: eventID})
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
		nil,
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &repositorymodel.StoredCalendar{ID: calendarID}, nil
				},
				findEventBySlugFn: func(ctx context.Context, gotUserID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
					return &repositorymodel.StoredEvent{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            domainvalue.StatusActive,
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*repositorymodel.StoredEvent, error) {
					return &repositorymodel.StoredEvent{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Status:            domainvalue.StatusActive,
					}, nil
				},
				listProposedDatesByEventFn: func(ctx context.Context, gotEventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
					return []*repositorymodel.StoredProposedDate{
						{ID: keepDateID, StartTime: start.Add(-time.Hour), EndTime: end.Add(-time.Hour), Priority: 1},
						{ID: deleteDateID, StartTime: start.Add(-2 * time.Hour), EndTime: end.Add(-2 * time.Hour), Priority: 2},
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*repositorymodel.StoredProposedDate, error) {
					mutations[id] = opt
					return &repositorymodel.StoredProposedDate{ID: id}, nil
				},
				deleteProposedDateFn: func(ctx context.Context, id uuid.UUID) error {
					deletedID = id
					return nil
				},
				createProposedDateFn: func(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*repositorymodel.StoredProposedDate, error) {
					t.Fatalf("CreateProposedDate should not be called")
					return nil, nil
				},
			},
		},
		nil,
	)

	err := uc.UpdateDraftedEvents(ctx, userID, "draft-event", "user@example.com", &appmodel.EventDraftUpdate{
		Title:       "Updated title",
		Location:    "Osaka",
		Description: "updated desc",
		Status:      domainvalue.StatusActive,
		ProposedDates: []appmodel.ProposedDate{
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
	if deleteMutation.SyncStatus == nil || *deleteMutation.SyncStatus != domainvalue.SyncStatusPending {
		t.Fatalf("unexpected deleted proposed date sync mutation: %#v", deleteMutation.SyncStatus)
	}
}

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
		nil,
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findEventBySlugFn: func(ctx context.Context, gotUserID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
					if gotUserID != userID || slug != "finalize-event" {
						t.Fatalf("unexpected find event args: %s %s", gotUserID, slug)
					}
					return &repositorymodel.StoredEvent{
						ID:                eventID,
						PrimaryCalendarID: calendarID,
						Title:             "Finalize",
						Status:            domainvalue.StatusActive,
					}, nil
				},
				readCalendarFn: func(ctx context.Context, gotCalendarID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
					if gotCalendarID != calendarID {
						t.Fatalf("unexpected calendar id: %s", gotCalendarID)
					}
					return &repositorymodel.StoredCalendar{
						ID:               calendarID,
						GoogleCalendarID: "primary-calendar",
					}, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*repositorymodel.StoredEvent, error) {
					if id != eventID {
						t.Fatalf("unexpected event id: %s", id)
					}
					failureMutation = opt
					return &repositorymodel.StoredEvent{ID: eventID}, nil
				},
			},
		},
		&fakeGoogleCalendarGateway{
			fetchEventsFn: func(ctx context.Context, userID uuid.UUID, calendars []*repositorymodel.StoredCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
				t.Fatalf("FetchEvents should not be called")
				return nil, nil
			},
			upsertEventFn: func(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
				return "", errors.New("google unavailable")
			},
		},
	)

	err := uc.FinalizeProposedDate(ctx, userID, "finalize-event", "user@example.com", &appmodel.ConfirmEvent{
		ConfirmDate: appmodel.ConfirmDate{
			Start: &start,
			End:   &end,
		},
	})
	if !internalErrors.IsKind(err, internalErrors.KindInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
	if failureMutation.SyncStatus == nil || *failureMutation.SyncStatus != domainvalue.SyncStatusFailed {
		t.Fatalf("unexpected failure sync status: %#v", failureMutation.SyncStatus)
	}
	if failureMutation.LastSyncError == nil || *failureMutation.LastSyncError == "" {
		t.Fatalf("expected last sync error to be recorded, got %#v", failureMutation.LastSyncError)
	}
}

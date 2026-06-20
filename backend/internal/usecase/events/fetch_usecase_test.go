package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

type fakeEventReader struct {
	findPrimaryCalendarFn          func(ctx context.Context, userID uuid.UUID) (*CalendarRecord, error)
	findAdjustaCandidateCalendarFn func(ctx context.Context, userID uuid.UUID) (*CalendarRecord, error)
	listCalendarsByUserFn          func(ctx context.Context, userID uuid.UUID) ([]*CalendarRecord, error)
	searchEventsFn                 func(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*EventRecord, error)
	findEventByIDFn                func(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*EventRecord, error)
}

func (f *fakeEventReader) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*CalendarRecord, error) {
	return f.findPrimaryCalendarFn(ctx, userID)
}

func (f *fakeEventReader) FindAdjustaCandidateCalendar(ctx context.Context, userID uuid.UUID) (*CalendarRecord, error) {
	if f.findAdjustaCandidateCalendarFn == nil {
		panic("FindAdjustaCandidateCalendar should not be called")
	}
	return f.findAdjustaCandidateCalendarFn(ctx, userID)
}

func (f *fakeEventReader) ListCalendarsByUser(ctx context.Context, userID uuid.UUID) ([]*CalendarRecord, error) {
	return f.listCalendarsByUserFn(ctx, userID)
}

func (f *fakeEventReader) SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*EventRecord, error) {
	return f.searchEventsFn(ctx, userID, calendarID, opt)
}

func (f *fakeEventReader) FindEventByID(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*EventRecord, error) {
	return f.findEventByIDFn(ctx, userID, eventID, withProposedDates)
}

type fakeGoogleCalendarGateway struct {
	fetchEventsFn func(ctx context.Context, userID uuid.UUID, calendars []*CalendarRecord, startTime, endTime time.Time) (*GoogleEventFetchResult, error)
	upsertEventFn func(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
}

func (f *fakeGoogleCalendarGateway) FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*CalendarRecord, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
	return f.fetchEventsFn(ctx, userID, calendars, startTime, endTime)
}

func (f *fakeGoogleCalendarGateway) UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
	return f.upsertEventFn(ctx, userID, calendarID, existingGoogleEventID, title, location, description, start, end)
}

func fakeReposFromReader(reader *fakeEventReader) EventTxRepositories {
	calendarRepo := &fakeCalendarRepository{reader: reader, calendars: map[uuid.UUID]*repoCalendar.Calendar{}}
	return EventTxRepositories{
		Calendar:     calendarRepo,
		Event:        &fakeEventRepository{reader: reader},
		ProposedDate: &fakeProposedDateRepository{},
		UserCalendar: &fakeUserCalendarRepository{reader: reader, calendarRepo: calendarRepo},
	}
}

func fakeReposFromTxStore(store *fakeEventTxStore) EventTxRepositories {
	calendarRepo := &fakeCalendarRepository{store: store, calendars: map[uuid.UUID]*repoCalendar.Calendar{}}
	return EventTxRepositories{
		Calendar:     calendarRepo,
		Event:        &fakeEventRepository{store: store},
		ProposedDate: &fakeProposedDateRepository{store: store},
		UserCalendar: &fakeUserCalendarRepository{store: store, calendarRepo: calendarRepo},
	}
}

type fakeCalendarRepository struct {
	reader    *fakeEventReader
	store     *fakeEventTxStore
	calendars map[uuid.UUID]*repoCalendar.Calendar
}

func (f *fakeCalendarRepository) Read(ctx context.Context, id uuid.UUID) (*repoCalendar.Calendar, error) {
	if f.store != nil && f.store.readCalendarFn != nil {
		record, err := f.store.readCalendarFn(ctx, id)
		return calendarRecordToDomain(record), err
	}
	if calendar, ok := f.calendars[id]; ok {
		return calendar, nil
	}
	return nil, errors.New("calendar not found")
}

func (f *fakeCalendarRepository) FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*repoCalendar.Calendar, error) {
	if f.reader == nil || f.reader.listCalendarsByUserFn == nil {
		return nil, errors.New("FilterByUserID should not be called")
	}
	records, err := f.reader.listCalendarsByUserFn(ctx, userID)
	if err != nil {
		return nil, err
	}
	calendars := make([]*repoCalendar.Calendar, 0, len(records))
	for _, record := range records {
		calendar := calendarRecordToDomain(record)
		f.calendars[calendar.ID] = calendar
		calendars = append(calendars, calendar)
	}
	return calendars, nil
}

func (f *fakeCalendarRepository) FindByFields(ctx context.Context, userID uuid.UUID, opt repoCalendar.CalendarQueryOptions) (*repoCalendar.Calendar, error) {
	if f.reader != nil && f.reader.findPrimaryCalendarFn != nil {
		record, err := f.reader.findPrimaryCalendarFn(ctx, userID)
		calendar := calendarRecordToDomain(record)
		if calendar != nil {
			f.calendars[calendar.ID] = calendar
		}
		return calendar, err
	}
	if f.store != nil && f.store.findPrimaryCalendarFn != nil {
		record, err := f.store.findPrimaryCalendarFn(ctx, userID)
		calendar := calendarRecordToDomain(record)
		if calendar != nil {
			f.calendars[calendar.ID] = calendar
		}
		return calendar, err
	}
	return nil, errors.New("FindByFields should not be called")
}

func (f *fakeCalendarRepository) FindByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
	return nil, errors.New("FindByGoogleCalendarID should not be called")
}

func (f *fakeCalendarRepository) FilterByFields(ctx context.Context, userID uuid.UUID, opt repoCalendar.CalendarQueryOptions) ([]*repoCalendar.Calendar, error) {
	return nil, errors.New("FilterByFields should not be called")
}

func (f *fakeCalendarRepository) Create(ctx context.Context, opt repoCalendar.CalendarMutationOptions) (*repoCalendar.Calendar, error) {
	return nil, errors.New("Create should not be called")
}

func (f *fakeCalendarRepository) Update(ctx context.Context, id uuid.UUID, opt repoCalendar.CalendarMutationOptions) (*repoCalendar.Calendar, error) {
	return nil, errors.New("Update should not be called")
}

func (f *fakeCalendarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("Delete should not be called")
}

func (f *fakeCalendarRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return errors.New("SoftDelete should not be called")
}

func (f *fakeCalendarRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return errors.New("Restore should not be called")
}

type fakeEventRepository struct {
	reader *fakeEventReader
	store  *fakeEventTxStore
}

func (f *fakeEventRepository) Read(ctx context.Context, id uuid.UUID, opt repoEvent.EventReadOptions) (*repoEvent.Event, error) {
	return nil, errors.New("Read should not be called")
}

func (f *fakeEventRepository) FilterByCalendarID(ctx context.Context, calendarID uuid.UUID, opt repoEvent.EventFilterOptions) ([]*repoEvent.Event, error) {
	return nil, errors.New("FilterByCalendarID should not be called")
}

func (f *fakeEventRepository) FindByIDAndUser(ctx context.Context, userID, eventID uuid.UUID, opt repoEvent.EventReadOptions) (*repoEvent.Event, error) {
	if f.reader != nil && f.reader.findEventByIDFn != nil {
		return f.reader.findEventByIDFn(ctx, userID, eventID, opt.WithProposedDates)
	}
	if f.store != nil && f.store.findEventByIDFn != nil {
		return f.store.findEventByIDFn(ctx, userID, eventID, opt.WithProposedDates)
	}
	return nil, errors.New("FindByIDAndUser should not be called")
}

func (f *fakeEventRepository) Create(ctx context.Context, userID uuid.UUID, opt repoEvent.EventCreateOptions, primaryCalendarID uuid.UUID) (*repoEvent.Event, error) {
	if f.store == nil || f.store.createEventFn == nil {
		return nil, errors.New("Create should not be called")
	}
	return f.store.createEventFn(ctx, userID, primaryCalendarID, opt.Title, opt.Location, opt.Description, time.Time{}, time.Time{})
}

func (f *fakeEventRepository) Update(ctx context.Context, id uuid.UUID, opt repoEvent.EventUpdateOptions) (*repoEvent.Event, error) {
	if f.store == nil || f.store.updateEventFn == nil {
		return nil, errors.New("Update should not be called")
	}
	return f.store.updateEventFn(ctx, id, opt)
}

func (f *fakeEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("Delete should not be called")
}

func (f *fakeEventRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if f.store == nil || f.store.softDeleteEventFn == nil {
		return errors.New("SoftDelete should not be called")
	}
	return f.store.softDeleteEventFn(ctx, id)
}

func (f *fakeEventRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return errors.New("Restore should not be called")
}

func (f *fakeEventRepository) SearchEvents(ctx context.Context, id, calendarID uuid.UUID, opt repoEvent.EventSearchOptions) ([]*repoEvent.Event, error) {
	if f.reader == nil || f.reader.searchEventsFn == nil {
		return nil, errors.New("SearchEvents should not be called")
	}
	return f.reader.searchEventsFn(ctx, id, calendarID, EventSearchOptions{
		WithProposedDates: opt.WithProposedDates,
		Title:             opt.Title,
		Location:          opt.Location,
		Description:       opt.Description,
		Status:            opt.Status,
		StartTimeGTE:      opt.ProposedDateStartGTE,
		StartTimeLTE:      opt.ProposedDateStartLTE,
		EndTimeGTE:        opt.ProposedDateEndGTE,
		EndTimeLTE:        opt.ProposedDateEndLTE,
		SortBy:            opt.SortBy,
		SortOrder:         opt.SortOrder,
	})
}

type fakeProposedDateRepository struct {
	store *fakeEventTxStore
}

func (f *fakeProposedDateRepository) Read(ctx context.Context, id uuid.UUID) (*repoProposedDate.ProposedDate, error) {
	return nil, errors.New("Read should not be called")
}

func (f *fakeProposedDateRepository) FilterByEventID(ctx context.Context, eventID uuid.UUID) ([]*repoProposedDate.ProposedDate, error) {
	if f.store == nil || f.store.listProposedDatesByEventFn == nil {
		return nil, errors.New("FilterByEventID should not be called")
	}
	return f.store.listProposedDatesByEventFn(ctx, eventID)
}

func (f *fakeProposedDateRepository) Create(ctx context.Context, opt repoProposedDate.ProposedDateCreateOptions, eventID uuid.UUID) (*repoProposedDate.ProposedDate, error) {
	if f.store == nil || f.store.createProposedDateFn == nil {
		return nil, errors.New("Create should not be called")
	}
	updateOpt := repoProposedDate.ProposedDateUpdateOptions{
		GoogleEventID: opt.GoogleEventID,
		StartTime:     &opt.StartTime,
		EndTime:       &opt.EndTime,
		Priority:      &opt.Priority,
		Status:        opt.Status,
		SyncStatus:    opt.SyncStatus,
		LastSyncedAt:  opt.LastSyncedAt,
		LastSyncError: opt.LastSyncError,
	}
	return f.store.createProposedDateFn(ctx, updateOpt, eventID)
}

func (f *fakeProposedDateRepository) Update(ctx context.Context, id uuid.UUID, opt repoProposedDate.ProposedDateUpdateOptions) (*repoProposedDate.ProposedDate, error) {
	if f.store == nil || f.store.updateProposedDateFn == nil {
		return nil, errors.New("Update should not be called")
	}
	return f.store.updateProposedDateFn(ctx, id, opt)
}

func (f *fakeProposedDateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("Delete should not be called")
}

func (f *fakeProposedDateRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if f.store == nil || f.store.deleteProposedDateFn == nil {
		return errors.New("SoftDelete should not be called")
	}
	return f.store.deleteProposedDateFn(ctx, id)
}

func (f *fakeProposedDateRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return errors.New("Restore should not be called")
}

func (f *fakeProposedDateRepository) CreateBulk(ctx context.Context, opts []repoProposedDate.ProposedDateCreateOptions, eventID uuid.UUID) ([]*repoProposedDate.ProposedDate, error) {
	if f.store == nil || f.store.createProposedDatesFn == nil {
		return nil, errors.New("CreateBulk should not be called")
	}
	selectedDates := make([]SelectedDate, 0, len(opts))
	for _, opt := range opts {
		selectedDates = append(selectedDates, SelectedDate{
			Start:    opt.StartTime,
			End:      opt.EndTime,
			Priority: opt.Priority,
		})
	}
	return f.store.createProposedDatesFn(ctx, selectedDates, eventID)
}

type fakeUserCalendarRepository struct {
	reader       *fakeEventReader
	store        *fakeEventTxStore
	calendarRepo *fakeCalendarRepository
}

func (f *fakeUserCalendarRepository) FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*repoUserCalendar.UserCalendar, error) {
	var (
		record *CalendarRecord
		err    error
	)
	if f.reader != nil && f.reader.findAdjustaCandidateCalendarFn != nil {
		record, err = f.reader.findAdjustaCandidateCalendarFn(ctx, userID)
	} else if f.store != nil && f.store.findAdjustaCandidateCalendarFn != nil {
		record, err = f.store.findAdjustaCandidateCalendarFn(ctx, userID)
	} else {
		return nil, errors.New("FilterByUserID should not be called")
	}
	if err != nil {
		return nil, err
	}
	calendar := calendarRecordToDomain(record)
	f.calendarRepo.calendars[calendar.ID] = calendar
	return []*repoUserCalendar.UserCalendar{
		{
			CalendarID:        calendar.ID,
			Role:              domainvalue.UserCalendarRoleAdjustaCandidate,
			SyncProposedDates: record.SyncProposedDates,
		},
	}, nil
}

func (f *fakeUserCalendarRepository) Ensure(ctx context.Context, userID, calendarID uuid.UUID, opt repoUserCalendar.UserCalendarQueryOptions) (*repoUserCalendar.UserCalendar, error) {
	return nil, errors.New("Ensure should not be called")
}

func (f *fakeUserCalendarRepository) SoftDeleteByUserAndCalendar(ctx context.Context, userID, calendarID uuid.UUID) error {
	return errors.New("SoftDeleteByUserAndCalendar should not be called")
}

func calendarRecordToDomain(record *CalendarRecord) *repoCalendar.Calendar {
	if record == nil {
		return nil
	}
	return &repoCalendar.Calendar{
		ID:               record.ID,
		GoogleCalendarID: record.GoogleCalendarID,
		Summary:          record.Summary,
		Description:      record.Description,
		Timezone:         record.Timezone,
	}
}

func applyEventMutation(record *EventRecord, opt EventMutation) {
	if opt.Title != nil {
		record.Title = *opt.Title
	}
	if opt.Location != nil {
		record.Location = *opt.Location
	}
	if opt.Description != nil {
		record.Description = *opt.Description
	}
	if opt.Status != nil {
		record.Status = *opt.Status
	}
	if opt.SyncStatus != nil {
		record.SyncStatus = *opt.SyncStatus
	}
	if opt.ConfirmedDateID != nil {
		record.ConfirmedDateID = *opt.ConfirmedDateID
	}
	if opt.ConfirmedGoogleEventID != nil {
		confirmedGoogleEventID := *opt.ConfirmedGoogleEventID
		record.ConfirmedGoogleEventID = &confirmedGoogleEventID
	}
	if opt.LastSyncedAt != nil {
		record.LastSyncedAt = opt.LastSyncedAt
	}
	if opt.ClearLastSyncedAt {
		record.LastSyncedAt = nil
	}
	if opt.LastSyncError != nil {
		lastSyncError := *opt.LastSyncError
		record.LastSyncError = &lastSyncError
	}
	if opt.ClearLastSyncError {
		record.LastSyncError = nil
	}
}

func applyProposedDateMutation(record *ProposedDateRecord, opt ProposedDateMutation) {
	if opt.GoogleEventID != nil {
		googleEventID := *opt.GoogleEventID
		record.GoogleEventID = &googleEventID
	}
	if opt.StartTime != nil {
		record.StartTime = *opt.StartTime
	}
	if opt.EndTime != nil {
		record.EndTime = *opt.EndTime
	}
	if opt.Priority != nil {
		record.Priority = *opt.Priority
	}
	if opt.Status != nil {
		record.Status = *opt.Status
	}
	if opt.SyncStatus != nil {
		record.SyncStatus = *opt.SyncStatus
	}
	if opt.LastSyncedAt != nil {
		record.LastSyncedAt = opt.LastSyncedAt
	}
	if opt.ClearLastSyncedAt {
		record.LastSyncedAt = nil
	}
	if opt.LastSyncError != nil {
		lastSyncError := *opt.LastSyncError
		record.LastSyncError = &lastSyncError
	}
	if opt.ClearLastSyncError {
		record.LastSyncError = nil
	}
}

func TestFetchAllGoogleEventsReturnsPartialContent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendars := []*CalendarRecord{
		{ID: uuid.New(), GoogleCalendarID: "cal-1", Summary: "Primary"},
	}
	events := []*FetchedGoogleEvent{
		{ID: "event-1", Summary: "Google event"},
	}

	uc := NewUsecase(
		fakeReposFromReader(&fakeEventReader{
			findPrimaryCalendarFn: func(ctx context.Context, userID uuid.UUID) (*CalendarRecord, error) {
				t.Fatalf("find primary calendar should not be called")
				return nil, nil
			},
			listCalendarsByUserFn: func(ctx context.Context, gotUserID uuid.UUID) ([]*CalendarRecord, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				return calendars, nil
			},
			searchEventsFn: func(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*EventRecord, error) {
				t.Fatalf("search events should not be called")
				return nil, nil
			},
			findEventByIDFn: func(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*EventRecord, error) {
				t.Fatalf("find by id should not be called")
				return nil, nil
			},
		}),
		nil,
		&fakeGoogleCalendarGateway{
			fetchEventsFn: func(ctx context.Context, gotUserID uuid.UUID, gotCalendars []*CalendarRecord, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
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
		fakeReposFromReader(&fakeEventReader{
			findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*CalendarRecord, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				return &CalendarRecord{ID: calendarID}, nil
			},
			listCalendarsByUserFn: func(ctx context.Context, userID uuid.UUID) ([]*CalendarRecord, error) {
				t.Fatalf("list calendars should not be called")
				return nil, nil
			},
			searchEventsFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, opt EventSearchOptions) ([]*EventRecord, error) {
				receivedOptions = opt
				return []*EventRecord{
					{
						ID:              eventID2,
						Title:           "Later",
						Status:          domainvalue.StatusConfirmed,
						ConfirmedDateID: dateID2,
						ProposedDates: []*ProposedDateRecord{
							{ID: dateID2, StartTime: now.Add(3 * time.Hour), EndTime: now.Add(4 * time.Hour)},
						},
					},
					{
						ID:              eventID1,
						Title:           "Sooner",
						Status:          domainvalue.StatusConfirmed,
						ConfirmedDateID: dateID1,
						ProposedDates: []*ProposedDateRecord{
							{ID: dateID1, StartTime: now.Add(1 * time.Hour), EndTime: now.Add(2 * time.Hour)},
						},
					},
				}, nil
			},
			findEventByIDFn: func(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*EventRecord, error) {
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
	active := domainvalue.StatusActive

	var receivedOptions EventSearchOptions

	uc := NewUsecase(
		fakeReposFromReader(&fakeEventReader{
			findPrimaryCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*CalendarRecord, error) {
				if gotUserID != userID {
					t.Fatalf("unexpected user id: %s", gotUserID)
				}
				return &CalendarRecord{ID: calendarID}, nil
			},
			listCalendarsByUserFn: func(ctx context.Context, userID uuid.UUID) ([]*CalendarRecord, error) {
				t.Fatalf("list calendars should not be called")
				return nil, nil
			},
			searchEventsFn: func(ctx context.Context, gotUserID, gotCalendarID uuid.UUID, opt EventSearchOptions) ([]*EventRecord, error) {
				receivedOptions = opt
				return []*EventRecord{
					{
						ID:     uuid.New(),
						Title:  "Needs action",
						Status: domainvalue.StatusActive,
						ProposedDates: []*ProposedDateRecord{
							{ID: uuid.New(), StartTime: now.Add(-1 * time.Hour), EndTime: now.Add(1 * time.Hour)},
						},
					},
				}, nil
			},
			findEventByIDFn: func(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*EventRecord, error) {
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
	if drafts[0].Status != domainvalue.StatusActive {
		t.Fatalf("unexpected draft status: %s", drafts[0].Status)
	}
}

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

	storedEvent := &EventRecord{
		ID:          eventID,
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "Discuss roadmap",
		Status:      domainvalue.StatusActive,
		SyncStatus:  domainvalue.SyncStatusPending,
		ProposedDates: []*ProposedDateRecord{
			{
				ID:         dateID1,
				StartTime:  start1,
				EndTime:    end1,
				Priority:   20,
				Status:     domainvalue.ProposedDateStatusActive,
				SyncStatus: domainvalue.SyncStatusPending,
			},
			{
				ID:            dateID2,
				GoogleEventID: &existingGoogleEventID,
				StartTime:     start2,
				EndTime:       end2,
				Priority:      10,
				Status:        domainvalue.ProposedDateStatusActive,
				SyncStatus:    domainvalue.SyncStatusPending,
			},
		},
	}

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*EventRecord, error) {
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
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*CalendarRecord, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					return &CalendarRecord{
						ID:                uuid.New(),
						GoogleCalendarID:  "adjusta-candidate",
						SyncProposedDates: true,
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*ProposedDateRecord, error) {
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
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*EventRecord, error) {
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
			fetchEventsFn: func(ctx context.Context, userID uuid.UUID, calendars []*CalendarRecord, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
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
	if detail.SyncStatus != domainvalue.SyncStatusSynced {
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
		if proposedDate.SyncStatus != domainvalue.SyncStatusSynced {
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

	storedEvent := &EventRecord{
		ID:          eventID,
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "Discuss roadmap",
		Status:      domainvalue.StatusActive,
		SyncStatus:  domainvalue.SyncStatusPending,
		ProposedDates: []*ProposedDateRecord{
			{
				ID:         dateID,
				StartTime:  start,
				EndTime:    end,
				Priority:   10,
				Status:     domainvalue.ProposedDateStatusActive,
				SyncStatus: domainvalue.SyncStatusPending,
			},
		},
	}

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*EventRecord, error) {
					if gotUserID != userID {
						t.Fatalf("unexpected user id: %s", gotUserID)
					}
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return storedEvent, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*CalendarRecord, error) {
					return &CalendarRecord{
						ID:                uuid.New(),
						GoogleCalendarID:  "adjusta-candidate",
						SyncProposedDates: true,
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*ProposedDateRecord, error) {
					if id != dateID {
						t.Fatalf("unexpected proposed date id: %s", id)
					}
					applyProposedDateMutation(storedEvent.ProposedDates[0], opt)
					return storedEvent.ProposedDates[0], nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*EventRecord, error) {
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
			fetchEventsFn: func(ctx context.Context, userID uuid.UUID, calendars []*CalendarRecord, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
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
	if detail.SyncStatus != domainvalue.SyncStatusFailed {
		t.Fatalf("unexpected event sync status: %s", detail.SyncStatus)
	}
	if detail.LastSyncError == nil || *detail.LastSyncError != "google unavailable" {
		t.Fatalf("unexpected event last sync error: %#v", detail.LastSyncError)
	}
	if len(detail.ProposedDates) != 1 {
		t.Fatalf("unexpected proposed dates: %#v", detail.ProposedDates)
	}
	if detail.ProposedDates[0].SyncStatus != domainvalue.SyncStatusFailed {
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

	storedEvent := &EventRecord{
		ID:          eventID,
		Title:       "Draft",
		Location:    "Tokyo",
		Description: "Discuss roadmap",
		Status:      domainvalue.StatusActive,
		SyncStatus:  domainvalue.SyncStatusNotSynced,
		ProposedDates: []*ProposedDateRecord{
			{
				ID:            dateID,
				GoogleEventID: &googleEventID,
				StartTime:     start,
				EndTime:       end,
				Priority:      10,
				Status:        domainvalue.ProposedDateStatusActive,
				SyncStatus:    domainvalue.SyncStatusNotSynced,
			},
		},
	}

	uc := NewUsecase(
		EventTxRepositories{},
		&fakeEventTransaction{
			store: &fakeEventTxStore{
				t: t,
				findEventByIDFn: func(ctx context.Context, gotUserID, gotEventID uuid.UUID, withProposedDates bool) (*EventRecord, error) {
					if gotEventID != eventID {
						t.Fatalf("unexpected event id: %s", gotEventID)
					}
					return storedEvent, nil
				},
				findAdjustaCandidateCalendarFn: func(ctx context.Context, gotUserID uuid.UUID) (*CalendarRecord, error) {
					return &CalendarRecord{
						ID:                uuid.New(),
						GoogleCalendarID:  "adjusta-candidate",
						SyncProposedDates: false,
					}, nil
				},
				updateProposedDateFn: func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*ProposedDateRecord, error) {
					t.Fatal("UpdateProposedDate should not be called when candidate sync is disabled")
					return nil, nil
				},
				updateEventFn: func(ctx context.Context, id uuid.UUID, opt EventMutation) (*EventRecord, error) {
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
			fetchEventsFn: func(ctx context.Context, userID uuid.UUID, calendars []*CalendarRecord, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
				t.Fatalf("fetch events should not be called")
				return nil, nil
			},
		},
	)

	detail, err := uc.FetchDraftedEventDetail(ctx, userID, "user@example.com", eventID)
	if err != nil {
		t.Fatalf("FetchDraftedEventDetail returned error: %v", err)
	}
	if detail.SyncStatus != domainvalue.SyncStatusNotSynced {
		t.Fatalf("unexpected event sync status: %s", detail.SyncStatus)
	}
	if len(detail.ProposedDates) != 1 || detail.ProposedDates[0].SyncStatus != domainvalue.SyncStatusNotSynced {
		t.Fatalf("unexpected proposed dates: %#v", detail.ProposedDates)
	}
}

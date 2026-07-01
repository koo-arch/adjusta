package events

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type fakeEventReader struct {
	findPrimaryCalendarFn          func(ctx context.Context, userID uuid.UUID) (*EventCalendar, error)
	findAdjustaCandidateCalendarFn func(ctx context.Context, userID uuid.UUID) (*EventCalendar, error)
	listCalendarsByUserFn          func(ctx context.Context, userID uuid.UUID) ([]*EventCalendar, error)
	searchEventsFn                 func(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*repoEvent.Event, error)
	findEventByIDFn                func(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*repoEvent.Event, error)
}

func (f *fakeEventReader) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*EventCalendar, error) {
	return f.findPrimaryCalendarFn(ctx, userID)
}

func (f *fakeEventReader) FindAdjustaCandidateCalendar(ctx context.Context, userID uuid.UUID) (*EventCalendar, error) {
	if f.findAdjustaCandidateCalendarFn == nil {
		panic("FindAdjustaCandidateCalendar should not be called")
	}
	return f.findAdjustaCandidateCalendarFn(ctx, userID)
}

func (f *fakeEventReader) ListCalendarsByUser(ctx context.Context, userID uuid.UUID) ([]*EventCalendar, error) {
	return f.listCalendarsByUserFn(ctx, userID)
}

func (f *fakeEventReader) SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*repoEvent.Event, error) {
	return f.searchEventsFn(ctx, userID, calendarID, opt)
}

func (f *fakeEventReader) FindEventByID(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*repoEvent.Event, error) {
	return f.findEventByIDFn(ctx, userID, eventID, withProposedDates)
}

type fakeGoogleCalendarGateway struct {
	fetchEventsFn func(ctx context.Context, userID uuid.UUID, calendars []*EventCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error)
	upsertEventFn func(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
}

func (f *fakeGoogleCalendarGateway) FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*EventCalendar, startTime, endTime time.Time) (*GoogleEventFetchResult, error) {
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
		return eventCalendarToDomain(record), err
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
		calendar := eventCalendarToDomain(record)
		f.calendars[calendar.ID] = calendar
		calendars = append(calendars, calendar)
	}
	return calendars, nil
}

func (f *fakeCalendarRepository) FindByFields(ctx context.Context, userID uuid.UUID, opt repoCalendar.CalendarQueryOptions) (*repoCalendar.Calendar, error) {
	if f.reader != nil && f.reader.findPrimaryCalendarFn != nil {
		record, err := f.reader.findPrimaryCalendarFn(ctx, userID)
		calendar := eventCalendarToDomain(record)
		if calendar != nil {
			f.calendars[calendar.ID] = calendar
		}
		return calendar, err
	}
	if f.store != nil && f.store.findPrimaryCalendarFn != nil {
		record, err := f.store.findPrimaryCalendarFn(ctx, userID)
		calendar := eventCalendarToDomain(record)
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
		record *EventCalendar
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
	calendar := eventCalendarToDomain(record)
	f.calendarRepo.calendars[calendar.ID] = calendar
	return []*repoUserCalendar.UserCalendar{
		{
			CalendarID:        calendar.ID,
			Role:              value.UserCalendarRoleAdjustaCandidate,
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

func eventCalendarToDomain(record *EventCalendar) *repoCalendar.Calendar {
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

func applyEventMutation(record *repoEvent.Event, opt EventMutation) {
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

func applyProposedDateMutation(record *repoProposedDate.ProposedDate, opt ProposedDateMutation) {
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

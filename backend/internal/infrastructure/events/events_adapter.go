package events

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type eventReader struct {
	repos infraRepository.Repositories
}

func NewEventReader(repos infraRepository.Repositories) usecaseEvents.EventReader {
	return &eventReader{repos: repos}
}

func (r *eventReader) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*usecaseEvents.CalendarRecord, error) {
	role := domainvalue.UserCalendarRolePrimary
	calendar, err := r.repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		Role: &role,
	})
	if err != nil {
		return nil, err
	}
	return toCalendarRecord(calendar), nil
}

func (r *eventReader) FindAdjustaCandidateCalendar(ctx context.Context, userID uuid.UUID) (*usecaseEvents.CalendarRecord, error) {
	return findAdjustaCandidateCalendarRecord(ctx, r.repos, userID)
}

func (r *eventReader) ListCalendarsByUser(ctx context.Context, userID uuid.UUID) ([]*usecaseEvents.CalendarRecord, error) {
	calendars, err := r.repos.Calendar.FilterByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return toCalendarRecords(calendars), nil
}

func (r *eventReader) SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt usecaseEvents.EventSearchOptions) ([]*usecaseEvents.EventRecord, error) {
	events, err := r.repos.Event.SearchEvents(ctx, userID, calendarID, toEventQueryOptions(opt))
	if err != nil {
		return nil, err
	}
	return toEventRecords(events), nil
}

func (r *eventReader) FindEventByID(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*usecaseEvents.EventRecord, error) {
	event, err := r.repos.Event.FindByIDAndUser(ctx, userID, eventID, repoEvent.EventQueryOptions{
		WithProposedDates: withProposedDates,
	})
	if err != nil {
		return nil, err
	}
	return toEventRecord(event), nil
}

type eventTransaction struct {
	uow infraRepository.UnitOfWork
}

func NewEventTransaction(uow infraRepository.UnitOfWork) usecaseEvents.EventTransaction {
	return &eventTransaction{
		uow: uow,
	}
}

func (t *eventTransaction) Do(ctx context.Context, fn func(store usecaseEvents.EventTxStore) error) error {
	return t.uow.Do(ctx, func(repos infraRepository.Repositories) error {
		return fn(&eventTxStore{
			repos: repos,
		})
	})
}

type eventTxStore struct {
	repos infraRepository.Repositories
}

func (s *eventTxStore) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*usecaseEvents.CalendarRecord, error) {
	role := domainvalue.UserCalendarRolePrimary
	calendar, err := s.repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		Role: &role,
	})
	if err != nil {
		return nil, err
	}
	return toCalendarRecord(calendar), nil
}

func (s *eventTxStore) FindAdjustaCandidateCalendar(ctx context.Context, userID uuid.UUID) (*usecaseEvents.CalendarRecord, error) {
	return findAdjustaCandidateCalendarRecord(ctx, s.repos, userID)
}

func (s *eventTxStore) FindEventByID(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*usecaseEvents.EventRecord, error) {
	event, err := s.repos.Event.FindByIDAndUser(ctx, userID, eventID, repoEvent.EventQueryOptions{
		WithProposedDates: withProposedDates,
	})
	if err != nil {
		return nil, err
	}
	return toEventRecord(event), nil
}

func (s *eventTxStore) ReadCalendar(ctx context.Context, calendarID uuid.UUID) (*usecaseEvents.CalendarRecord, error) {
	calendar, err := s.repos.Calendar.Read(ctx, calendarID, repoCalendar.CalendarQueryOptions{})
	if err != nil {
		return nil, err
	}
	return toCalendarRecord(calendar), nil
}

func (s *eventTxStore) CreateEvent(ctx context.Context, userID, primaryCalendarID uuid.UUID, title, location, description string, start, end time.Time) (*usecaseEvents.EventRecord, error) {
	event, err := s.repos.Event.Create(ctx, userID, repoEvent.EventCreateOptions{
		Title:       title,
		Location:    location,
		Description: description,
	}, primaryCalendarID)
	if err != nil {
		return nil, err
	}
	return toEventRecord(event), nil
}

func (s *eventTxStore) UpdateEvent(ctx context.Context, id uuid.UUID, opt usecaseEvents.EventMutation) (*usecaseEvents.EventRecord, error) {
	event, err := s.repos.Event.Update(ctx, id, repoEvent.EventQueryOptions{
		Title:                  opt.Title,
		Location:               opt.Location,
		Description:            opt.Description,
		Status:                 opt.Status,
		SyncStatus:             opt.SyncStatus,
		ConfirmedDateID:        opt.ConfirmedDateID,
		ConfirmedGoogleEventID: opt.ConfirmedGoogleEventID,
		LastSyncedAt:           opt.LastSyncedAt,
		ClearLastSyncedAt:      opt.ClearLastSyncedAt,
		LastSyncError:          opt.LastSyncError,
		ClearLastSyncError:     opt.ClearLastSyncError,
	})
	if err != nil {
		return nil, err
	}
	return toEventRecord(event), nil
}

func (s *eventTxStore) SoftDeleteEvent(ctx context.Context, id uuid.UUID) error {
	return s.repos.Event.SoftDelete(ctx, id)
}

func (s *eventTxStore) ListProposedDatesByEvent(ctx context.Context, eventID uuid.UUID) ([]*usecaseEvents.ProposedDateRecord, error) {
	dates, err := s.repos.ProposedDate.FilterByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return toProposedDateRecords(dates), nil
}

func (s *eventTxStore) CreateProposedDates(ctx context.Context, selectedDates []usecaseEvents.SelectedDate, eventID uuid.UUID) ([]*usecaseEvents.ProposedDateRecord, error) {
	dates, err := s.repos.ProposedDate.CreateBulk(ctx, toProposedDateCreateOptionsList(selectedDates), eventID)
	if err != nil {
		return nil, err
	}
	return toProposedDateRecords(dates), nil
}

func (s *eventTxStore) UpdateProposedDate(ctx context.Context, id uuid.UUID, opt usecaseEvents.ProposedDateMutation) (*usecaseEvents.ProposedDateRecord, error) {
	date, err := s.repos.ProposedDate.Update(ctx, id, toProposedDateUpdateOptions(opt))
	if err != nil {
		return nil, err
	}
	return toProposedDateRecord(date), nil
}

func (s *eventTxStore) DeleteProposedDate(ctx context.Context, id uuid.UUID) error {
	return s.repos.ProposedDate.SoftDelete(ctx, id)
}

func (s *eventTxStore) CreateProposedDate(ctx context.Context, opt usecaseEvents.ProposedDateMutation, eventID uuid.UUID) (*usecaseEvents.ProposedDateRecord, error) {
	createOpt, err := toProposedDateCreateOptions(opt)
	if err != nil {
		return nil, err
	}

	date, err := s.repos.ProposedDate.Create(ctx, createOpt, eventID)
	if err != nil {
		return nil, err
	}
	return toProposedDateRecord(date), nil
}

func toEventQueryOptions(opt usecaseEvents.EventSearchOptions) repoEvent.EventQueryOptions {
	return repoEvent.EventQueryOptions{
		Title:                opt.Title,
		Location:             opt.Location,
		Description:          opt.Description,
		Status:               opt.Status,
		WithProposedDates:    opt.WithProposedDates,
		ProposedDateStartGTE: opt.StartTimeGTE,
		ProposedDateStartLTE: opt.StartTimeLTE,
		ProposedDateEndGTE:   opt.EndTimeGTE,
		ProposedDateEndLTE:   opt.EndTimeLTE,
		SortBy:               opt.SortBy,
		SortOrder:            opt.SortOrder,
	}
}

func toProposedDateCreateOptions(opt usecaseEvents.ProposedDateMutation) (repoProposedDate.ProposedDateCreateOptions, error) {
	if opt.Start == nil || opt.End == nil || opt.Priority == nil {
		return repoProposedDate.ProposedDateCreateOptions{}, errors.New("start, end, and priority are required to create proposed date")
	}

	return repoProposedDate.ProposedDateCreateOptions{
		GoogleEventID: opt.GoogleEventID,
		StartTime:     *opt.Start,
		EndTime:       *opt.End,
		Priority:      *opt.Priority,
		Status:        opt.Status,
		SyncStatus:    opt.SyncStatus,
		LastSyncedAt:  opt.LastSyncedAt,
		LastSyncError: opt.LastSyncError,
	}, nil
}

func toProposedDateCreateOptionsList(selectedDates []usecaseEvents.SelectedDate) []repoProposedDate.ProposedDateCreateOptions {
	opts := make([]repoProposedDate.ProposedDateCreateOptions, 0, len(selectedDates))
	for _, selectedDate := range selectedDates {
		status := domainvalue.ProposedDateStatusActive
		opts = append(opts, repoProposedDate.ProposedDateCreateOptions{
			StartTime: selectedDate.Start,
			EndTime:   selectedDate.End,
			Priority:  selectedDate.Priority,
			Status:    &status,
		})
	}
	return opts
}

func toProposedDateUpdateOptions(opt usecaseEvents.ProposedDateMutation) repoProposedDate.ProposedDateUpdateOptions {
	return repoProposedDate.ProposedDateUpdateOptions{
		GoogleEventID:      opt.GoogleEventID,
		StartTime:          opt.Start,
		EndTime:            opt.End,
		Priority:           opt.Priority,
		Status:             opt.Status,
		SyncStatus:         opt.SyncStatus,
		LastSyncedAt:       opt.LastSyncedAt,
		ClearLastSyncedAt:  opt.ClearLastSyncedAt,
		LastSyncError:      opt.LastSyncError,
		ClearLastSyncError: opt.ClearLastSyncError,
	}
}

func toCalendarRecord(calendar *repoCalendar.Calendar) *usecaseEvents.CalendarRecord {
	return toCalendarRecordWithSync(calendar, false)
}

func toCalendarRecordWithSync(calendar *repoCalendar.Calendar, syncProposedDates bool) *usecaseEvents.CalendarRecord {
	if calendar == nil {
		return nil
	}

	return &usecaseEvents.CalendarRecord{
		ID:                calendar.ID,
		GoogleCalendarID:  calendar.GoogleCalendarID,
		Summary:           calendar.Summary,
		Description:       calendar.Description,
		Timezone:          calendar.Timezone,
		SyncProposedDates: syncProposedDates,
	}
}

func toCalendarRecords(calendars []*repoCalendar.Calendar) []*usecaseEvents.CalendarRecord {
	records := make([]*usecaseEvents.CalendarRecord, 0, len(calendars))
	for _, calendar := range calendars {
		records = append(records, toCalendarRecord(calendar))
	}
	return records
}

func findAdjustaCandidateCalendarRecord(ctx context.Context, repos infraRepository.Repositories, userID uuid.UUID) (*usecaseEvents.CalendarRecord, error) {
	userCalendars, err := repos.UserCalendar.FilterByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, userCalendar := range userCalendars {
		if userCalendar.Role != domainvalue.UserCalendarRoleAdjustaCandidate {
			continue
		}

		calendar, err := repos.Calendar.Read(ctx, userCalendar.CalendarID, repoCalendar.CalendarQueryOptions{})
		if err != nil {
			return nil, err
		}

		return toCalendarRecordWithSync(calendar, userCalendar.SyncProposedDates), nil
	}

	return nil, repoerr.ErrNotFound
}

func toEventRecord(event *repoEvent.Event) *usecaseEvents.EventRecord {
	if event == nil {
		return nil
	}

	return &usecaseEvents.EventRecord{
		ID:                     event.ID,
		PrimaryCalendarID:      event.PrimaryCalendarID,
		Title:                  event.Title,
		Location:               event.Location,
		Description:            event.Description,
		Status:                 event.Status,
		ConfirmedDateID:        event.ConfirmedDateID,
		ConfirmedGoogleEventID: event.ConfirmedGoogleEventID,
		SyncStatus:             event.SyncStatus,
		LastSyncedAt:           event.LastSyncedAt,
		LastSyncError:          event.LastSyncError,
		ProposedDates:          toProposedDateRecords(event.ProposedDates),
	}
}

func toEventRecords(events []*repoEvent.Event) []*usecaseEvents.EventRecord {
	records := make([]*usecaseEvents.EventRecord, 0, len(events))
	for _, event := range events {
		records = append(records, toEventRecord(event))
	}
	return records
}

func toProposedDateRecord(date *repoProposedDate.ProposedDate) *usecaseEvents.ProposedDateRecord {
	if date == nil {
		return nil
	}

	return &usecaseEvents.ProposedDateRecord{
		ID:            date.ID,
		EventID:       date.EventID,
		GoogleEventID: date.GoogleEventID,
		StartTime:     date.StartTime,
		EndTime:       date.EndTime,
		Priority:      date.Priority,
		Status:        date.Status,
		SyncStatus:    date.SyncStatus,
		LastSyncedAt:  date.LastSyncedAt,
		LastSyncError: date.LastSyncError,
	}
}

func toProposedDateRecords(dates []*repoProposedDate.ProposedDate) []*usecaseEvents.ProposedDateRecord {
	records := make([]*usecaseEvents.ProposedDateRecord, 0, len(dates))
	for _, date := range dates {
		records = append(records, toProposedDateRecord(date))
	}
	return records
}

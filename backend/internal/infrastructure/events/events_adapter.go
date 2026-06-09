package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	infraGoogleCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/googlecalendar"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type eventReader struct {
	repos infraRepository.Repositories
}

func NewEventReader(repos infraRepository.Repositories) usecaseEvents.EventReader {
	return &eventReader{repos: repos}
}

func (r *eventReader) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
	role := domainvalue.UserCalendarRolePrimary
	return r.repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		Role: &role,
	})
}

func (r *eventReader) ListCalendarsByUser(ctx context.Context, userID uuid.UUID) ([]*repositorymodel.StoredCalendar, error) {
	return r.repos.Calendar.FilterByUserID(ctx, userID)
}

func (r *eventReader) SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt usecaseEvents.EventSearchOptions) ([]*repositorymodel.StoredEvent, error) {
	return r.repos.Event.SearchEvents(ctx, userID, calendarID, toEventQueryOptions(opt))
}

func (r *eventReader) FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
	return r.repos.Event.FindBySlugAndUser(ctx, userID, slug, repoEvent.EventQueryOptions{
		WithProposedDates: withProposedDates,
	})
}

type eventTransaction struct {
	uow         infraRepository.UnitOfWork
	calendarApp *infraGoogleCalendar.GoogleCalendarManager
}

func NewEventTransaction(uow infraRepository.UnitOfWork, calendarApp *infraGoogleCalendar.GoogleCalendarManager) usecaseEvents.EventTransaction {
	return &eventTransaction{
		uow:         uow,
		calendarApp: calendarApp,
	}
}

func (t *eventTransaction) Do(ctx context.Context, fn func(store usecaseEvents.EventTxStore) error) error {
	return t.uow.Do(ctx, func(repos infraRepository.Repositories) error {
		return fn(&eventTxStore{
			repos:       repos,
			calendarApp: t.calendarApp,
		})
	})
}

type eventTxStore struct {
	repos       infraRepository.Repositories
	calendarApp *infraGoogleCalendar.GoogleCalendarManager
}

func (s *eventTxStore) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
	role := domainvalue.UserCalendarRolePrimary
	return s.repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		Role: &role,
	})
}

func (s *eventTxStore) FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*repositorymodel.StoredEvent, error) {
	return s.repos.Event.FindBySlugAndUser(ctx, userID, slug, repoEvent.EventQueryOptions{
		WithProposedDates: withProposedDates,
	})
}

func (s *eventTxStore) ReadCalendar(ctx context.Context, calendarID uuid.UUID) (*repositorymodel.StoredCalendar, error) {
	return s.repos.Calendar.Read(ctx, calendarID, repoCalendar.CalendarQueryOptions{})
}

func (s *eventTxStore) CreateEvent(ctx context.Context, userID, primaryCalendarID uuid.UUID, title, location, description string, start, end time.Time) (*repositorymodel.StoredEvent, error) {
	googleEvent := s.calendarApp.ConvertToCalendarEvent(nil, title, location, description, start, end)
	return s.repos.Event.Create(ctx, userID, googleEvent, primaryCalendarID)
}

func (s *eventTxStore) UpdateEvent(ctx context.Context, id uuid.UUID, opt usecaseEvents.EventMutation) (*repositorymodel.StoredEvent, error) {
	return s.repos.Event.Update(ctx, id, repoEvent.EventQueryOptions{
		Title:                  opt.Title,
		Location:               opt.Location,
		Description:            opt.Description,
		Status:                 opt.Status,
		SyncStatus:             opt.SyncStatus,
		ConfirmedDateID:        opt.ConfirmedDateID,
		GoogleEventID:          opt.GoogleEventID,
		ConfirmedGoogleEventID: opt.ConfirmedGoogleEventID,
		LastSyncedAt:           opt.LastSyncedAt,
		ClearLastSyncedAt:      opt.ClearLastSyncedAt,
		LastSyncError:          opt.LastSyncError,
		ClearLastSyncError:     opt.ClearLastSyncError,
	})
}

func (s *eventTxStore) SoftDeleteEvent(ctx context.Context, id uuid.UUID) error {
	return s.repos.Event.SoftDelete(ctx, id)
}

func (s *eventTxStore) ListProposedDatesByEvent(ctx context.Context, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
	return s.repos.ProposedDate.FilterByEventID(ctx, eventID)
}

func (s *eventTxStore) CreateProposedDates(ctx context.Context, selectedDates []appmodel.SelectedDate, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
	return s.repos.ProposedDate.CreateBulk(ctx, selectedDates, eventID)
}

func (s *eventTxStore) UpdateProposedDate(ctx context.Context, id uuid.UUID, opt usecaseEvents.ProposedDateMutation) (*repositorymodel.StoredProposedDate, error) {
	return s.repos.ProposedDate.Update(ctx, id, toProposedDateQueryOptions(opt))
}

func (s *eventTxStore) DeleteProposedDate(ctx context.Context, id uuid.UUID) error {
	return s.repos.ProposedDate.SoftDelete(ctx, id)
}

func (s *eventTxStore) CreateProposedDate(ctx context.Context, opt usecaseEvents.ProposedDateMutation, eventID uuid.UUID) (*repositorymodel.StoredProposedDate, error) {
	return s.repos.ProposedDate.Create(ctx, toProposedDateQueryOptions(opt), eventID)
}

func (s *eventTxStore) DecrementPriorityExceptID(ctx context.Context, eventID, excludeID uuid.UUID) error {
	return s.repos.ProposedDate.DecrementPriorityExceptID(ctx, eventID, excludeID)
}

func (s *eventTxStore) ReorderPriority(ctx context.Context, eventID uuid.UUID) error {
	return s.repos.ProposedDate.ReorderPriority(ctx, eventID)
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

func toProposedDateQueryOptions(opt usecaseEvents.ProposedDateMutation) repoProposedDate.ProposedDateQueryOptions {
	return repoProposedDate.ProposedDateQueryOptions{
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

package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/repo/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type eventReader struct {
	repos Repositories
}

func NewEventReader(repos Repositories) usecaseEvents.EventReader {
	return &eventReader{repos: repos}
}

func (r *eventReader) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*models.StoredCalendar, error) {
	isPrimary := true
	return r.repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	})
}

func (r *eventReader) ListGoogleCalendarInfosByUser(ctx context.Context, userID uuid.UUID) ([]*models.GoogleCalendarInfo, error) {
	return r.repos.GoogleCalendarInfo.ListByUser(ctx, userID)
}

func (r *eventReader) SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt usecaseEvents.EventSearchOptions) ([]*models.StoredEvent, error) {
	return r.repos.Event.SearchEvents(ctx, userID, calendarID, toEventQueryOptions(opt))
}

func (r *eventReader) FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*models.StoredEvent, error) {
	return r.repos.Event.FindBySlugAndUser(ctx, userID, slug, repoEvent.EventQueryOptions{
		WithProposedDates: withProposedDates,
	})
}

type eventTransaction struct {
	uow         UnitOfWork
	calendarApp *appCalendar.GoogleCalendarManager
}

func NewEventTransaction(uow UnitOfWork, calendarApp *appCalendar.GoogleCalendarManager) usecaseEvents.EventTransaction {
	return &eventTransaction{
		uow:         uow,
		calendarApp: calendarApp,
	}
}

func (t *eventTransaction) Do(ctx context.Context, fn func(store usecaseEvents.EventTxStore) error) error {
	return t.uow.Do(ctx, func(repos Repositories) error {
		return fn(&eventTxStore{
			repos:       repos,
			calendarApp: t.calendarApp,
		})
	})
}

type eventTxStore struct {
	repos       Repositories
	calendarApp *appCalendar.GoogleCalendarManager
}

func (s *eventTxStore) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*models.StoredCalendar, error) {
	isPrimary := true
	return s.repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	})
}

func (s *eventTxStore) FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*models.StoredEvent, error) {
	return s.repos.Event.FindBySlugAndUser(ctx, userID, slug, repoEvent.EventQueryOptions{
		WithProposedDates: withProposedDates,
	})
}

func (s *eventTxStore) CreateEvent(ctx context.Context, calendarID uuid.UUID, title, location, description string, start, end time.Time) (*models.StoredEvent, error) {
	googleEvent := s.calendarApp.ConvertToCalendarEvent(nil, title, location, description, start, end)
	return s.repos.Event.Create(ctx, googleEvent, calendarID)
}

func (s *eventTxStore) UpdateEvent(ctx context.Context, id uuid.UUID, opt usecaseEvents.EventMutation) (*models.StoredEvent, error) {
	return s.repos.Event.Update(ctx, id, repoEvent.EventQueryOptions{
		Summary:         opt.Title,
		Location:        opt.Location,
		Description:     opt.Description,
		Status:          opt.Status,
		ConfirmedDateID: opt.ConfirmedDateID,
		GoogleEventID:   opt.GoogleEventID,
	})
}

func (s *eventTxStore) SoftDeleteEvent(ctx context.Context, id uuid.UUID) error {
	return s.repos.Event.SoftDelete(ctx, id)
}

func (s *eventTxStore) ListProposedDatesByEvent(ctx context.Context, eventID uuid.UUID) ([]*models.StoredProposedDate, error) {
	return s.repos.ProposedDate.FilterByEventID(ctx, eventID)
}

func (s *eventTxStore) CreateProposedDates(ctx context.Context, selectedDates []models.SelectedDate, eventID uuid.UUID) ([]*models.StoredProposedDate, error) {
	return s.repos.ProposedDate.CreateBulk(ctx, selectedDates, eventID)
}

func (s *eventTxStore) UpdateProposedDate(ctx context.Context, id uuid.UUID, opt usecaseEvents.ProposedDateMutation) (*models.StoredProposedDate, error) {
	return s.repos.ProposedDate.Update(ctx, id, toProposedDateQueryOptions(opt))
}

func (s *eventTxStore) DeleteProposedDate(ctx context.Context, id uuid.UUID) error {
	return s.repos.ProposedDate.Delete(ctx, id)
}

func (s *eventTxStore) CreateProposedDate(ctx context.Context, opt usecaseEvents.ProposedDateMutation, eventID uuid.UUID) (*models.StoredProposedDate, error) {
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
		Summary:              opt.Title,
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
		StartTime: opt.Start,
		EndTime:   opt.End,
		Priority:  opt.Priority,
	}
}

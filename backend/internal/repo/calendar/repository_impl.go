package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/ent/user"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repo/infraerr"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type CalendarRepositoryImpl struct {
	client *ent.Client
}

func NewCalendarRepository(client *ent.Client) *CalendarRepositoryImpl {
	return &CalendarRepositoryImpl{
		client: client,
	}
}

func (r *CalendarRepositoryImpl) WithTx(tx transaction.Tx) CalendarRepository {
	return &CalendarRepositoryImpl{client: tx.Client()}
}

func (r *CalendarRepositoryImpl) Read(ctx context.Context, id uuid.UUID, opt CalendarQueryOptions) (*models.StoredCalendar, error) {
	findCalendar := r.client.Calendar.Query()
	entity, err := findCalendar.Where(calendar.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toStoredCalendar(entity), nil
}

func (r *CalendarRepositoryImpl) FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*models.StoredCalendar, error) {
	entities, err := r.client.Calendar.Query().
		Where(calendar.HasUserWith(user.ID(userID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredCalendars(entities), nil
}

func (r *CalendarRepositoryImpl) FindByFields(ctx context.Context, userID uuid.UUID, opt CalendarQueryOptions) (*models.StoredCalendar, error) {
	if !opt.WithEvents && opt.WithProposedDates {
		return nil, fmt.Errorf("WithDates is only available when withEvents is true")
	}

	findCalendar := r.client.Calendar.Query()
	query := r.applyCalendarQueryOptions(findCalendar, userID, opt)
	entity, err := query.Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toStoredCalendar(entity), nil
}

func (r *CalendarRepositoryImpl) FilterByFields(ctx context.Context, userID uuid.UUID, opt CalendarQueryOptions) ([]*models.StoredCalendar, error) {
	filterCalendar := r.client.Calendar.Query()

	if !opt.WithEvents && opt.WithProposedDates {
		return nil, fmt.Errorf("WithDates is only available when withEvents is true")
	}
	query := r.applyCalendarQueryOptions(filterCalendar, userID, opt)
	entities, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredCalendars(entities), nil
}

func (r *CalendarRepositoryImpl) Create(ctx context.Context, userID uuid.UUID) (*models.StoredCalendar, error) {
	createCalendar := r.client.Calendar.Create()
	entity, err := createCalendar.
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredCalendar(entity), nil
}

func (r *CalendarRepositoryImpl) Update(ctx context.Context, id uuid.UUID) (*models.StoredCalendar, error) {
	entity, err := r.client.Calendar.UpdateOneID(id).Save(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toStoredCalendar(entity), nil
}

func (r *CalendarRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.Calendar.DeleteOneID(id).Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *CalendarRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.client.Calendar.UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *CalendarRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	err := r.client.Calendar.UpdateOneID(id).
		SetNillableDeletedAt(nil).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *CalendarRepositoryImpl) applyCalendarQueryOptions(query *ent.CalendarQuery, userID uuid.UUID, opt CalendarQueryOptions) *ent.CalendarQuery {
	query = query.Where(calendar.HasUserWith(user.IDEQ(userID)))

	if opt.GoogleCalendarID != nil {
		query = query.Where(calendar.HasGoogleCalendarInfosWith(googlecalendarinfo.GoogleCalendarIDEQ(*opt.GoogleCalendarID)))
	}
	if opt.Summary != nil {
		query = query.Where(calendar.HasGoogleCalendarInfosWith(googlecalendarinfo.SummaryEQ(*opt.Summary)))
	}
	if opt.IsPrimary != nil {
		query = query.Where(calendar.HasGoogleCalendarInfosWith(googlecalendarinfo.IsPrimaryEQ(*opt.IsPrimary)))
	}

	return query
}

func toStoredCalendar(entity *ent.Calendar) *models.StoredCalendar {
	if entity == nil {
		return nil
	}

	return &models.StoredCalendar{
		ID: entity.ID,
	}
}

func toStoredCalendars(entities []*ent.Calendar) []*models.StoredCalendar {
	storedCalendars := make([]*models.StoredCalendar, 0, len(entities))
	for _, entity := range entities {
		storedCalendars = append(storedCalendars, toStoredCalendar(entity))
	}
	return storedCalendars
}

package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	dbUserCalendar "github.com/koo-arch/adjusta-backend/ent/usercalendar"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
)

type CalendarRepository = repoCalendar.CalendarRepository
type CalendarQueryOptions = repoCalendar.CalendarQueryOptions
type CalendarMutationOptions = repoCalendar.CalendarMutationOptions

type CalendarRepositoryImpl struct {
	client *ent.Client
}

func NewCalendarRepository(client *ent.Client) *CalendarRepositoryImpl {
	return &CalendarRepositoryImpl{
		client: client,
	}
}

func (r *CalendarRepositoryImpl) Read(ctx context.Context, id uuid.UUID) (*repoCalendar.Calendar, error) {
	findCalendar := r.client.Calendar.Query()
	entity, err := findCalendar.Where(calendar.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toCalendar(entity), nil
}

func (r *CalendarRepositoryImpl) FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*repoCalendar.Calendar, error) {
	entities, err := r.client.Calendar.Query().
		Where(calendar.HasUserCalendarsWith(dbUserCalendar.UserIDEQ(userID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toCalendars(entities), nil
}

func (r *CalendarRepositoryImpl) FindByFields(ctx context.Context, userID uuid.UUID, opt CalendarQueryOptions) (*repoCalendar.Calendar, error) {
	if !opt.WithEvents && opt.WithProposedDates {
		return nil, fmt.Errorf("WithDates is only available when withEvents is true")
	}

	findCalendar := r.client.Calendar.Query()
	query := r.applyCalendarQueryOptions(findCalendar, userID, opt)
	entity, err := query.Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toCalendar(entity), nil
}

func (r *CalendarRepositoryImpl) FindByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
	entity, err := r.client.Calendar.Query().
		Where(calendar.GoogleCalendarIDEQ(googleCalendarID)).
		Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toCalendar(entity), nil
}

func (r *CalendarRepositoryImpl) FilterByFields(ctx context.Context, userID uuid.UUID, opt CalendarQueryOptions) ([]*repoCalendar.Calendar, error) {
	filterCalendar := r.client.Calendar.Query()

	if !opt.WithEvents && opt.WithProposedDates {
		return nil, fmt.Errorf("WithDates is only available when withEvents is true")
	}
	query := r.applyCalendarQueryOptions(filterCalendar, userID, opt)
	entities, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return toCalendars(entities), nil
}

func (r *CalendarRepositoryImpl) Create(ctx context.Context, opt CalendarMutationOptions) (*repoCalendar.Calendar, error) {
	createCalendar := r.client.Calendar.Create()
	applyCalendarCreateOptions(createCalendar, opt)
	entity, err := createCalendar.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toCalendar(entity), nil
}

func (r *CalendarRepositoryImpl) Update(ctx context.Context, id uuid.UUID, opt CalendarMutationOptions) (*repoCalendar.Calendar, error) {
	updateCalendar := r.client.Calendar.UpdateOneID(id)
	applyCalendarUpdateOptions(updateCalendar, opt)
	entity, err := updateCalendar.Save(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toCalendar(entity), nil
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
	query = query.Where(calendar.HasUserCalendarsWith(dbUserCalendar.UserIDEQ(userID)))

	if opt.GoogleCalendarID != nil {
		query = query.Where(calendar.GoogleCalendarIDEQ(*opt.GoogleCalendarID))
	}
	if opt.Summary != nil {
		query = query.Where(calendar.SummaryEQ(*opt.Summary))
	}
	if opt.Role != nil {
		query = query.Where(calendar.HasUserCalendarsWith(
			dbUserCalendar.UserIDEQ(userID),
			dbUserCalendar.RoleEQ(dbUserCalendar.Role(*opt.Role)),
		))
	}

	return query
}

func applyCalendarCreateOptions(create *ent.CalendarCreate, opt CalendarMutationOptions) {
	create.SetNillableGoogleCalendarID(opt.GoogleCalendarID)
	create.SetNillableSummary(opt.Summary)
	create.SetNillableDescription(opt.Description)
	create.SetNillableTimezone(opt.Timezone)
}

func applyCalendarUpdateOptions(update *ent.CalendarUpdateOne, opt CalendarMutationOptions) {
	update.SetNillableGoogleCalendarID(opt.GoogleCalendarID)
	update.SetNillableSummary(opt.Summary)
	update.SetNillableDescription(opt.Description)
	update.SetNillableTimezone(opt.Timezone)
}

func toCalendar(entity *ent.Calendar) *repoCalendar.Calendar {
	if entity == nil {
		return nil
	}

	return &repoCalendar.Calendar{
		ID:               entity.ID,
		GoogleCalendarID: valueOrEmpty(entity.GoogleCalendarID),
		Summary:          valueOrEmpty(entity.Summary),
		Description:      entity.Description,
		Timezone:         entity.Timezone,
	}
}

func toCalendars(entities []*ent.Calendar) []*repoCalendar.Calendar {
	calendars := make([]*repoCalendar.Calendar, 0, len(entities))
	for _, entity := range entities {
		calendars = append(calendars, toCalendar(entity))
	}
	return calendars
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

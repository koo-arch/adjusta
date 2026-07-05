package usercalendar

import (
	"context"
	"time"

	"github.com/google/uuid"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/mixins"
	dbUserCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/usercalendar"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
)

type UserCalendarRepository = repoUserCalendar.UserCalendarRepository
type UserCalendarQueryOptions = repoUserCalendar.UserCalendarQueryOptions

type UserCalendarRepositoryImpl struct {
	client *ent.Client
}

func NewUserCalendarRepository(client *ent.Client) *UserCalendarRepositoryImpl {
	return &UserCalendarRepositoryImpl{client: client}
}

func (r *UserCalendarRepositoryImpl) FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*repoUserCalendar.UserCalendar, error) {
	entities, err := r.client.UserCalendar.Query().
		Where(dbUserCalendar.UserIDEQ(userID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	userCalendars := make([]*repoUserCalendar.UserCalendar, 0, len(entities))
	for _, entity := range entities {
		userCalendars = append(userCalendars, toModelUserCalendar(entity))
	}
	return userCalendars, nil
}

func (r *UserCalendarRepositoryImpl) Ensure(ctx context.Context, userID, calendarID uuid.UUID, opt UserCalendarQueryOptions) (*repoUserCalendar.UserCalendar, error) {
	entity, err := r.client.UserCalendar.Query().
		Where(
			dbUserCalendar.UserIDEQ(userID),
			dbUserCalendar.CalendarIDEQ(calendarID),
		).
		Only(mixins.SkipSoftDelete(ctx))
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}

		create := r.client.UserCalendar.Create().
			SetUserID(userID).
			SetCalendarID(calendarID)
		applyUserCalendarCreateOptions(create, opt)
		entity, err = create.Save(ctx)
		if err != nil {
			return nil, err
		}
		return toModelUserCalendar(entity), nil
	}

	update := r.client.UserCalendar.UpdateOneID(entity.ID)
	if entity.DeletedAt != nil {
		update = update.SetNillableDeletedAt(nil)
	}
	applyUserCalendarUpdateOptions(update, opt)
	entity, err = update.Save(mixins.SkipSoftDelete(ctx))
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelUserCalendar(entity), nil
}

func (r *UserCalendarRepositoryImpl) SoftDeleteByUserAndCalendar(ctx context.Context, userID, calendarID uuid.UUID) error {
	_, err := r.client.UserCalendar.Query().
		Where(
			dbUserCalendar.UserIDEQ(userID),
			dbUserCalendar.CalendarIDEQ(calendarID),
		).
		Only(ctx)
	if err != nil {
		return infraerr.MapNotFound(err)
	}

	err = r.client.UserCalendar.Update().
		Where(
			dbUserCalendar.UserIDEQ(userID),
			dbUserCalendar.CalendarIDEQ(calendarID),
		).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func applyUserCalendarCreateOptions(create *ent.UserCalendarCreate, opt UserCalendarQueryOptions) {
	if opt.Role != nil {
		create.SetRole(dbUserCalendar.Role(*opt.Role))
	}
	if opt.IsVisible != nil {
		create.SetIsVisible(*opt.IsVisible)
	}
	if opt.SyncProposedDates != nil {
		create.SetSyncProposedDates(*opt.SyncProposedDates)
	}
}

func applyUserCalendarUpdateOptions(update *ent.UserCalendarUpdateOne, opt UserCalendarQueryOptions) {
	if opt.Role != nil {
		update.SetRole(dbUserCalendar.Role(*opt.Role))
	}
	if opt.IsVisible != nil {
		update.SetIsVisible(*opt.IsVisible)
	}
	if opt.SyncProposedDates != nil {
		update.SetSyncProposedDates(*opt.SyncProposedDates)
	}
}

func toModelUserCalendar(entity *ent.UserCalendar) *repoUserCalendar.UserCalendar {
	if entity == nil {
		return nil
	}

	return &repoUserCalendar.UserCalendar{
		ID:                entity.ID,
		UserID:            entity.UserID,
		CalendarID:        entity.CalendarID,
		Role:              value.UserCalendarRole(entity.Role),
		IsVisible:         entity.IsVisible,
		SyncProposedDates: entity.SyncProposedDates,
	}
}

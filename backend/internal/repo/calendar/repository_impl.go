package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/user"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
	"github.com/koo-arch/adjusta-backend/ent/googlecalendarinfo"
)

type CalendarRepositoryImpl struct {
	client *ent.Client
}

func NewCalendarRepository(client *ent.Client) *CalendarRepositoryImpl {
	return &CalendarRepositoryImpl{
		client: client,
	}
}

func (r *CalendarRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt CalendarQueryOptions) (*ent.Calendar, error) {
	findCalendar := r.client.Calendar.Query()
	if tx != nil {
		findCalendar = tx.Calendar.Query()
	}
	
	if opt.WithGoogleCalendarInfo {
		findCalendar = findCalendar.WithGoogleCalendarInfos()
	}

	return findCalendar.Only(ctx)
}

func (r *CalendarRepositoryImpl) FilterByUserID(ctx context.Context, tx *ent.Tx, userID uuid.UUID) ([]*ent.Calendar, error) {
	fileterCalendar := r.client.Calendar.Query()
	if tx != nil {
		fileterCalendar = tx.Calendar.Query()
	}
	return fileterCalendar.
		Where(calendar.HasUserWith(user.ID(userID))).
		All(ctx)
}

func (r *CalendarRepositoryImpl) FindByFields(ctx context.Context, tx *ent.Tx, userID uuid.UUID, opt CalendarQueryOptions) (*ent.Calendar, error) {
	if !opt.WithEvents && opt.WithProposedDates {
		return nil, fmt.Errorf("WithDates is only available when withEvents is true")
	}

	findCalendar := r.client.Calendar.Query()
	if tx != nil {
		findCalendar = tx.Calendar.Query()
	}

	if opt.WithGoogleCalendarInfo {
		findCalendar = findCalendar.WithGoogleCalendarInfos()
	}

	query := r.applyCalendarQueryOptions(findCalendar, userID, opt)

	return query.Only(ctx)
}

func (r *CalendarRepositoryImpl) FilterByFields(ctx context.Context, tx *ent.Tx, userID uuid.UUID, opt CalendarQueryOptions) ([]*ent.Calendar, error) {
	filterCalendar := r.client.Calendar.Query()
	
	if !opt.WithEvents && opt.WithProposedDates {
		return nil, fmt.Errorf("WithDates is only available when withEvents is true")
	}
	if tx != nil {
		filterCalendar = tx.Calendar.Query()
	}

	if opt.WithGoogleCalendarInfo {
		filterCalendar = filterCalendar.WithGoogleCalendarInfos()
	}
	
	query := r.applyCalendarQueryOptions(filterCalendar, userID, opt)

	return query.All(ctx)
}

func (r *CalendarRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, entUser *ent.User, entGoogleCalendar *ent.GoogleCalendarInfo) (*ent.Calendar, error) {
	createCalendar := r.client.Calendar.Create()
	if tx != nil {
		createCalendar = tx.Calendar.Create()
	}

	if entGoogleCalendar != nil {
		createCalendar = createCalendar.AddGoogleCalendarInfos(entGoogleCalendar)
	}

	return createCalendar.
		SetUser(entUser).
		Save(ctx)
}

func (r *CalendarRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.Calendar, error) {
	updateCalendar := r.client.Calendar.UpdateOneID(id)
	if tx != nil {
		updateCalendar = tx.Calendar.UpdateOneID(id)
	}
	return updateCalendar.
		Save(ctx)
}

func (r *CalendarRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.Calendar.DeleteOneID(id).Exec(ctx)
	}
	return r.client.Calendar.DeleteOneID(id).Exec(ctx)
}

func (r *CalendarRepositoryImpl) SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	softDeleteCalendar := r.client.Calendar.UpdateOneID(id)
	if tx != nil {
		softDeleteCalendar = tx.Calendar.UpdateOneID(id)
	}
	return softDeleteCalendar.
		SetDeletedAt(time.Now()).
		Exec(ctx)
}

func (r *CalendarRepositoryImpl) Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	restoreCalendar := r.client.Calendar.UpdateOneID(id)
	if tx != nil {
		restoreCalendar = tx.Calendar.UpdateOneID(id)
	}
	return restoreCalendar.
		SetNillableDeletedAt(nil).
		Exec(ctx)
}

func (r *CalendarRepositoryImpl) SoftDeleteWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	// カレンダーを論理削除
	if err := r.SoftDelete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to soft delete calendar: %w", err)
	}
	// 関連するイベントIDを取得
	eventIDs, err := tx.Event.Query().Where(event.HasCalendarWith(calendar.ID(id))).Select(event.FieldID).IDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to query events: %w", err)
	}


	if len(eventIDs) > 0 {
		// 関連するイベントを論理削除
		softDeleteEvent := tx.Event.Update().Where(event.IDIn(eventIDs...))
		if err := softDeleteEvent.SetDeletedAt(time.Now()).Exec(ctx); err != nil {
			return fmt.Errorf("failed to soft delete events: %w", err)
		}

		// 関連する候補日程を論理削除
		softDeleteProposedDate := tx.ProposedDate.Update().Where(proposeddate.HasEventWith(event.IDIn(eventIDs...)))
		if err := softDeleteProposedDate.SetDeletedAt(time.Now()).Exec(ctx); err != nil {
			return fmt.Errorf("failed to soft delete proposed dates: %w", err)
		}
	}

	return nil
}

func (r *CalendarRepositoryImpl) RestoreWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}

	// カレンダーを復元
	if err := r.Restore(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to restore calendar: %w", err)
	}

	// 関連するイベントを復元
	eventIDs, err := tx.Event.
		Query().
		Where(event.HasCalendarWith(calendar.ID(id))).
		IDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to query events: %w", err)
	}
	if len(eventIDs) > 0{
		if err := tx.Event.
			Update().
			Where(event.IDIn(eventIDs...)).
			SetNillableDeletedAt(nil).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to restore events: %w", err)
		}

		// 関連する候補日程を復元
		proposedDateIDs, err := tx.ProposedDate.
			Query().
			Where(proposeddate.HasEventWith(event.IDIn(eventIDs...))).
			IDs(ctx)
		if err != nil {
			return fmt.Errorf("failed to query proposed dates: %w", err)
		}
		if len(proposedDateIDs) > 0 {
			if err := tx.ProposedDate.
				Update().
				Where(proposeddate.IDIn(proposedDateIDs...)).
				SetNillableDeletedAt(nil).
				Exec(ctx); err != nil {
				return fmt.Errorf("failed to restore proposed dates: %w", err)
			}
		}
	}

	return nil
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

	if opt.WithEvents {
		query = query.WithEvents(func(eventQuery *ent.EventQuery) {
			if opt.EventOffset > 0 {
				eventQuery = eventQuery.Offset(opt.EventOffset)
			}
			if opt.EventLimit > 0 {
				eventQuery = eventQuery.Limit(opt.EventLimit)
			}

			if opt.WithProposedDates {
				eventQuery = eventQuery.WithProposedDates(func(dateQuery *ent.ProposedDateQuery) {
					if opt.ProposedDateOffset > 0 {
						dateQuery = dateQuery.Offset(opt.ProposedDateOffset)
					}
					if opt.ProposedDateLimit > 0 {
						dateQuery = dateQuery.Limit(opt.ProposedDateLimit)
					}
				})
			}
		})
	}

	return query
}
package event

import (
	"context"
	"time"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/user"
	dbCalendar "github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
)

type EventRepositoryImpl struct {
	client *ent.Client
}

func NewEventRepository(client *ent.Client) *EventRepositoryImpl {
	return &EventRepositoryImpl{
		client: client,
	}
}

func (r *EventRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt EventQueryOptions) (*ent.Event, error) {
	query := r.client.Event.Query()
	if tx != nil {
		query = tx.Event.Query()
	}

	if opt.WithProposedDates {
		query = query.WithProposedDates()
	}

	return query.Where(event.IDEQ(id)).Only(ctx)
}

func (r *EventRepositoryImpl) FilterByCalendarID(ctx context.Context, tx *ent.Tx, calendarID uuid.UUID, opt EventQueryOptions) ([]*ent.Event, error) {
	filterEvent := r.client.Event.Query()
	if tx != nil {
		filterEvent = tx.Event.Query()
	}

	filterEvent = filterEvent.Where(event.HasCalendarWith(dbCalendar.IDEQ(calendarID)))

	// イベントに対するオフセットとリミットを適用
	if opt.EventOffset > 0 {
		filterEvent = filterEvent.Offset(opt.EventOffset)
	}
	if opt.EventLimit > 0 {
		filterEvent = filterEvent.Limit(opt.EventLimit)
	}

	// イベントの提案日に対するオフセットとリミットを適用
	if opt.WithProposedDates {
		filterEvent = filterEvent.WithProposedDates(func(query *ent.ProposedDateQuery) {
			if opt.ProposedDateOffset > 0 {
				query.Offset(opt.ProposedDateOffset)
			}
			if opt.ProposedDateLimit > 0 {
				query.Limit(opt.ProposedDateLimit)
			}
		})
	}

	return filterEvent.All(ctx)
}

func (r *EventRepositoryImpl) FindBySlugAndUser(ctx context.Context, tx *ent.Tx, userID uuid.UUID, slug string, opt EventQueryOptions) (*ent.Event, error) {
	query := r.client.Event.Query()
	if tx != nil {
		query = tx.Event.Query()
	}

	if opt.WithProposedDates {
		query = query.WithProposedDates()
	}

	return query.
		Where(
			event.SlugEQ(slug),
			event.HasCalendarWith(dbCalendar.HasUserWith(user.IDEQ(userID))),
		).
		Only(ctx)
}

func (r *EventRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, googleEvent *calendar.Event, entCalendar *ent.Calendar) (*ent.Event, error) {
	eventCreate := r.client.Event.Create()
	if tx != nil {
		eventCreate = tx.Event.Create()
	}

	if googleEvent.Id != "" {
		eventCreate = eventCreate.SetGoogleEventID(googleEvent.Id)
	}

	eventCreate = eventCreate.
		SetSummary(googleEvent.Summary).
		SetDescription(googleEvent.Description).
		SetLocation(googleEvent.Location).
		SetCalendar(entCalendar)

	return eventCreate.Save(ctx)
}

func (r *EventRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt EventQueryOptions) (*ent.Event, error) {
	eventUpdate := r.client.Event.UpdateOneID(id)
	if tx != nil {
		eventUpdate = tx.Event.UpdateOneID(id)
	}

	if opt.Summary != nil {
		eventUpdate = eventUpdate.SetSummary(*opt.Summary)
	}

	if opt.Location != nil {
		eventUpdate = eventUpdate.SetLocation(*opt.Location)
	}

	if opt.Description != nil {
		eventUpdate = eventUpdate.SetDescription(*opt.Description)
	}

	if opt.Status != nil {
		status := event.Status(*opt.Status)
		eventUpdate = eventUpdate.SetStatus(status)
	}

	if opt.ConfirmedDateID != nil {
		eventUpdate = eventUpdate.SetConfirmedDateID(*opt.ConfirmedDateID)
	}

	if opt.GoogleEventID != nil {
		eventUpdate = eventUpdate.SetGoogleEventID(*opt.GoogleEventID)
	}

	return eventUpdate.Save(ctx)
}

func (r *EventRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.Event.DeleteOneID(id).Exec(ctx)
	}
	return r.client.Event.DeleteOneID(id).Exec(ctx)
}

func (r *EventRepositoryImpl) SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	softDeleteEvent := r.client.Event.UpdateOneID(id)
	if tx != nil {
		softDeleteEvent = tx.Event.UpdateOneID(id)
	}
	return softDeleteEvent.
		SetDeletedAt(time.Now()).
		Exec(ctx)
}

func (r *EventRepositoryImpl) Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	restoreEvent := r.client.Event.UpdateOneID(id)
	if tx != nil {
		restoreEvent = tx.Event.UpdateOneID(id)
	}
	return restoreEvent.
		SetNillableDeletedAt(nil).
		Exec(ctx)
}

func (r *EventRepositoryImpl) SoftDeleteWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}

	// イベントを論理削除
	if err := r.SoftDelete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to soft delete event: %w", err)
	}

	// 関連する提案日を論理削除
	proposedDateIDs, err := tx.ProposedDate.
		Query().
		Where(proposeddate.HasEventWith(event.IDEQ(id))).
		IDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to query proposed dates: %w", err)
	}
	if len(proposedDateIDs) > 0 {
		if err := tx.ProposedDate.Update().
			Where(proposeddate.IDIn(proposedDateIDs...)).
			SetDeletedAt(time.Now()).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to soft delete proposed dates: %w", err)
		}
	}

	return nil
}

func (r *EventRepositoryImpl) RestoreWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}

	// イベントを復元
	if err := r.Restore(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to restore event: %w", err)
	}

	// 関連する提案日を復元
	proposedDateIDs, err := tx.ProposedDate.
		Query().
		Where(proposeddate.HasEventWith(event.IDEQ(id))).
		IDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to query proposed dates: %w", err)
	}
	if len(proposedDateIDs) > 0 {
		if err := tx.ProposedDate.Update().
			Where(proposeddate.IDIn(proposedDateIDs...)).
			SetNillableDeletedAt(nil).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to restore proposed dates: %w", err)
		}
	}

	return nil
}

func (r *EventRepositoryImpl) SearchEvents(ctx context.Context, tx *ent.Tx, id, calendarID uuid.UUID, opt EventQueryOptions) ([]*ent.Event, error) {
	query := r.client.Event.Query()
	if tx != nil {
		query = tx.Event.Query()
	}

	query = query.Where(event.HasCalendarWith(dbCalendar.IDEQ(calendarID)))

	if opt.Summary != nil {
		query = query.Where(event.SummaryContains(*opt.Summary))
	}

	if opt.Location != nil {
		query = query.Where(event.LocationContains(*opt.Location))
	}

	if opt.Description != nil {
		query = query.Where(event.DescriptionContains(*opt.Description))
	}

	if opt.Status != nil {
		query = query.Where(event.StatusEQ(event.Status(*opt.Status)))
	}

	if opt.ConfirmedDateID != nil {
		query = query.Where(event.ConfirmedDateIDEQ(*opt.ConfirmedDateID))
	}

	if opt.GoogleEventID != nil {
		query = query.Where(event.GoogleEventIDEQ(*opt.GoogleEventID))
	}

	// イベントに対するオフセットとリミットを適用
	if opt.EventOffset > 0 {
		query = query.Offset(opt.EventOffset)
	}
	if opt.EventLimit > 0 {
		query = query.Limit(opt.EventLimit)
	}

	// イベントの提案日に対するオフセットとリミットを適用
	if opt.WithProposedDates {
		query = query.WithProposedDates(func(query *ent.ProposedDateQuery) {
			if opt.SortBy != "" {
				switch opt.SortBy {
				case "ProposedDateStart":
					if opt.SortOrder == "desc" {
						query = query.Order(ent.Desc(proposeddate.FieldStartTime))
					} else {
						query = query.Order(ent.Asc(proposeddate.FieldStartTime))
					}
				case "ProposedDateEnd":
					if opt.SortOrder == "desc" {
						query = query.Order(ent.Desc(proposeddate.FieldEndTime))
					} else {
						query = query.Order(ent.Asc(proposeddate.FieldEndTime))
					}
				case "ProposedDatePriority":
					if opt.SortOrder == "desc" {
						query = query.Order(ent.Desc(proposeddate.FieldPriority))
					} else {
						query = query.Order(ent.Asc(proposeddate.FieldPriority))
					}
				default:
				// デフォルトは StartTime 昇順
				query = query.Order(ent.Asc(proposeddate.FieldStartTime))
				}
			}

			if opt.ProposedDateOffset > 0 {
				query = query.Offset(opt.ProposedDateOffset)
			}
			if opt.ProposedDateLimit > 0 {
				query = query.Limit(opt.ProposedDateLimit)
			}

			if opt.ProposedDateStartGTE != nil {
				query = query.Where(proposeddate.StartTimeGTE(*opt.ProposedDateStartGTE))
			}

			if opt.ProposedDateStartLTE != nil {
				query = query.Where(proposeddate.StartTimeLTE(*opt.ProposedDateStartLTE))
			}

			if opt.ProposedDateEndGTE != nil {
				query = query.Where(proposeddate.EndTimeGTE(*opt.ProposedDateEndGTE))
			}

			if opt.ProposedDateEndLTE != nil {
				query = query.Where(proposeddate.EndTimeLTE(*opt.ProposedDateEndLTE))
			}
		})
	}

	return query.All(ctx)
}
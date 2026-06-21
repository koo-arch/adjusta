package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
)

type EventRepository = repoEvent.EventRepository
type EventCreateOptions = repoEvent.EventCreateOptions
type EventFilterOptions = repoEvent.EventFilterOptions
type EventReadOptions = repoEvent.EventReadOptions
type EventSearchOptions = repoEvent.EventSearchOptions
type EventUpdateOptions = repoEvent.EventUpdateOptions

type EventRepositoryImpl struct {
	client *ent.Client
}

func NewEventRepository(client *ent.Client) *EventRepositoryImpl {
	return &EventRepositoryImpl{
		client: client,
	}
}

func (r *EventRepositoryImpl) Read(ctx context.Context, id uuid.UUID, opt EventReadOptions) (*repoEvent.Event, error) {
	query := r.client.Event.Query()

	if opt.WithProposedDates {
		query = query.WithProposedDates()
	}

	entity, err := query.Where(event.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toEvent(entity), nil
}

func (r *EventRepositoryImpl) FilterByCalendarID(ctx context.Context, calendarID uuid.UUID, opt EventFilterOptions) ([]*repoEvent.Event, error) {
	filterEvent := r.client.Event.Query()

	filterEvent = filterEvent.Where(event.PrimaryCalendarIDEQ(calendarID))

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

	entities, err := filterEvent.All(ctx)
	if err != nil {
		return nil, err
	}
	return toEvents(entities), nil
}

func (r *EventRepositoryImpl) FindByIDAndUser(ctx context.Context, userID, eventID uuid.UUID, opt EventReadOptions) (*repoEvent.Event, error) {
	query := r.client.Event.Query()

	if opt.WithProposedDates {
		query = query.WithProposedDates()
	}

	entity, err := query.
		Where(
			event.IDEQ(eventID),
			event.UserIDEQ(userID),
		).
		Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toEvent(entity), nil
}

func (r *EventRepositoryImpl) Create(ctx context.Context, userID uuid.UUID, opt EventCreateOptions, primaryCalendarID uuid.UUID) (*repoEvent.Event, error) {
	eventCreate := r.client.Event.Create()

	eventCreate = eventCreate.
		SetUserID(userID).
		SetPrimaryCalendarID(primaryCalendarID).
		SetTitle(opt.Title).
		SetDescription(opt.Description).
		SetLocation(opt.Location)

	entity, err := eventCreate.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toEvent(entity), nil
}

func (r *EventRepositoryImpl) Update(ctx context.Context, id uuid.UUID, opt EventUpdateOptions) (*repoEvent.Event, error) {
	eventUpdate := r.client.Event.UpdateOneID(id)

	if opt.Title != nil {
		eventUpdate = eventUpdate.SetTitle(*opt.Title)
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

	if opt.SyncStatus != nil {
		syncStatus := event.SyncStatus(*opt.SyncStatus)
		eventUpdate = eventUpdate.SetSyncStatus(syncStatus)
	}

	if opt.ConfirmedDateID != nil {
		eventUpdate = eventUpdate.SetConfirmedDateID(*opt.ConfirmedDateID)
	}

	if opt.ConfirmedGoogleEventID != nil {
		eventUpdate = eventUpdate.SetConfirmedGoogleEventID(*opt.ConfirmedGoogleEventID)
	}

	if opt.ClearLastSyncedAt {
		eventUpdate = eventUpdate.ClearLastSyncedAt()
	}
	if opt.LastSyncedAt != nil {
		eventUpdate = eventUpdate.SetLastSyncedAt(*opt.LastSyncedAt)
	}

	if opt.ClearLastSyncError {
		eventUpdate = eventUpdate.ClearLastSyncError()
	}
	if opt.LastSyncError != nil {
		eventUpdate = eventUpdate.SetLastSyncError(*opt.LastSyncError)
	}

	entity, err := eventUpdate.Save(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toEvent(entity), nil
}

func (r *EventRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.Event.DeleteOneID(id).Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *EventRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.client.Event.UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *EventRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	err := r.client.Event.UpdateOneID(id).
		SetNillableDeletedAt(nil).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *EventRepositoryImpl) SearchEvents(ctx context.Context, userID, calendarID uuid.UUID, opt EventSearchOptions) ([]*repoEvent.Event, error) {
	query := r.client.Event.Query()

	query = query.Where(
		event.UserIDEQ(userID),
		event.PrimaryCalendarIDEQ(calendarID),
	)

	if opt.Title != nil {
		query = query.Where(event.TitleContains(*opt.Title))
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

	if opt.SyncStatus != nil {
		query = query.Where(event.SyncStatusEQ(event.SyncStatus(*opt.SyncStatus)))
	}

	if opt.ConfirmedDateID != nil {
		query = query.Where(event.ConfirmedDateIDEQ(*opt.ConfirmedDateID))
	}

	if opt.ConfirmedGoogleEventID != nil {
		query = query.Where(event.ConfirmedGoogleEventIDEQ(*opt.ConfirmedGoogleEventID))
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

	entities, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return toEvents(entities), nil
}

func toEvent(entity *ent.Event) *repoEvent.Event {
	if entity == nil {
		return nil
	}

	return &repoEvent.Event{
		UserID:                 entity.UserID,
		ID:                     entity.ID,
		PrimaryCalendarID:      entity.PrimaryCalendarID,
		Title:                  entity.Title,
		Location:               entity.Location,
		Description:            entity.Description,
		Status:                 value.EventStatus(entity.Status),
		ConfirmedDateID:        entity.ConfirmedDateID,
		ConfirmedGoogleEventID: entity.ConfirmedGoogleEventID,
		SyncStatus:             value.SyncStatus(entity.SyncStatus),
		LastSyncedAt:           entity.LastSyncedAt,
		LastSyncError:          entity.LastSyncError,
		ProposedDates:          toEventProposedDates(entity.Edges.ProposedDates),
	}
}

func toEvents(entities []*ent.Event) []*repoEvent.Event {
	events := make([]*repoEvent.Event, 0, len(entities))
	for _, entity := range entities {
		events = append(events, toEvent(entity))
	}
	return events
}

func toEventProposedDates(entities []*ent.ProposedDate) []*repoProposedDate.ProposedDate {
	dates := make([]*repoProposedDate.ProposedDate, 0, len(entities))
	for _, entity := range entities {
		dates = append(dates, &repoProposedDate.ProposedDate{
			ID:            entity.ID,
			EventID:       entity.EventID,
			GoogleEventID: entity.GoogleEventID,
			StartTime:     entity.StartTime,
			EndTime:       entity.EndTime,
			Priority:      entity.Priority,
			Status:        value.ProposedDateStatus(entity.Status),
			SyncStatus:    value.SyncStatus(entity.SyncStatus),
			LastSyncedAt:  entity.LastSyncedAt,
			LastSyncError: entity.LastSyncError,
		})
	}
	return dates
}

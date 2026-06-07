package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	dbCalendar "github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
	dbUserCalendar "github.com/koo-arch/adjusta-backend/ent/usercalendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"google.golang.org/api/calendar/v3"
)

type EventRepository = repoEvent.EventRepository
type EventQueryOptions = repoEvent.EventQueryOptions

type EventRepositoryImpl struct {
	client *ent.Client
}

func NewEventRepository(client *ent.Client) *EventRepositoryImpl {
	return &EventRepositoryImpl{
		client: client,
	}
}

func (r *EventRepositoryImpl) WithTx(tx transaction.Tx) EventRepository {
	return &EventRepositoryImpl{client: tx.Client()}
}

func (r *EventRepositoryImpl) Read(ctx context.Context, id uuid.UUID, opt EventQueryOptions) (*repositorymodel.StoredEvent, error) {
	query := r.client.Event.Query()

	if opt.WithProposedDates {
		query = query.WithProposedDates()
	}

	entity, err := query.Where(event.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toStoredEvent(entity), nil
}

func (r *EventRepositoryImpl) FilterByCalendarID(ctx context.Context, calendarID uuid.UUID, opt EventQueryOptions) ([]*repositorymodel.StoredEvent, error) {
	filterEvent := r.client.Event.Query()

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

	entities, err := filterEvent.All(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredEvents(entities), nil
}

func (r *EventRepositoryImpl) FindBySlugAndUser(ctx context.Context, userID uuid.UUID, slug string, opt EventQueryOptions) (*repositorymodel.StoredEvent, error) {
	query := r.client.Event.Query()

	if opt.WithProposedDates {
		query = query.WithProposedDates()
	}

	entity, err := query.
		Where(
			event.SlugEQ(slug),
			event.HasCalendarWith(dbCalendar.HasUserCalendarsWith(dbUserCalendar.UserIDEQ(userID))),
		).
		Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toStoredEvent(entity), nil
}

func (r *EventRepositoryImpl) Create(ctx context.Context, googleEvent *calendar.Event, calendarID uuid.UUID) (*repositorymodel.StoredEvent, error) {
	eventCreate := r.client.Event.Create()

	if googleEvent.Id != "" {
		eventCreate = eventCreate.SetGoogleEventID(googleEvent.Id)
	}

	eventCreate = eventCreate.
		SetSummary(googleEvent.Summary).
		SetDescription(googleEvent.Description).
		SetLocation(googleEvent.Location).
		SetCalendarID(calendarID)

	entity, err := eventCreate.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredEvent(entity), nil
}

func (r *EventRepositoryImpl) Update(ctx context.Context, id uuid.UUID, opt EventQueryOptions) (*repositorymodel.StoredEvent, error) {
	eventUpdate := r.client.Event.UpdateOneID(id)

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

	entity, err := eventUpdate.Save(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toStoredEvent(entity), nil
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

func (r *EventRepositoryImpl) SearchEvents(ctx context.Context, id, calendarID uuid.UUID, opt EventQueryOptions) ([]*repositorymodel.StoredEvent, error) {
	query := r.client.Event.Query()

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

	entities, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredEvents(entities), nil
}

func toStoredEvent(entity *ent.Event) *repositorymodel.StoredEvent {
	if entity == nil {
		return nil
	}

	return &repositorymodel.StoredEvent{
		ID:              entity.ID,
		Summary:         entity.Summary,
		Location:        entity.Location,
		Description:     entity.Description,
		Status:          domainvalue.EventStatus(entity.Status),
		ConfirmedDateID: entity.ConfirmedDateID,
		GoogleEventID:   entity.GoogleEventID,
		Slug:            entity.Slug,
		ProposedDates:   toStoredEventProposedDates(entity.Edges.ProposedDates),
	}
}

func toStoredEvents(entities []*ent.Event) []*repositorymodel.StoredEvent {
	storedEvents := make([]*repositorymodel.StoredEvent, 0, len(entities))
	for _, entity := range entities {
		storedEvents = append(storedEvents, toStoredEvent(entity))
	}
	return storedEvents
}

func toStoredEventProposedDates(entities []*ent.ProposedDate) []*repositorymodel.StoredProposedDate {
	storedDates := make([]*repositorymodel.StoredProposedDate, 0, len(entities))
	for _, entity := range entities {
		storedDates = append(storedDates, &repositorymodel.StoredProposedDate{
			ID:        entity.ID,
			EventID:   entity.EventID,
			StartTime: entity.StartTime,
			EndTime:   entity.EndTime,
			Priority:  entity.Priority,
		})
	}
	return storedDates
}

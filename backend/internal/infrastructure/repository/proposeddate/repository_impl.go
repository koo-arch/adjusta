package proposeddate

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type ProposedDateRepository = repoProposedDate.ProposedDateRepository
type ProposedDateQueryOptions = repoProposedDate.ProposedDateQueryOptions

type ProposedDateRepositoryImpl struct {
	client *ent.Client
}

func NewProposedDateRepository(client *ent.Client) *ProposedDateRepositoryImpl {
	return &ProposedDateRepositoryImpl{
		client: client,
	}
}

func (r *ProposedDateRepositoryImpl) WithTx(tx transaction.Tx) ProposedDateRepository {
	return &ProposedDateRepositoryImpl{client: tx.Client()}
}

func (r *ProposedDateRepositoryImpl) Read(ctx context.Context, id uuid.UUID) (*repositorymodel.StoredProposedDate, error) {
	entity, err := r.client.ProposedDate.Get(ctx, id)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toStoredProposedDate(entity), nil
}

func (r *ProposedDateRepositoryImpl) FilterByEventID(ctx context.Context, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
	entities, err := r.client.ProposedDate.Query().
		Where(proposeddate.HasEventWith(event.IDEQ(eventID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredProposedDates(entities), nil
}

func (r *ProposedDateRepositoryImpl) ExclusionEventID(ctx context.Context, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
	entities, err := r.client.ProposedDate.Query().
		Where(proposeddate.Not(proposeddate.HasEventWith(event.IDEQ(eventID)))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredProposedDates(entities), nil
}

func (r *ProposedDateRepositoryImpl) Create(ctx context.Context, opt ProposedDateQueryOptions, eventID uuid.UUID) (*repositorymodel.StoredProposedDate, error) {
	proposedDateCreate := r.client.ProposedDate.Create()

	proposedDateCreate = proposedDateCreate.
		SetStartTime(*opt.StartTime).
		SetEndTime(*opt.EndTime).
		SetPriority(*opt.Priority).
		SetEventID(eventID)

	if opt.GoogleEventID != nil {
		proposedDateCreate = proposedDateCreate.SetGoogleEventID(*opt.GoogleEventID)
	}
	if opt.Status != nil {
		proposedDateCreate = proposedDateCreate.SetStatus(proposeddate.Status(*opt.Status))
	}
	if opt.SyncStatus != nil {
		proposedDateCreate = proposedDateCreate.SetSyncStatus(proposeddate.SyncStatus(*opt.SyncStatus))
	}
	if opt.LastSyncedAt != nil {
		proposedDateCreate = proposedDateCreate.SetLastSyncedAt(*opt.LastSyncedAt)
	}
	if opt.LastSyncError != nil {
		proposedDateCreate = proposedDateCreate.SetLastSyncError(*opt.LastSyncError)
	}

	entity, err := proposedDateCreate.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredProposedDate(entity), nil
}

func (r *ProposedDateRepositoryImpl) Update(ctx context.Context, id uuid.UUID, opt ProposedDateQueryOptions) (*repositorymodel.StoredProposedDate, error) {
	proposedDateUpdate := r.client.ProposedDate.UpdateOneID(id)

	if opt.StartTime != nil {
		proposedDateUpdate = proposedDateUpdate.SetStartTime(*opt.StartTime)
	}

	if opt.EndTime != nil {
		proposedDateUpdate = proposedDateUpdate.SetEndTime(*opt.EndTime)
	}

	if opt.Priority != nil {
		proposedDateUpdate = proposedDateUpdate.SetPriority(*opt.Priority)
	}

	if opt.GoogleEventID != nil {
		proposedDateUpdate = proposedDateUpdate.SetGoogleEventID(*opt.GoogleEventID)
	}

	if opt.Status != nil {
		proposedDateUpdate = proposedDateUpdate.SetStatus(proposeddate.Status(*opt.Status))
	}

	if opt.SyncStatus != nil {
		proposedDateUpdate = proposedDateUpdate.SetSyncStatus(proposeddate.SyncStatus(*opt.SyncStatus))
	}

	if opt.ClearLastSyncedAt {
		proposedDateUpdate = proposedDateUpdate.ClearLastSyncedAt()
	}
	if opt.LastSyncedAt != nil {
		proposedDateUpdate = proposedDateUpdate.SetLastSyncedAt(*opt.LastSyncedAt)
	}

	if opt.ClearLastSyncError {
		proposedDateUpdate = proposedDateUpdate.ClearLastSyncError()
	}
	if opt.LastSyncError != nil {
		proposedDateUpdate = proposedDateUpdate.SetLastSyncError(*opt.LastSyncError)
	}

	entity, err := proposedDateUpdate.Save(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toStoredProposedDate(entity), nil
}

func (r *ProposedDateRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.ProposedDate.DeleteOneID(id).Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *ProposedDateRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.client.ProposedDate.UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *ProposedDateRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	err := r.client.ProposedDate.UpdateOneID(id).
		SetNillableDeletedAt(nil).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *ProposedDateRepositoryImpl) CreateBulk(ctx context.Context, selectedDates []appmodel.SelectedDate, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error) {
	var proposedDateCreates []*ent.ProposedDateCreate

	for _, selectedDate := range selectedDates {
		proposedDateCreate := r.client.ProposedDate.Create()

		proposedDateCreate = proposedDateCreate.
			SetStartTime(selectedDate.Start).
			SetEndTime(selectedDate.End).
			SetPriority(selectedDate.Priority).
			SetStatus(proposeddate.Status(domainvalue.ProposedDateStatusActive)).
			SetEventID(eventID)

		proposedDateCreates = append(proposedDateCreates, proposedDateCreate)
	}

	entities, err := r.client.ProposedDate.CreateBulk(proposedDateCreates...).Save(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredProposedDates(entities), nil
}

func toStoredProposedDate(entity *ent.ProposedDate) *repositorymodel.StoredProposedDate {
	if entity == nil {
		return nil
	}

	return &repositorymodel.StoredProposedDate{
		ID:            entity.ID,
		EventID:       entity.EventID,
		GoogleEventID: entity.GoogleEventID,
		StartTime:     entity.StartTime,
		EndTime:       entity.EndTime,
		Priority:      entity.Priority,
		Status:        domainvalue.ProposedDateStatus(entity.Status),
		SyncStatus:    domainvalue.SyncStatus(entity.SyncStatus),
		LastSyncedAt:  entity.LastSyncedAt,
		LastSyncError: entity.LastSyncError,
	}
}

func toStoredProposedDates(entities []*ent.ProposedDate) []*repositorymodel.StoredProposedDate {
	storedDates := make([]*repositorymodel.StoredProposedDate, 0, len(entities))
	for _, entity := range entities {
		storedDates = append(storedDates, toStoredProposedDate(entity))
	}
	return storedDates
}

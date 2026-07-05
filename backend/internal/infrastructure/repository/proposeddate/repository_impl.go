package proposeddate

import (
	"context"
	"time"

	"github.com/google/uuid"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/event"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/proposeddate"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
)

type ProposedDateRepository = repoProposedDate.ProposedDateRepository
type ProposedDateCreateOptions = repoProposedDate.ProposedDateCreateOptions
type ProposedDateUpdateOptions = repoProposedDate.ProposedDateUpdateOptions

type ProposedDateRepositoryImpl struct {
	client *ent.Client
}

func NewProposedDateRepository(client *ent.Client) *ProposedDateRepositoryImpl {
	return &ProposedDateRepositoryImpl{
		client: client,
	}
}

func (r *ProposedDateRepositoryImpl) Read(ctx context.Context, id uuid.UUID) (*repoProposedDate.ProposedDate, error) {
	entity, err := r.client.ProposedDate.Get(ctx, id)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toProposedDate(entity), nil
}

func (r *ProposedDateRepositoryImpl) FilterByEventID(ctx context.Context, eventID uuid.UUID) ([]*repoProposedDate.ProposedDate, error) {
	entities, err := r.client.ProposedDate.Query().
		Where(proposeddate.HasEventWith(event.IDEQ(eventID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toProposedDates(entities), nil
}

func (r *ProposedDateRepositoryImpl) Create(ctx context.Context, opt ProposedDateCreateOptions, eventID uuid.UUID) (*repoProposedDate.ProposedDate, error) {
	proposedDateCreate := r.client.ProposedDate.Create()

	proposedDateCreate = proposedDateCreate.
		SetStartTime(opt.StartTime).
		SetEndTime(opt.EndTime).
		SetPriority(opt.Priority).
		SetEventID(eventID)

	proposedDateCreate = applyProposedDateCreateOptions(proposedDateCreate, opt)

	entity, err := proposedDateCreate.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toProposedDate(entity), nil
}

func applyProposedDateCreateOptions(proposedDateCreate *ent.ProposedDateCreate, opt ProposedDateCreateOptions) *ent.ProposedDateCreate {
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

	return proposedDateCreate
}

func (r *ProposedDateRepositoryImpl) Update(ctx context.Context, id uuid.UUID, opt ProposedDateUpdateOptions) (*repoProposedDate.ProposedDate, error) {
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
	return toProposedDate(entity), nil
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

func (r *ProposedDateRepositoryImpl) CreateBulk(ctx context.Context, opts []ProposedDateCreateOptions, eventID uuid.UUID) ([]*repoProposedDate.ProposedDate, error) {
	var proposedDateCreates []*ent.ProposedDateCreate

	for _, opt := range opts {
		proposedDateCreate := r.client.ProposedDate.Create()

		proposedDateCreate = proposedDateCreate.
			SetStartTime(opt.StartTime).
			SetEndTime(opt.EndTime).
			SetPriority(opt.Priority).
			SetEventID(eventID)

		proposedDateCreate = applyProposedDateCreateOptions(proposedDateCreate, opt)

		proposedDateCreates = append(proposedDateCreates, proposedDateCreate)
	}

	entities, err := r.client.ProposedDate.CreateBulk(proposedDateCreates...).Save(ctx)
	if err != nil {
		return nil, err
	}
	return toProposedDates(entities), nil
}

func toProposedDate(entity *ent.ProposedDate) *repoProposedDate.ProposedDate {
	if entity == nil {
		return nil
	}

	return &repoProposedDate.ProposedDate{
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
	}
}

func toProposedDates(entities []*ent.ProposedDate) []*repoProposedDate.ProposedDate {
	dates := make([]*repoProposedDate.ProposedDate, 0, len(entities))
	for _, entity := range entities {
		dates = append(dates, toProposedDate(entity))
	}
	return dates
}

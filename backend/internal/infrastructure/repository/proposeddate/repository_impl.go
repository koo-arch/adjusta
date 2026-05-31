package proposeddate

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
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
			SetEventID(eventID)

		proposedDateCreates = append(proposedDateCreates, proposedDateCreate)
	}

	entities, err := r.client.ProposedDate.CreateBulk(proposedDateCreates...).Save(ctx)
	if err != nil {
		return nil, err
	}
	return toStoredProposedDates(entities), nil
}

func (r *ProposedDateRepositoryImpl) DecrementPriorityExceptID(ctx context.Context, eventID, excludeID uuid.UUID) error {
	_, err := r.client.ProposedDate.Update().Where(
		proposeddate.HasEventWith(event.IDEQ(eventID)),
		proposeddate.IDNEQ(excludeID),
	).
		AddPriority(1).
		Save(ctx)

	return err
}

func (r *ProposedDateRepositoryImpl) ReorderPriority(ctx context.Context, eventID uuid.UUID) error {
	// eventIDに紐づくProposedDateをpriority順に取得
	proposedDates, err := r.client.ProposedDate.Query().
		Where(proposeddate.HasEventWith(event.IDEQ(eventID))).
		Order(ent.Asc(proposeddate.FieldPriority)).
		All(ctx)
	if err != nil {
		return err
	}

	// priorityが連番でない時に振り直す
	if !r.isSequential(proposedDates) {
		return r.updateToSequentialPriority(ctx, proposedDates)
	}

	return nil
}

func (r *ProposedDateRepositoryImpl) isSequential(proposedDates []*ent.ProposedDate) bool {
	for i, proposedDate := range proposedDates {
		if proposedDate.Priority != i+1 {
			return false
		}
	}
	return true
}

func (r *ProposedDateRepositoryImpl) updateToSequentialPriority(ctx context.Context, proposedDates []*ent.ProposedDate) error {
	for i, proposedDate := range proposedDates {
		priority := i + 1
		_, err := r.Update(ctx, proposedDate.ID, ProposedDateQueryOptions{
			Priority: &priority,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func toStoredProposedDate(entity *ent.ProposedDate) *repositorymodel.StoredProposedDate {
	if entity == nil {
		return nil
	}

	return &repositorymodel.StoredProposedDate{
		ID:        entity.ID,
		EventID:   entity.EventID,
		StartTime: entity.StartTime,
		EndTime:   entity.EndTime,
		Priority:  entity.Priority,
	}
}

func toStoredProposedDates(entities []*ent.ProposedDate) []*repositorymodel.StoredProposedDate {
	storedDates := make([]*repositorymodel.StoredProposedDate, 0, len(entities))
	for _, entity := range entities {
		storedDates = append(storedDates, toStoredProposedDate(entity))
	}
	return storedDates
}

package proposeddate

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type ProposedDateRepositoryImpl struct {
	client *ent.Client
}

func NewProposedDateRepository(client *ent.Client) *ProposedDateRepositoryImpl {
	return &ProposedDateRepositoryImpl{
		client: client,
	}
}

func (r *ProposedDateRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.ProposedDate, error) {
	if tx != nil {
		return tx.ProposedDate.Get(ctx, id)
	}
	return r.client.ProposedDate.Get(ctx, id)
}

func (r *ProposedDateRepositoryImpl) FilterByEventID(ctx context.Context, tx *ent.Tx, eventID uuid.UUID) ([]*ent.ProposedDate, error) {
	filterProposedDate := r.client.ProposedDate.Query()
	if tx != nil {
		filterProposedDate = tx.ProposedDate.Query()
	}
	return filterProposedDate.
		Where(proposeddate.HasEventWith(event.IDEQ(eventID))).
		All(ctx)
}

func (r* ProposedDateRepositoryImpl) ExclusionEventID(ctx context.Context, tx *ent.Tx, eventID uuid.UUID) ([]*ent.ProposedDate, error) {
	filterProposedDate := r.client.ProposedDate.Query()
	if tx != nil {
		filterProposedDate = tx.ProposedDate.Query()
	}
	return filterProposedDate.
		Where(proposeddate.Not(proposeddate.HasEventWith(event.IDEQ(eventID)))).
		All(ctx)
}

func (r *ProposedDateRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, opt ProposedDateQueryOptions, entEvent *ent.Event) (*ent.ProposedDate, error) {
	proposedDateCreate := r.client.ProposedDate.Create()
	if tx != nil {
		proposedDateCreate = tx.ProposedDate.Create()
	}

	proposedDateCreate = proposedDateCreate.
		SetStartTime(*opt.StartTime).
		SetEndTime(*opt.EndTime).
		SetPriority(*opt.Priority).
		SetEvent(entEvent)

	return proposedDateCreate.Save(ctx)
}

func (r *ProposedDateRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt ProposedDateQueryOptions) (*ent.ProposedDate, error) {
	proposedDateUpdate := r.client.ProposedDate.UpdateOneID(id)
	if tx != nil {
		proposedDateUpdate = tx.ProposedDate.UpdateOneID(id)
	}

	if opt.StartTime!= nil {
		proposedDateUpdate = proposedDateUpdate.SetStartTime(*opt.StartTime)
	}

	if opt.EndTime != nil {
		proposedDateUpdate = proposedDateUpdate.SetEndTime(*opt.EndTime)
	}

	if opt.Priority != nil {
		proposedDateUpdate = proposedDateUpdate.SetPriority(*opt.Priority)
	}

	return proposedDateUpdate.Save(ctx)
}

func (r *ProposedDateRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.ProposedDate.DeleteOneID(id).Exec(ctx)
	}
	return r.client.ProposedDate.DeleteOneID(id).Exec(ctx)
}

func (r *ProposedDateRepositoryImpl) SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	softDeleteProposedDate := r.client.ProposedDate.UpdateOneID(id)
	if tx != nil {
		softDeleteProposedDate = tx.ProposedDate.UpdateOneID(id)
	}
	return softDeleteProposedDate.
		SetDeletedAt(time.Now()).
		Exec(ctx)
}

func (r *ProposedDateRepositoryImpl) Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	restoreProposedDate := r.client.ProposedDate.UpdateOneID(id)
	if tx != nil {
		restoreProposedDate = tx.ProposedDate.UpdateOneID(id)
	}
	return restoreProposedDate.
		SetNillableDeletedAt(nil).
		Exec(ctx)
}

func (r *ProposedDateRepositoryImpl) CreateBulk(ctx context.Context, tx *ent.Tx, selectedDates []models.SelectedDate, entEvent *ent.Event) ([]*ent.ProposedDate, error) {
	var proposedDateCreates []*ent.ProposedDateCreate

	for _, selectedDate := range selectedDates {
		proposedDateCreate := r.client.ProposedDate.Create()
		if tx != nil {
			proposedDateCreate = tx.ProposedDate.Create()
		}

		proposedDateCreate = proposedDateCreate.
			SetStartTime(selectedDate.Start).
			SetEndTime(selectedDate.End).
			SetPriority(selectedDate.Priority).
			SetEvent(entEvent)

		proposedDateCreates = append(proposedDateCreates, proposedDateCreate)
	}

	if (tx != nil) {
		return tx.ProposedDate.CreateBulk(proposedDateCreates...).Save(ctx)
	}

	return r.client.ProposedDate.CreateBulk(proposedDateCreates...).Save(ctx)
}


func (r *ProposedDateRepositoryImpl) DecrementPriorityExceptID(ctx context.Context, tx *ent.Tx, eventID, excludeID uuid.UUID) error {
	update := r.client.ProposedDate.Update()
	if tx != nil {
		update = tx.ProposedDate.Update()
	}

	_, err := update.Where(
		proposeddate.HasEventWith(event.IDEQ(eventID)),
		proposeddate.IDNEQ(excludeID),
	).
		AddPriority(1).
		Save(ctx)

	return err
}

func (r *ProposedDateRepositoryImpl) ReorderPriority(ctx context.Context, tx *ent.Tx, eventID uuid.UUID) error {
	query := r.client.ProposedDate.Query()
	if tx != nil {
		query = tx.ProposedDate.Query()
	}

	// eventIDに紐づくProposedDateをpriority順に取得
	proposedDates, err := query.
		Where(proposeddate.HasEventWith(event.IDEQ(eventID))).
		Order(ent.Asc(proposeddate.FieldPriority)).
		All(ctx)
	if err != nil {
		return err
	}

	// priorityが連番でない時に振り直す
	if !r.isSequential(proposedDates) {
		return r.updateToSequentialPriority(ctx, tx, proposedDates)
	}

	return nil
}

func (r *ProposedDateRepositoryImpl) isSequential(proposedDates []*ent.ProposedDate) bool {
	for i, proposedDate := range proposedDates {
		if proposedDate.Priority != i + 1 {
			return false
		}
	}
	return true
}

func (r *ProposedDateRepositoryImpl) updateToSequentialPriority (ctx context.Context, tx *ent.Tx, proposedDates []*ent.ProposedDate) error {
	for i, proposedDate := range proposedDates {
		priority := i + 1
		_, err := r.Update(ctx, tx, proposedDate.ID, ProposedDateQueryOptions{
			Priority: &priority,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
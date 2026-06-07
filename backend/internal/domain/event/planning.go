package event

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

var ErrInvalidDateRange = errors.New("invalid date range")

type DraftProposedDate struct {
	ID       *uuid.UUID
	Start    time.Time
	End      time.Time
	Priority int
}

type ExistingProposedDate struct {
	ID       uuid.UUID
	Start    time.Time
	End      time.Time
	Priority int
}

type ProposedDateCreate struct {
	Start    time.Time
	End      time.Time
	Priority int
}

type ProposedDateUpdate struct {
	ID       uuid.UUID
	Start    time.Time
	End      time.Time
	Priority int
}

type ProposedDateChangeSet struct {
	Creates []ProposedDateCreate
	Updates []ProposedDateUpdate
	Deletes []uuid.UUID
}

type PriorityAdjustment string

const (
	PriorityAdjustmentNone             PriorityAdjustment = ""
	PriorityAdjustmentIncrementOthers  PriorityAdjustment = "increment_others"
	PriorityAdjustmentReorderRemaining PriorityAdjustment = "reorder_remaining"
)

type ConfirmationPlan struct {
	Status             domainvalue.EventStatus
	Create             *ProposedDateCreate
	Update             *ProposedDateUpdate
	PriorityAdjustment PriorityAdjustment
}

func BuildProposedDateChangeSet(requested []DraftProposedDate, existing []ExistingProposedDate) ProposedDateChangeSet {
	changeSet := ProposedDateChangeSet{
		Creates: make([]ProposedDateCreate, 0),
		Updates: make([]ProposedDateUpdate, 0),
		Deletes: make([]uuid.UUID, 0),
	}

	requestedByID := make(map[uuid.UUID]DraftProposedDate, len(requested))
	for _, date := range requested {
		if date.ID == nil {
			changeSet.Creates = append(changeSet.Creates, ProposedDateCreate{
				Start:    date.Start,
				End:      date.End,
				Priority: date.Priority,
			})
			continue
		}

		requestedByID[*date.ID] = date
	}

	for _, date := range existing {
		if requested, ok := requestedByID[date.ID]; ok {
			changeSet.Updates = append(changeSet.Updates, ProposedDateUpdate{
				ID:       date.ID,
				Start:    requested.Start,
				End:      requested.End,
				Priority: requested.Priority,
			})
			delete(requestedByID, date.ID)
			continue
		}

		changeSet.Deletes = append(changeSet.Deletes, date.ID)
	}

	return changeSet
}

func BuildConfirmationPlan(confirmDate DraftProposedDate) (*ConfirmationPlan, error) {
	if !confirmDate.Start.Before(confirmDate.End) {
		return nil, ErrInvalidDateRange
	}

	plan := &ConfirmationPlan{
		Status: domainvalue.StatusConfirmed,
	}

	if confirmDate.ID == nil {
		plan.Create = &ProposedDateCreate{
			Start:    confirmDate.Start,
			End:      confirmDate.End,
			Priority: 1,
		}
		plan.PriorityAdjustment = PriorityAdjustmentIncrementOthers
		return plan, nil
	}

	plan.Update = &ProposedDateUpdate{
		ID:       *confirmDate.ID,
		Start:    confirmDate.Start,
		End:      confirmDate.End,
		Priority: 0,
	}
	plan.PriorityAdjustment = PriorityAdjustmentReorderRemaining

	return plan, nil
}

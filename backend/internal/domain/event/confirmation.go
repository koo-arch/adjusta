package event

import (
	"errors"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

var ErrInvalidDateRange = errors.New("invalid date range")

type ConfirmationChangeSet struct {
	Status          domainvalue.EventStatus
	Create          *ProposedDateCreate
	Update          *ProposedDateUpdate
	MarkNotSelected []uuid.UUID
}

func PlanConfirmationChanges(confirmDate DraftProposedDate, existing []ExistingProposedDate) (*ConfirmationChangeSet, error) {
	if !confirmDate.Start.Before(confirmDate.End) {
		return nil, ErrInvalidDateRange
	}

	changeSet := &ConfirmationChangeSet{
		Status:          domainvalue.StatusConfirmed,
		MarkNotSelected: make([]uuid.UUID, 0, len(existing)),
	}

	highestPriority := nextHighestPriorityExcluding(existing, confirmDate.ID)
	if confirmDate.ID == nil {
		changeSet.Create = &ProposedDateCreate{
			Start:    confirmDate.Start,
			End:      confirmDate.End,
			Priority: highestPriority,
		}
	} else {
		changeSet.Update = &ProposedDateUpdate{
			ID:       *confirmDate.ID,
			Start:    confirmDate.Start,
			End:      confirmDate.End,
			Priority: highestPriority,
		}
	}

	for _, date := range existing {
		if confirmDate.ID != nil && date.ID == *confirmDate.ID {
			continue
		}
		changeSet.MarkNotSelected = append(changeSet.MarkNotSelected, date.ID)
	}

	return changeSet, nil
}

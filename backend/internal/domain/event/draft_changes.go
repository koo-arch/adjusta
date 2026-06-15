package event

import (
	"time"

	"github.com/google/uuid"
)

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

func PlanProposedDateChanges(requested []DraftProposedDate, existing []ExistingProposedDate) ProposedDateChangeSet {
	changeSet := ProposedDateChangeSet{
		Creates: make([]ProposedDateCreate, 0),
		Updates: make([]ProposedDateUpdate, 0),
		Deletes: make([]uuid.UUID, 0),
	}

	existingByID := make(map[uuid.UUID]ExistingProposedDate, len(existing))
	for _, date := range existing {
		existingByID[date.ID] = date
	}

	requested = reconcileRequestedPriorities(requested, existing)
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
		requestedDate, ok := requestedByID[date.ID]
		if !ok {
			changeSet.Deletes = append(changeSet.Deletes, date.ID)
			continue
		}

		if date.Start.Equal(requestedDate.Start) && date.End.Equal(requestedDate.End) && date.Priority == requestedDate.Priority {
			continue
		}

		changeSet.Updates = append(changeSet.Updates, ProposedDateUpdate{
			ID:       date.ID,
			Start:    requestedDate.Start,
			End:      requestedDate.End,
			Priority: requestedDate.Priority,
		})
	}

	// If an ID was requested but no longer exists in DB, treat it as a create fallback.
	for id, date := range requestedByID {
		if _, ok := existingByID[id]; ok {
			continue
		}
		changeSet.Creates = append(changeSet.Creates, ProposedDateCreate{
			Start:    date.Start,
			End:      date.End,
			Priority: date.Priority,
		})
	}

	return changeSet
}

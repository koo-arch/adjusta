package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type SyncChange struct {
	Status             value.SyncStatus
	LastSyncedAt       *time.Time
	ClearLastSyncedAt  bool
	LastSyncError      *string
	ClearLastSyncError bool
}

type EventChange struct {
	Status                 *value.EventStatus
	ConfirmedDateID        *uuid.UUID
	ConfirmedGoogleEventID *string
	Sync                   SyncChange
}

type ProposedDateChange struct {
	GoogleEventID *string
	Start         *time.Time
	End           *time.Time
	Priority      *int
	Status        *value.ProposedDateStatus
	Sync          SyncChange
}

func NewPendingEventChange(status *value.EventStatus) EventChange {
	return EventChange{
		Status: status,
		Sync: SyncChange{
			Status: value.SyncStatusPending,
		},
	}
}

func NewNotSyncedEventChange(status *value.EventStatus) EventChange {
	return EventChange{
		Status: status,
		Sync: SyncChange{
			Status:             value.SyncStatusNotSynced,
			ClearLastSyncError: true,
		},
	}
}

func NewDraftEventChange(status *value.EventStatus, syncExternally bool) EventChange {
	if syncExternally {
		return NewPendingEventChange(status)
	}
	return NewNotSyncedEventChange(status)
}

func NewPendingEventSyncChange() EventChange {
	return NewPendingEventChange(nil)
}

func NewPendingProposedDateChange(start, end *time.Time, priority *int, status *value.ProposedDateStatus) ProposedDateChange {
	return ProposedDateChange{
		Start:    start,
		End:      end,
		Priority: priority,
		Status:   status,
		Sync: SyncChange{
			Status: value.SyncStatusPending,
		},
	}
}

func NewNotSyncedProposedDateChange(start, end *time.Time, priority *int, status *value.ProposedDateStatus) ProposedDateChange {
	return ProposedDateChange{
		Start:    start,
		End:      end,
		Priority: priority,
		Status:   status,
		Sync: SyncChange{
			Status:             value.SyncStatusNotSynced,
			ClearLastSyncError: true,
		},
	}
}

func NewDraftProposedDateChange(start, end *time.Time, priority *int, status *value.ProposedDateStatus, syncExternally bool) ProposedDateChange {
	if syncExternally {
		return NewPendingProposedDateChange(start, end, priority, status)
	}
	return NewNotSyncedProposedDateChange(start, end, priority, status)
}

func NewConfirmedProposedDateChange(start, end *time.Time, priority *int) ProposedDateChange {
	status := value.ProposedDateStatusConfirmed
	return NewPendingProposedDateChange(start, end, priority, &status)
}

func NewNotSelectedProposedDateChange() ProposedDateChange {
	status := value.ProposedDateStatusNotSelected
	return NewPendingProposedDateChange(nil, nil, nil, &status)
}

func NewPendingProposedDateSyncChange() ProposedDateChange {
	return NewPendingProposedDateChange(nil, nil, nil, nil)
}

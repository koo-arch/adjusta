package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

type SyncChange struct {
	Status             domainvalue.SyncStatus
	LastSyncedAt       *time.Time
	ClearLastSyncedAt  bool
	LastSyncError      *string
	ClearLastSyncError bool
}

type EventChange struct {
	Status                 *domainvalue.EventStatus
	ConfirmedDateID        *uuid.UUID
	GoogleEventID          *string
	ConfirmedGoogleEventID *string
	Sync                   SyncChange
}

type ProposedDateChange struct {
	Start    *time.Time
	End      *time.Time
	Priority *int
	Status   *domainvalue.ProposedDateStatus
	Sync     SyncChange
}

func NewPendingEventChange(status *domainvalue.EventStatus) EventChange {
	return EventChange{
		Status: status,
		Sync: SyncChange{
			Status: domainvalue.SyncStatusPending,
		},
	}
}

func NewNotSyncedEventChange(status *domainvalue.EventStatus) EventChange {
	return EventChange{
		Status: status,
		Sync: SyncChange{
			Status:             domainvalue.SyncStatusNotSynced,
			ClearLastSyncError: true,
		},
	}
}

func NewSyncedEventChange(status domainvalue.EventStatus, confirmedDateID uuid.UUID, googleEventID string, syncedAt time.Time) EventChange {
	return EventChange{
		Status:                 &status,
		ConfirmedDateID:        &confirmedDateID,
		GoogleEventID:          &googleEventID,
		ConfirmedGoogleEventID: &googleEventID,
		Sync: SyncChange{
			Status:             domainvalue.SyncStatusSynced,
			LastSyncedAt:       &syncedAt,
			ClearLastSyncError: true,
		},
	}
}

func NewFailedEventChange(syncErr error) EventChange {
	lastSyncError := syncErr.Error()

	return EventChange{
		Sync: SyncChange{
			Status:        domainvalue.SyncStatusFailed,
			LastSyncError: &lastSyncError,
		},
	}
}

func NewPendingProposedDateChange(start, end *time.Time, priority *int, status *domainvalue.ProposedDateStatus) ProposedDateChange {
	return ProposedDateChange{
		Start:    start,
		End:      end,
		Priority: priority,
		Status:   status,
		Sync: SyncChange{
			Status: domainvalue.SyncStatusPending,
		},
	}
}

func NewNotSyncedProposedDateChange(start, end *time.Time, priority *int, status *domainvalue.ProposedDateStatus) ProposedDateChange {
	return ProposedDateChange{
		Start:    start,
		End:      end,
		Priority: priority,
		Status:   status,
		Sync: SyncChange{
			Status:             domainvalue.SyncStatusNotSynced,
			ClearLastSyncError: true,
		},
	}
}

func ResolveGoogleEventID(confirmedGoogleEventID *string, googleEventID string) string {
	if confirmedGoogleEventID != nil && *confirmedGoogleEventID != "" {
		return *confirmedGoogleEventID
	}
	return googleEventID
}

func ResolveReusableGoogleEventID(confirmDateID *uuid.UUID, confirmedGoogleEventID *string, requestedGoogleEventID, googleEventID string) *string {
	if confirmDateID == nil {
		return nil
	}
	if confirmedGoogleEventID != nil && *confirmedGoogleEventID != "" {
		return confirmedGoogleEventID
	}
	if requestedGoogleEventID != "" {
		return &requestedGoogleEventID
	}
	if googleEventID != "" {
		return &googleEventID
	}
	return nil
}

package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func NewSyncedEventChange(status value.EventStatus, confirmedDateID uuid.UUID, googleEventID string, syncedAt time.Time) EventChange {
	return EventChange{
		Status:                 &status,
		ConfirmedDateID:        &confirmedDateID,
		ConfirmedGoogleEventID: &googleEventID,
		Sync: SyncChange{
			Status:             value.SyncStatusSynced,
			LastSyncedAt:       &syncedAt,
			ClearLastSyncError: true,
		},
	}
}

func NewFailedEventChange(syncErr error) EventChange {
	lastSyncError := syncErr.Error()

	return EventChange{
		Sync: SyncChange{
			Status:        value.SyncStatusFailed,
			LastSyncError: &lastSyncError,
		},
	}
}

func NewSyncedProposedDateChange(googleEventID string, syncedAt time.Time) ProposedDateChange {
	return ProposedDateChange{
		GoogleEventID: &googleEventID,
		Sync: SyncChange{
			Status:             value.SyncStatusSynced,
			LastSyncedAt:       &syncedAt,
			ClearLastSyncError: true,
		},
	}
}

func NewFailedProposedDateChange(syncErr error) ProposedDateChange {
	lastSyncError := syncErr.Error()

	return ProposedDateChange{
		Sync: SyncChange{
			Status:        value.SyncStatusFailed,
			LastSyncError: &lastSyncError,
		},
	}
}

func NewSyncedEventSyncChange(syncedAt time.Time) EventChange {
	return EventChange{
		Sync: SyncChange{
			Status:             value.SyncStatusSynced,
			LastSyncedAt:       &syncedAt,
			ClearLastSyncError: true,
		},
	}
}

func ResolveGoogleEventID(confirmedGoogleEventID *string) string {
	if confirmedGoogleEventID != nil && *confirmedGoogleEventID != "" {
		return *confirmedGoogleEventID
	}
	return ""
}

func ResolveReusableGoogleEventID(confirmDateID *uuid.UUID, confirmedGoogleEventID *string, requestedGoogleEventID string) *string {
	if confirmDateID == nil {
		return nil
	}
	if confirmedGoogleEventID != nil && *confirmedGoogleEventID != "" {
		return confirmedGoogleEventID
	}
	if requestedGoogleEventID != "" {
		return &requestedGoogleEventID
	}
	return nil
}

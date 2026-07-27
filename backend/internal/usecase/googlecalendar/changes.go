package googlecalendar

import (
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	domainProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func eventNotSyncedUpdate() domainEvent.EventUpdateOptions {
	return domainEvent.EventUpdateOptions{
		SyncStatus:         syncStatus(value.SyncStatusNotSynced),
		ClearLastSyncError: true,
	}
}

func confirmedEventSyncedUpdate(confirmedDateID uuid.UUID, googleEventID string, syncedAt time.Time) domainEvent.EventUpdateOptions {
	status := value.StatusConfirmed
	return domainEvent.EventUpdateOptions{
		Status:                 &status,
		ConfirmedDateID:        &confirmedDateID,
		ConfirmedGoogleEventID: &googleEventID,
		SyncStatus:             syncStatus(value.SyncStatusSynced),
		LastSyncedAt:           &syncedAt,
		ClearLastSyncError:     true,
	}
}

func eventSyncedUpdate(syncedAt time.Time) domainEvent.EventUpdateOptions {
	return domainEvent.EventUpdateOptions{
		SyncStatus:         syncStatus(value.SyncStatusSynced),
		LastSyncedAt:       &syncedAt,
		ClearLastSyncError: true,
	}
}

func eventFailedUpdate(syncErr error) domainEvent.EventUpdateOptions {
	lastSyncError := syncErr.Error()
	return domainEvent.EventUpdateOptions{
		SyncStatus:    syncStatus(value.SyncStatusFailed),
		LastSyncError: &lastSyncError,
	}
}

func proposedDateNotSyncedUpdate() domainProposedDate.ProposedDateUpdateOptions {
	return domainProposedDate.ProposedDateUpdateOptions{
		SyncStatus:         syncStatus(value.SyncStatusNotSynced),
		ClearLastSyncError: true,
	}
}

func proposedDateSyncedWithGoogleEventUpdate(googleEventID string, syncedAt time.Time) domainProposedDate.ProposedDateUpdateOptions {
	return domainProposedDate.ProposedDateUpdateOptions{
		GoogleEventID:      &googleEventID,
		SyncStatus:         syncStatus(value.SyncStatusSynced),
		LastSyncedAt:       &syncedAt,
		ClearLastSyncError: true,
	}
}

func proposedDateSyncedUpdate(syncedAt time.Time) domainProposedDate.ProposedDateUpdateOptions {
	return domainProposedDate.ProposedDateUpdateOptions{
		SyncStatus:         syncStatus(value.SyncStatusSynced),
		LastSyncedAt:       &syncedAt,
		ClearLastSyncError: true,
	}
}

func proposedDateFailedUpdate(syncErr error) domainProposedDate.ProposedDateUpdateOptions {
	lastSyncError := syncErr.Error()
	return domainProposedDate.ProposedDateUpdateOptions{
		SyncStatus:    syncStatus(value.SyncStatusFailed),
		LastSyncError: &lastSyncError,
	}
}

func syncStatus(status value.SyncStatus) *value.SyncStatus {
	return &status
}

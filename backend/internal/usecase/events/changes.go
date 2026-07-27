package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func draftEventUpdate(status *value.EventStatus, syncExternally bool) EventMutation {
	if syncExternally {
		return eventPendingUpdate(status)
	}
	return EventMutation{
		Status:             status,
		SyncStatus:         syncStatus(value.SyncStatusNotSynced),
		ClearLastSyncError: true,
	}
}

func eventSyncPendingUpdate() EventMutation {
	return eventPendingUpdate(nil)
}

func confirmedEventPendingUpdate(confirmedDateID uuid.UUID) EventMutation {
	status := value.StatusConfirmed
	return EventMutation{
		Status:          &status,
		ConfirmedDateID: &confirmedDateID,
		SyncStatus:      syncStatus(value.SyncStatusPending),
	}
}

func draftProposedDateUpdate(start, end *time.Time, priority *int, status *value.ProposedDateStatus, syncExternally bool) ProposedDateMutation {
	update := ProposedDateMutation{
		StartTime: start,
		EndTime:   end,
		Priority:  priority,
		Status:    status,
	}
	if syncExternally {
		update.SyncStatus = syncStatus(value.SyncStatusPending)
		return update
	}
	update.SyncStatus = syncStatus(value.SyncStatusNotSynced)
	update.ClearLastSyncError = true
	return update
}

func confirmedProposedDateUpdate(start, end *time.Time, priority *int) ProposedDateMutation {
	status := value.ProposedDateStatusConfirmed
	return ProposedDateMutation{
		StartTime:  start,
		EndTime:    end,
		Priority:   priority,
		Status:     &status,
		SyncStatus: syncStatus(value.SyncStatusPending),
	}
}

func notSelectedProposedDateUpdate() ProposedDateMutation {
	status := value.ProposedDateStatusNotSelected
	return ProposedDateMutation{
		Status:     &status,
		SyncStatus: syncStatus(value.SyncStatusPending),
	}
}

func proposedDateSyncPendingUpdate() ProposedDateMutation {
	return ProposedDateMutation{
		SyncStatus: syncStatus(value.SyncStatusPending),
	}
}

func eventPendingUpdate(status *value.EventStatus) EventMutation {
	return EventMutation{
		Status:     status,
		SyncStatus: syncStatus(value.SyncStatusPending),
	}
}

func syncStatus(status value.SyncStatus) *value.SyncStatus {
	return &status
}

func googleEventIDOrEmpty(id *string) string {
	if id == nil {
		return ""
	}
	return *id
}

package events

import (
	"context"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
)

func (uc *Usecase) recordEventSyncFailure(ctx context.Context, store EventTxStore, eventID uuid.UUID, syncErr error) error {
	_, err := store.UpdateEvent(ctx, eventID, mergeEventChange(EventMutation{}, domainEvent.NewFailedEventChange(syncErr)))
	return err
}

func mergeEventChange(mutation EventMutation, change domainEvent.EventChange) EventMutation {
	mutation.Status = change.Status
	mutation.ConfirmedDateID = change.ConfirmedDateID
	mutation.GoogleEventID = change.GoogleEventID
	mutation.ConfirmedGoogleEventID = change.ConfirmedGoogleEventID
	if change.Sync.Status != "" {
		syncStatus := change.Sync.Status
		mutation.SyncStatus = &syncStatus
	}
	mutation.LastSyncedAt = change.Sync.LastSyncedAt
	mutation.ClearLastSyncedAt = change.Sync.ClearLastSyncedAt
	mutation.LastSyncError = change.Sync.LastSyncError
	mutation.ClearLastSyncError = change.Sync.ClearLastSyncError

	return mutation
}

func buildProposedDateMutation(change domainEvent.ProposedDateChange) ProposedDateMutation {
	mutation := ProposedDateMutation{
		GoogleEventID:      change.GoogleEventID,
		Start:              change.Start,
		End:                change.End,
		Priority:           change.Priority,
		Status:             change.Status,
		LastSyncedAt:       change.Sync.LastSyncedAt,
		ClearLastSyncedAt:  change.Sync.ClearLastSyncedAt,
		LastSyncError:      change.Sync.LastSyncError,
		ClearLastSyncError: change.Sync.ClearLastSyncError,
	}
	if change.Sync.Status != "" {
		syncStatus := change.Sync.Status
		mutation.SyncStatus = &syncStatus
	}

	return mutation
}

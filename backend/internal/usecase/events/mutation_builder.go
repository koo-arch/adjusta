package events

import (
	"context"
	"errors"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type EventMutation = domainEvent.EventUpdateOptions
type ProposedDateMutation = repoProposedDate.ProposedDateUpdateOptions

func (uc *Usecase) recordEventSyncFailure(ctx context.Context, repos EventTxRepositories, eventID uuid.UUID, syncErr error) error {
	_, err := repos.Event.Update(ctx, eventID, mergeEventChange(EventMutation{}, domainEvent.NewFailedEventChange(syncErr)))
	return err
}

func mergeEventChange(mutation EventMutation, change domainEvent.EventChange) EventMutation {
	mutation.Status = change.Status
	mutation.ConfirmedDateID = change.ConfirmedDateID
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

func toProposedDateCreateOptions(opt ProposedDateMutation) (repoProposedDate.ProposedDateCreateOptions, error) {
	if opt.StartTime == nil || opt.EndTime == nil || opt.Priority == nil {
		return repoProposedDate.ProposedDateCreateOptions{}, errors.New("start, end, and priority are required to create proposed date")
	}

	return repoProposedDate.ProposedDateCreateOptions{
		GoogleEventID: opt.GoogleEventID,
		StartTime:     *opt.StartTime,
		EndTime:       *opt.EndTime,
		Priority:      *opt.Priority,
		Status:        opt.Status,
		SyncStatus:    opt.SyncStatus,
		LastSyncedAt:  opt.LastSyncedAt,
		LastSyncError: opt.LastSyncError,
	}, nil
}

func toProposedDateCreateOptionsList(selectedDates []SelectedDate) []repoProposedDate.ProposedDateCreateOptions {
	opts := make([]repoProposedDate.ProposedDateCreateOptions, 0, len(selectedDates))
	for _, selectedDate := range selectedDates {
		status := value.ProposedDateStatusActive
		opts = append(opts, repoProposedDate.ProposedDateCreateOptions{
			StartTime: selectedDate.Start,
			EndTime:   selectedDate.End,
			Priority:  selectedDate.Priority,
			Status:    &status,
		})
	}
	return opts
}

func buildProposedDateMutation(change domainEvent.ProposedDateChange) ProposedDateMutation {
	mutation := ProposedDateMutation{
		GoogleEventID:      change.GoogleEventID,
		StartTime:          change.Start,
		EndTime:            change.End,
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

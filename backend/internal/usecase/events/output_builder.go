package events

import (
	"sort"

	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func buildProposedDateOutputs(storedDates []*ProposedDateRecord) []ProposedDateOutput {
	proposedDates := make([]ProposedDateOutput, 0, len(storedDates))
	for _, storedDate := range storedDates {
		proposedDates = append(proposedDates, ProposedDateOutput{
			ID:            &storedDate.ID,
			GoogleEventID: storedDate.GoogleEventID,
			Start:         &storedDate.StartTime,
			End:           &storedDate.EndTime,
			Priority:      storedDate.Priority,
			Status:        storedDate.Status,
			SyncStatus:    storedDate.SyncStatus,
			LastSyncedAt:  storedDate.LastSyncedAt,
			LastSyncError: storedDate.LastSyncError,
		})
	}

	sort.Slice(proposedDates, func(i, j int) bool {
		return proposedDates[i].Priority > proposedDates[j].Priority
	})

	return proposedDates
}

func buildEventDraftDetailOutput(storedEvent *EventRecord) (*EventDraftDetailOutput, error) {
	if storedEvent.ProposedDates == nil {
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return &EventDraftDetailOutput{
		ID:                     storedEvent.ID,
		Title:                  storedEvent.Title,
		Location:               storedEvent.Location,
		Description:            storedEvent.Description,
		Status:                 storedEvent.Status,
		SyncStatus:             storedEvent.SyncStatus,
		ConfirmedDateID:        &storedEvent.ConfirmedDateID,
		GoogleEventID:          domainEvent.ResolveGoogleEventID(storedEvent.ConfirmedGoogleEventID),
		ConfirmedGoogleEventID: storedEvent.ConfirmedGoogleEventID,
		LastSyncedAt:           storedEvent.LastSyncedAt,
		LastSyncError:          storedEvent.LastSyncError,
		ProposedDates:          buildProposedDateOutputs(storedEvent.ProposedDates),
	}, nil
}

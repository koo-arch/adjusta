package events

import (
	"errors"

	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type EventMutation = domainEvent.EventUpdateOptions
type ProposedDateMutation = repoProposedDate.ProposedDateUpdateOptions

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

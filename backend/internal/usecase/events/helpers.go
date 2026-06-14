package events

import (
	"context"
	"log"
	"sort"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) findPrimaryCalendar(ctx context.Context, finder PrimaryCalendarFinder, userID uuid.UUID, email string) (*CalendarRecord, error) {
	storedCalendar, err := finder.FindPrimaryCalendar(ctx, userID)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedCalendar, nil
}

func buildProposedDates(storedDates []*ProposedDateRecord) []appmodel.ProposedDate {
	proposedDates := make([]appmodel.ProposedDate, 0, len(storedDates))
	for _, storedDate := range storedDates {
		proposedDates = append(proposedDates, appmodel.ProposedDate{
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

func buildEventDraftDetail(storedEvent *EventRecord) (*appmodel.EventDraftDetail, error) {
	if storedEvent.ProposedDates == nil {
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return &appmodel.EventDraftDetail{
		ID:                     storedEvent.ID,
		Title:                  storedEvent.Title,
		Location:               storedEvent.Location,
		Description:            storedEvent.Description,
		Status:                 storedEvent.Status,
		SyncStatus:             storedEvent.SyncStatus,
		ConfirmedDateID:        &storedEvent.ConfirmedDateID,
		GoogleEventID:          eventGoogleEventID(storedEvent),
		ConfirmedGoogleEventID: storedEvent.ConfirmedGoogleEventID,
		LastSyncedAt:           storedEvent.LastSyncedAt,
		LastSyncError:          storedEvent.LastSyncError,
		Slug:                   storedEvent.Slug,
		ProposedDates:          buildProposedDates(storedEvent.ProposedDates),
	}, nil
}

func normalizeUsecaseError(err error, fallbackMessage string) error {
	return internalErrors.NormalizeAPIError(err, fallbackMessage)
}

func withPendingEventSync(mutation EventMutation) EventMutation {
	syncStatus := domainvalue.SyncStatusPending
	mutation.SyncStatus = &syncStatus
	return mutation
}

func withSyncedEventSync(mutation EventMutation) EventMutation {
	syncStatus := domainvalue.SyncStatusSynced
	mutation.SyncStatus = &syncStatus
	mutation.ClearLastSyncError = true
	return mutation
}

func withFailedEventSync(mutation EventMutation, syncErr error) EventMutation {
	syncStatus := domainvalue.SyncStatusFailed
	lastSyncError := syncErr.Error()
	mutation.SyncStatus = &syncStatus
	mutation.LastSyncError = &lastSyncError
	return mutation
}

func withPendingProposedDateSync(mutation ProposedDateMutation) ProposedDateMutation {
	syncStatus := domainvalue.SyncStatusPending
	mutation.SyncStatus = &syncStatus
	return mutation
}

func (uc *Usecase) markEventSyncFailed(ctx context.Context, store EventTxStore, eventID uuid.UUID, syncErr error) error {
	_, err := store.UpdateEvent(ctx, eventID, withFailedEventSync(EventMutation{}, syncErr))
	return err
}

func eventGoogleEventID(storedEvent *EventRecord) string {
	if storedEvent == nil {
		return ""
	}
	if storedEvent.ConfirmedGoogleEventID != nil && *storedEvent.ConfirmedGoogleEventID != "" {
		return *storedEvent.ConfirmedGoogleEventID
	}
	return storedEvent.GoogleEventID
}

func toDomainDraftProposedDate(date appmodel.ProposedDate) (domainEvent.DraftProposedDate, error) {
	if date.Start == nil || date.End == nil {
		return domainEvent.DraftProposedDate{}, internalErrors.NewBadRequestError("候補日程が不正です")
	}

	return domainEvent.DraftProposedDate{
		ID:       date.ID,
		Start:    *date.Start,
		End:      *date.End,
		Priority: date.Priority,
	}, nil
}

func toDomainDraftProposedDates(dates []appmodel.ProposedDate) ([]domainEvent.DraftProposedDate, error) {
	converted := make([]domainEvent.DraftProposedDate, 0, len(dates))
	for _, date := range dates {
		convertedDate, err := toDomainDraftProposedDate(date)
		if err != nil {
			return nil, err
		}
		converted = append(converted, convertedDate)
	}
	return converted, nil
}

func toDomainConfirmationDate(date appmodel.ConfirmDate) (domainEvent.DraftProposedDate, error) {
	if date.Start == nil || date.End == nil {
		return domainEvent.DraftProposedDate{}, internalErrors.NewBadRequestError("確定候補日程が不正です")
	}

	return domainEvent.DraftProposedDate{
		ID:       date.ID,
		Start:    *date.Start,
		End:      *date.End,
		Priority: date.Priority,
	}, nil
}

func normalizeSelectedDatesPriorities(dates []appmodel.SelectedDate) []appmodel.SelectedDate {
	normalized := make([]appmodel.SelectedDate, 0, len(dates))
	for i, date := range dates {
		date.Priority = domainEvent.PriorityValueForOrder(i, len(dates))
		normalized = append(normalized, date)
	}
	return normalized
}

func toDomainExistingProposedDates(dates []*ProposedDateRecord) []domainEvent.ExistingProposedDate {
	converted := make([]domainEvent.ExistingProposedDate, 0, len(dates))
	for _, date := range dates {
		converted = append(converted, domainEvent.ExistingProposedDate{
			ID:       date.ID,
			Start:    date.StartTime,
			End:      date.EndTime,
			Priority: date.Priority,
		})
	}
	return converted
}

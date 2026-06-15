package events

import (
	"context"
	"log"
	"sort"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) loadPrimaryCalendar(ctx context.Context, finder PrimaryCalendarFinder, userID uuid.UUID, email string) (*CalendarRecord, error) {
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

func buildAppProposedDates(storedDates []*ProposedDateRecord) []appmodel.ProposedDate {
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

func buildAppEventDraftDetail(storedEvent *EventRecord) (*appmodel.EventDraftDetail, error) {
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
		GoogleEventID:          domainEvent.ResolveGoogleEventID(storedEvent.ConfirmedGoogleEventID, storedEvent.GoogleEventID),
		ConfirmedGoogleEventID: storedEvent.ConfirmedGoogleEventID,
		LastSyncedAt:           storedEvent.LastSyncedAt,
		LastSyncError:          storedEvent.LastSyncError,
		Slug:                   storedEvent.Slug,
		ProposedDates:          buildAppProposedDates(storedEvent.ProposedDates),
	}, nil
}

func mapUsecaseError(err error, fallbackMessage string) error {
	return internalErrors.NormalizeAPIError(err, fallbackMessage)
}

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

func toDomainDraftDate(date appmodel.ProposedDate) (domainEvent.DraftProposedDate, error) {
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

func toDomainDraftDateList(dates []appmodel.ProposedDate) ([]domainEvent.DraftProposedDate, error) {
	converted := make([]domainEvent.DraftProposedDate, 0, len(dates))
	for _, date := range dates {
		convertedDate, err := toDomainDraftDate(date)
		if err != nil {
			return nil, err
		}
		converted = append(converted, convertedDate)
	}
	return converted, nil
}

func toDomainConfirmationDraftDate(date appmodel.ConfirmDate) (domainEvent.DraftProposedDate, error) {
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

func assignSelectedDatePriorities(dates []appmodel.SelectedDate) []appmodel.SelectedDate {
	assigned := make([]appmodel.SelectedDate, 0, len(dates))
	for i, date := range dates {
		date.Priority = domainEvent.PriorityForOrder(i, len(dates))
		assigned = append(assigned, date)
	}
	return assigned
}

func toDomainExistingDateList(dates []*ProposedDateRecord) []domainEvent.ExistingProposedDate {
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

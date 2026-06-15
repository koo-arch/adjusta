package events

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) UpdateDraftedEvents(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *appmodel.EventDraftUpdate) error {
	var committedErr error

	err := uc.tx.Do(ctx, func(store EventTxStore) error {
		if _, err := uc.loadPrimaryCalendar(ctx, store, userID, email); err != nil {
			return err
		}

		storedEvent, err := store.FindEventBySlug(ctx, userID, slug, false)
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		eventOptions := EventMutation{
			Title:       &eventReq.Title,
			Location:    &eventReq.Location,
			Description: &eventReq.Description,
		}
		if eventReq.Status != domainvalue.StatusConfirmed {
			eventOptions = mergeEventChange(eventOptions, domainEvent.NewPendingEventChange(&eventReq.Status))
		}
		storedEvent, err = store.UpdateEvent(ctx, storedEvent.ID, eventOptions)
		if err != nil {
			log.Printf("failed to update event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		existingDates, err := store.ListProposedDatesByEvent(ctx, storedEvent.ID)
		if err != nil {
			log.Printf("failed to get proposed dates for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		if eventReq.Status == domainvalue.StatusConfirmed {
			storedCalendar, err := store.ReadCalendar(ctx, storedEvent.PrimaryCalendarID)
			if err != nil {
				log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
				return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
			}

			confirmDate := appmodel.ConfirmDate{
				ID:       eventReq.ProposedDates[0].ID,
				Start:    eventReq.ProposedDates[0].Start,
				End:      eventReq.ProposedDates[0].End,
				Priority: eventReq.ProposedDates[0].Priority,
			}
			confirmEvent := appmodel.ConfirmEvent{
				ConfirmDate: confirmDate,
			}

			googleEventID, err := uc.handleGoogleEvent(ctx, userID, storedCalendar.GoogleCalendarID, storedEvent, &confirmEvent)
			if err != nil {
				log.Printf("failed to handle google event for account: %s, error: %v", email, err)
				if syncErr := uc.recordEventSyncFailure(ctx, store, storedEvent.ID, err); syncErr != nil {
					log.Printf("failed to mark sync failure for account: %s, error: %v", email, syncErr)
					return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
				}
				committedErr = internalErrors.NormalizeAPIError(err, "Googleカレンダー更新時にエラーが発生しました")
				return nil
			}

			if err := uc.confirmEventDate(ctx, store, googleEventID, &confirmEvent, storedEvent); err != nil {
				log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
				return mapUsecaseError(err, internalErrors.InternalErrorMessage)
			}
		}

		if err := uc.updateProposedDates(ctx, store, eventReq, storedEvent, existingDates); err != nil {
			return mapUsecaseError(err, internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running update drafted event transaction: %v", err)
		return mapUsecaseError(err, internalErrors.InternalErrorMessage)
	}
	if committedErr != nil {
		return committedErr
	}

	return nil
}

func (uc *Usecase) updateProposedDates(ctx context.Context, store EventTxStore, eventReq *appmodel.EventDraftUpdate, storedEvent *EventRecord, existingDates []*ProposedDateRecord) error {
	requestedDates, err := toDomainDraftDateList(eventReq.ProposedDates)
	if err != nil {
		return err
	}

	changeSet := domainEvent.PlanProposedDateChanges(requestedDates, toDomainExistingDateList(existingDates))

	for _, date := range changeSet.Updates {
		dateOptions := buildProposedDateMutation(domainEvent.NewPendingProposedDateChange(&date.Start, &date.End, &date.Priority, nil))
		if _, err := store.UpdateProposedDate(ctx, date.ID, dateOptions); err != nil {
			return fmt.Errorf("failed to update proposed date for account: %s, error: %w", date.ID, err)
		}
	}

	for _, dateID := range changeSet.Deletes {
		if _, err := store.UpdateProposedDate(ctx, dateID, buildProposedDateMutation(domainEvent.NewPendingProposedDateChange(nil, nil, nil, nil))); err != nil {
			return fmt.Errorf("failed to mark proposed date sync pending for account: %s, error: %w", dateID, err)
		}
		if err := store.DeleteProposedDate(ctx, dateID); err != nil {
			return fmt.Errorf("failed to delete proposed date for account: %s, error: %w", dateID, err)
		}
	}

	for _, date := range changeSet.Creates {
		dateOptions := buildProposedDateMutation(domainEvent.NewPendingProposedDateChange(&date.Start, &date.End, &date.Priority, nil))
		if _, err := store.CreateProposedDate(ctx, dateOptions, storedEvent.ID); err != nil {
			return fmt.Errorf("failed to create proposed date for event: %s, error: %w", storedEvent.ID, err)
		}
	}

	return nil
}

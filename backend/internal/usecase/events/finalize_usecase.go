package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FinalizeProposedDate(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, email string, confirmation ConfirmationRequest) error {
	var committedErr error

	err := uc.tx.Do(ctx, func(store EventTxStore) error {
		storedEvent, err := store.FindEventByID(ctx, userID, eventID, false)
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		storedCalendar, err := store.ReadCalendar(ctx, storedEvent.PrimaryCalendarID)
		if err != nil {
			log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		googleEventID, err := uc.handleGoogleEvent(ctx, userID, storedCalendar.GoogleCalendarID, storedEvent, confirmation)
		if err != nil {
			log.Printf("failed to handle google event for account: %s, error: %v", email, err)
			if syncErr := uc.recordEventSyncFailure(ctx, store, storedEvent.ID, err); syncErr != nil {
				log.Printf("failed to mark sync failure for account: %s, error: %v", email, syncErr)
				return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
			}
			committedErr = internalErrors.NormalizeAPIError(err, "サーバーでエラーが発生しました")
			return nil
		}

		if err := uc.confirmEventDate(ctx, store, googleEventID, confirmation, storedEvent); err != nil {
			log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
			return mapUsecaseError(err, internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running finalize proposed date transaction: %v", err)
		return mapUsecaseError(err, internalErrors.InternalErrorMessage)
	}
	if committedErr != nil {
		return committedErr
	}

	return nil
}

func (uc *Usecase) handleGoogleEvent(ctx context.Context, userID uuid.UUID, calendarID string, storedEvent *EventRecord, confirmation ConfirmationRequest) (*string, error) {
	existingGoogleEventID := domainEvent.ResolveReusableGoogleEventID(
		confirmation.ID,
		storedEvent.ConfirmedGoogleEventID,
		confirmation.GoogleEventID,
	)

	googleEventID, err := uc.googleCalendar.UpsertEvent(
		ctx,
		userID,
		calendarID,
		existingGoogleEventID,
		storedEvent.Title,
		storedEvent.Location,
		storedEvent.Description,
		*confirmation.Start,
		*confirmation.End,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert google event: %w", err)
	}

	return &googleEventID, nil
}

func (uc *Usecase) confirmEventDate(ctx context.Context, store EventTxStore, googleEventID *string, confirmation ConfirmationRequest, storedEvent *EventRecord) error {
	confirmDate, err := toDomainConfirmationRequest(confirmation)
	if err != nil {
		return err
	}

	existingDates, err := store.ListProposedDatesByEvent(ctx, storedEvent.ID)
	if err != nil {
		return fmt.Errorf("failed to list proposed dates error: %w", err)
	}

	changeSet, err := domainEvent.PlanConfirmationChanges(confirmDate, toDomainExistingDateList(existingDates))
	if err != nil {
		return internalErrors.NewBadRequestError("確定候補日程が不正です")
	}

	confirmDateID := confirmation.ID
	if confirmation.ID == nil {
		confirmedStatus := domainvalue.ProposedDateStatusConfirmed
		dateOptions := buildProposedDateMutation(domainEvent.NewPendingProposedDateChange(
			&changeSet.Create.Start,
			&changeSet.Create.End,
			&changeSet.Create.Priority,
			&confirmedStatus,
		))

		storedDate, err := store.CreateProposedDate(ctx, dateOptions, storedEvent.ID)
		if err != nil {
			return fmt.Errorf("failed to create proposed date error: %w", err)
		}
		confirmDateID = &storedDate.ID
	}

	if confirmation.ID != nil {
		confirmedStatus := domainvalue.ProposedDateStatusConfirmed
		dateOptions := buildProposedDateMutation(domainEvent.NewPendingProposedDateChange(
			nil,
			nil,
			&changeSet.Update.Priority,
			&confirmedStatus,
		))

		if _, err := store.UpdateProposedDate(ctx, *confirmation.ID, dateOptions); err != nil {
			return fmt.Errorf("failed to update proposed date error: %w", err)
		}
	}

	if err := uc.markUnselectedProposedDates(ctx, store, changeSet.MarkNotSelected); err != nil {
		return fmt.Errorf("failed to update proposed date statuses: %w", err)
	}

	eventOptions := mergeEventChange(EventMutation{}, domainEvent.NewSyncedEventChange(changeSet.Status, *confirmDateID, *googleEventID, time.Now()))
	if _, err := store.UpdateEvent(ctx, storedEvent.ID, eventOptions); err != nil {
		return fmt.Errorf("failed to update event status error: %w", err)
	}

	return nil
}

func (uc *Usecase) markUnselectedProposedDates(ctx context.Context, store EventTxStore, proposedDateIDs []uuid.UUID) error {
	notSelected := domainvalue.ProposedDateStatusNotSelected
	for _, proposedDateID := range proposedDateIDs {
		if _, err := store.UpdateProposedDate(ctx, proposedDateID, buildProposedDateMutation(domainEvent.NewPendingProposedDateChange(nil, nil, nil, &notSelected))); err != nil {
			return err
		}
	}

	return nil
}

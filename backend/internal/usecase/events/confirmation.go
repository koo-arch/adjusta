package events

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FinalizeProposedDate(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, email string, confirmation ConfirmationRequest) error {
	var outboxMessageID uuid.UUID
	err := uc.tx.DoEvent(ctx, func(repos EventTxRepositories) error {
		storedEvent, err := repos.Event.FindByIDAndUser(ctx, userID, eventID, domainEvent.EventReadOptions{
			WithProposedDates: false,
		})
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		if err := uc.confirmEventDate(ctx, repos, confirmation, storedEvent); err != nil {
			log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
			return mapUsecaseError(err, internalErrors.InternalErrorMessage)
		}
		outboxMessageID, err = enqueueGoogleCalendarSync(ctx, repos, storedEvent.ID)
		if err != nil {
			return fmt.Errorf("failed to enqueue confirmed event sync: %w", err)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running finalize proposed date transaction: %v", err)
		return mapUsecaseError(err, internalErrors.InternalErrorMessage)
	}
	uc.dispatchOutboxMessage(ctx, outboxMessageID)
	return nil
}

func (uc *Usecase) confirmEventDate(ctx context.Context, repos EventTxRepositories, confirmation ConfirmationRequest, storedEvent *domainEvent.Event) error {
	confirmDate, err := toDomainConfirmationRequest(confirmation)
	if err != nil {
		return err
	}

	existingDates, err := repos.ProposedDate.FilterByEventID(ctx, storedEvent.ID)
	if err != nil {
		return fmt.Errorf("failed to list proposed dates error: %w", err)
	}

	changeSet, err := domainEvent.PlanConfirmationChanges(confirmDate, toDomainExistingDateList(existingDates))
	if err != nil {
		return internalErrors.NewBadRequestError("確定候補日程が不正です")
	}

	confirmDateID := confirmation.ID
	if confirmation.ID == nil {
		dateOptions := confirmedProposedDateUpdate(
			&changeSet.Create.Start,
			&changeSet.Create.End,
			&changeSet.Create.Priority,
		)

		createOptions, err := toProposedDateCreateOptions(dateOptions)
		if err != nil {
			return err
		}

		storedDate, err := repos.ProposedDate.Create(ctx, createOptions, storedEvent.ID)
		if err != nil {
			return fmt.Errorf("failed to create proposed date error: %w", err)
		}
		confirmDateID = &storedDate.ID
	}

	if confirmation.ID != nil {
		dateOptions := confirmedProposedDateUpdate(
			nil,
			nil,
			&changeSet.Update.Priority,
		)

		if _, err := repos.ProposedDate.Update(ctx, *confirmation.ID, dateOptions); err != nil {
			return fmt.Errorf("failed to update proposed date error: %w", err)
		}
	}

	if err := uc.markUnselectedProposedDates(ctx, repos, changeSet.MarkNotSelected); err != nil {
		return fmt.Errorf("failed to update proposed date statuses: %w", err)
	}

	eventOptions := confirmedEventPendingUpdate(*confirmDateID)
	if _, err := repos.Event.Update(ctx, storedEvent.ID, eventOptions); err != nil {
		return fmt.Errorf("failed to update event status error: %w", err)
	}

	return nil
}

func (uc *Usecase) markUnselectedProposedDates(ctx context.Context, repos EventTxRepositories, proposedDateIDs []uuid.UUID) error {
	for _, proposedDateID := range proposedDateIDs {
		if _, err := repos.ProposedDate.Update(ctx, proposedDateID, notSelectedProposedDateUpdate()); err != nil {
			return err
		}
	}

	return nil
}

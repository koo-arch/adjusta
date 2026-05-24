package events

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FinalizeProposedDate(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *models.ConfirmEvent) error {
	err := uc.tx.Do(ctx, func(store EventTxStore) error {
		storedEvent, err := store.FindEventBySlug(ctx, userID, slug, false)
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		googleEventID, err := uc.handleGoogleEvent(ctx, userID, storedEvent, eventReq)
		if err != nil {
			log.Printf("failed to handle google event for account: %s, error: %v", email, err)
			return internalErrors.NormalizeAPIError(err, "サーバーでエラーが発生しました")
		}

		if err := uc.confirmEventDate(ctx, store, googleEventID, eventReq, storedEvent); err != nil {
			log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running finalize proposed date transaction: %v", err)
		return normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return nil
}

func (uc *Usecase) handleGoogleEvent(ctx context.Context, userID uuid.UUID, storedEvent *models.StoredEvent, eventReq *models.ConfirmEvent) (*string, error) {
	var existingGoogleEventID *string
	if eventReq.ConfirmDate.ID != nil && storedEvent.GoogleEventID != "" {
		existingGoogleEventID = &storedEvent.GoogleEventID
	}

	googleEventID, err := uc.googleCalendar.UpsertEvent(
		ctx,
		userID,
		existingGoogleEventID,
		storedEvent.Summary,
		storedEvent.Location,
		storedEvent.Description,
		*eventReq.ConfirmDate.Start,
		*eventReq.ConfirmDate.End,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert google event: %w", err)
	}

	return &googleEventID, nil
}

func (uc *Usecase) confirmEventDate(ctx context.Context, store EventTxStore, googleEventID *string, eventReq *models.ConfirmEvent, storedEvent *models.StoredEvent) error {
	priority := 1
	dateOptions := ProposedDateMutation{
		Priority: &priority,
	}

	confirmDateID := eventReq.ConfirmDate.ID
	if eventReq.ConfirmDate.ID == nil {
		dateOptions.Start = eventReq.ConfirmDate.Start
		dateOptions.End = eventReq.ConfirmDate.End

		storedDate, err := store.CreateProposedDate(ctx, dateOptions, storedEvent.ID)
		if err != nil {
			return fmt.Errorf("failed to create proposed date error: %w", err)
		}
		confirmDateID = &storedDate.ID

		if err := store.DecrementPriorityExceptID(ctx, storedEvent.ID, storedDate.ID); err != nil {
			return fmt.Errorf("failed to decrement priority error: %w", err)
		}
	}

	if eventReq.ConfirmDate.ID != nil {
		zero := 0
		dateOptions.Priority = &zero

		if _, err := store.UpdateProposedDate(ctx, *eventReq.ConfirmDate.ID, dateOptions); err != nil {
			return fmt.Errorf("failed to update proposed date error: %w", err)
		}

		if err := store.ReorderPriority(ctx, storedEvent.ID); err != nil {
			return fmt.Errorf("failed to reorder priority error: %w", err)
		}
	}

	status := models.StatusConfirmed
	eventOptions := EventMutation{
		Status:          &status,
		ConfirmedDateID: confirmDateID,
		GoogleEventID:   googleEventID,
	}
	if _, err := store.UpdateEvent(ctx, storedEvent.ID, eventOptions); err != nil {
		return fmt.Errorf("failed to update event status error: %w", err)
	}

	return nil
}

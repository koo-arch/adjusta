package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func (uc *Usecase) DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) error {
	err := uc.tx.Do(ctx, func(store EventTxStore) error {
		if _, err := uc.loadPrimaryCalendar(ctx, store, userID, email); err != nil {
			return err
		}

		if _, err := store.UpdateEvent(ctx, eventID, mergeEventChange(EventMutation{}, domainEvent.NewPendingEventChange(nil))); err != nil {
			log.Printf("failed to mark event sync pending for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError("イベント削除時にエラーが発生しました")
		}

		if err := store.SoftDeleteEvent(ctx, eventID); err != nil {
			log.Printf("failed to delete event for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError("イベント削除時にエラーが発生しました")
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running delete drafted event transaction: %v", err)
		return mapUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return nil
}

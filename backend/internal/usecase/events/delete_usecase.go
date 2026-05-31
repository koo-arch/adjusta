package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func (uc *Usecase) DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *appmodel.EventDraftDetail) error {
	err := uc.tx.Do(ctx, func(store EventTxStore) error {
		if _, err := uc.findPrimaryCalendar(ctx, store, userID, email); err != nil {
			return err
		}

		if err := store.SoftDeleteEvent(ctx, eventReq.ID); err != nil {
			log.Printf("failed to delete event for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError("イベント削除時にエラーが発生しました")
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running delete drafted event transaction: %v", err)
		return normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return nil
}

package events

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	internalRepo "github.com/koo-arch/adjusta-backend/internal/repo"
)

func (uc *Usecase) DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftDetail) error {
	err := uc.uow.Do(ctx, func(repos internalRepo.Repositories) error {
		if _, err := uc.findPrimaryCalendar(ctx, repos, userID, email); err != nil {
			return err
		}

		if err := repos.Event.SoftDelete(ctx, eventReq.ID); err != nil {
			log.Printf("failed to delete event for account: %s, error: %v", email, err)
			return internalErrors.NewAPIError(http.StatusInternalServerError, "イベント削除時にエラーが発生しました")
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running delete drafted event transaction: %v", err)
		return normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return nil
}

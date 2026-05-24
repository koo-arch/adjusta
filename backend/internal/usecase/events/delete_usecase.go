package events

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

func (uc *Usecase) DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftDetail) error {
	tx, err := uc.client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, "エラーが発生しました")
	}

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}

	_, err = uc.calendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return internalErrors.NewAPIError(http.StatusInternalServerError, "カレンダー取得時にエラーが発生しました")
	}

	if err := uc.eventRepo.SoftDelete(ctx, tx, eventReq.ID); err != nil {
		log.Printf("failed to delete event for account: %s, error: %v", email, err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, "イベント削除時にエラーが発生しました")
	}

	return nil
}

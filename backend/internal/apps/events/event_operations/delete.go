package event_operations

import (
	"context"
	"net/http"
	"log"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

type EventDeleteManager struct {
	event *events.EventManager
}

func NewEventDeleteManager(event *events.EventManager) *EventDeleteManager {
	return &EventDeleteManager{
		event: event,
	}
}

func (edm *EventDeleteManager) DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftDetail) error {
	tx, err := edm.event.Client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, "エラーが発生しました")
	}

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}

	_, err = edm.event.CalendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return internalErrors.NewAPIError(http.StatusInternalServerError, "カレンダー取得時にエラーが発生しました")
	}


	// データベースからイベントを削除
	err = edm.event.EventRepo.SoftDelete(ctx, tx, eventReq.ID)
	if err != nil {
		log.Printf("failed to delete event for account: %s, error: %v", email, err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, "イベント削除時にエラーが発生しました")
	}

	return nil
}

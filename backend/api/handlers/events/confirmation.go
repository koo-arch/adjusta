package events

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/dto"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/validation"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

func (eh *Handler) EventFinalizeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventID, ok := parseEventIDParam(c)
		if !ok {
			return
		}

		var confirmEvent dto.ConfirmEvent
		if err := c.ShouldBindJSON(&confirmEvent); err != nil {
			respond.BadRequest(c, "リクエストのデータ形式が不正です")
			return
		}

		if err := validation.FinalizeValidation(&confirmEvent); err != nil {
			log.Printf("failed to validate confirm event: %v", err)
			respond.Error(c, err, "イベントの確定に失敗しました")
			return
		}

		eventUsecase := eh.confirmationUsecase

		confirmation := usecaseEvents.ConfirmationRequest{
			ID:            confirmEvent.ConfirmDate.ID,
			GoogleEventID: confirmEvent.ConfirmDate.GoogleEventID,
			Start:         confirmEvent.ConfirmDate.Start,
			End:           confirmEvent.ConfirmDate.End,
			Priority:      confirmEvent.ConfirmDate.Priority,
		}

		err = eventUsecase.FinalizeProposedDate(ctx, userid, eventID, email, confirmation)
		if err != nil {
			log.Printf("failed to finalize event: %v", err)
			respond.Error(c, err, "イベントの確定に失敗しました")
			return
		}

		respond.OKMessage(c, "success")
	}
}

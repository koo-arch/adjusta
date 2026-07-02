package events

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
)

func (eh *Handler) FetchEventDraftDetailHandler() gin.HandlerFunc {
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

		eventUsecase := eh.detailUsecase

		draftedEvent, err := eventUsecase.FetchDraftedEventDetail(ctx, userid, email, eventID)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "イベント詳細の取得に失敗しました")
			return
		}

		respond.OK(c, toEventDraftDetailResponse(draftedEvent))
	}
}

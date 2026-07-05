package events

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
)

func (eh *Handler) FetchUpcomingEventsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := eh.agendaUsecase

		daysBefore := 3
		upcomingEvents, err := eventUsecase.FetchUpcomingEvents(ctx, userid, email, daysBefore)
		if err != nil {
			log.Printf("failed to fetch upcoming events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, toUpcomingEventResponses(upcomingEvents))
	}
}

func (eh *Handler) FetchNeedsActionDraftsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := eh.agendaUsecase

		daysBefore := 3
		needsActionDrafts, err := eventUsecase.FetchNeedsActionDrafts(ctx, userid, email, daysBefore)
		if err != nil {
			log.Printf("failed to fetch needs action events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, toNeedsActionDraftResponses(needsActionDrafts))
	}
}

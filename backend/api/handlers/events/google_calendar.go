package events

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/internal/errors"
)

func (eh *Handler) FetchEventListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, "ユーザー情報確認時にエラーが発生しました。")
			return
		}

		events, err := eh.googleCalendarUsecase.FetchAllGoogleEvents(ctx, userid, email)
		if err, ok := err.(*errors.APIError); ok && err.Kind == errors.KindPartial {
			respond.Partial(c, gin.H{
				"events":  toGoogleEventResponses(events),
				"warning": err.Details,
			})
			return
		}
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "Googleカレンダーのイベント取得に失敗しました")
			return
		}

		respond.OK(c, gin.H{
			"events": toGoogleEventResponses(events),
		})
	}
}

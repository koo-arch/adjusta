package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
)

type CalendarMiddleware struct {
	calendarSyncUsecase CalendarSyncUsecase
}

func NewCalendarMiddleware(calendarSyncUsecase CalendarSyncUsecase) *CalendarMiddleware {
	return &CalendarMiddleware{
		calendarSyncUsecase: calendarSyncUsecase,
	}
}

func (cm *CalendarMiddleware) SyncGoogleCalendars() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			log.Printf("failed to extract user info for account: %s, %v", email, err)
			respond.Error(c, err, "ユーザー情報確認時にエラーが発生しました。")
			return
		}

		calendarUsecase := cm.calendarSyncUsecase
		calendarList, err := calendarUsecase.SyncGoogleCalendars(ctx, userid, email)
		if err != nil {
			log.Printf("failed to register calendar list for account: %s, error: %v", email, err)
			respond.Error(c, err, "Googleカレンダーの同期に失敗しました")
			return
		}

		c.Set("calendarList", calendarList)
		c.Next()
	}
}

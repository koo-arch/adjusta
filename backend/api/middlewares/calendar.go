package middlewares

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/utils"
)

type CalendarMiddleware struct {
	middleware *Middleware
}

func NewCalendarMiddleware(middleware *Middleware) *CalendarMiddleware {
	return &CalendarMiddleware{
		middleware: middleware,
	}
}

func (cm *CalendarMiddleware) SyncGoogleCalendars() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			log.Printf("failed to extract user info for account: %s, %v", email, err)
			utils.HandleAPIError(c, err, "ユーザー情報確認時にエラーが発生しました。")
			return
		}

		// キャッシュにある場合はそれを使う
		cache := cm.middleware.Server.Cache
		cacheKey := fmt.Sprintf("calendars:%s", userid)
		if cacheCalendar, found := cache.CalendarCache.Get(cacheKey); found {
			c.Set("calendarList", cacheCalendar)
			c.Next()
			c.Abort()
			return
		}

		calendarUsecase := cm.middleware.Server.CalendarSyncUsecase
		calendarList, err := calendarUsecase.SyncGoogleCalendars(ctx, userid, email)
		if err != nil {
			log.Printf("failed to register calendar list for account: %s, error: %v", email, err)
			utils.HandleAPIError(c, err, "Googleカレンダーの同期に失敗しました")
			return
		}

		cache.CalendarCache.Set(cacheKey, calendarList, 5*time.Hour)
		c.Set("calendarList", calendarList)
		c.Next()
	}
}

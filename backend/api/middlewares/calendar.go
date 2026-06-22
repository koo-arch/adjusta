package middlewares

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
	infraCache "github.com/koo-arch/adjusta-backend/internal/infrastructure/cache"
)

type CalendarMiddleware struct {
	cache               *infraCache.Cache
	calendarSyncService CalendarSyncService
}

func NewCalendarMiddleware(cache *infraCache.Cache, calendarSyncService CalendarSyncService) *CalendarMiddleware {
	return &CalendarMiddleware{
		cache:               cache,
		calendarSyncService: calendarSyncService,
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

		// キャッシュにある場合はそれを使う
		cache := cm.cache
		cacheKey := fmt.Sprintf("calendars:%s", userid)
		if cacheCalendar, found := cache.CalendarCache.Get(cacheKey); found {
			c.Set("calendarList", cacheCalendar)
			c.Next()
			c.Abort()
			return
		}

		calendarUsecase := cm.calendarSyncService
		calendarList, err := calendarUsecase.SyncGoogleCalendars(ctx, userid, email)
		if err != nil {
			log.Printf("failed to register calendar list for account: %s, error: %v", email, err)
			respond.Error(c, err, "Googleカレンダーの同期に失敗しました")
			return
		}

		cache.CalendarCache.Set(cacheKey, calendarList, 5*time.Hour)
		c.Set("calendarList", calendarList)
		c.Next()
	}
}

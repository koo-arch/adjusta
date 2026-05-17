package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/ent"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"github.com/koo-arch/adjusta-backend/utils"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
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

		userRepo := cm.middleware.Server.UserRepo
		entUser, err := userRepo.Read(ctx, nil, userid, user.UserQueryOptions{})
		if err != nil {
			log.Printf("failed to get user info for account: %s, %v", email, err)
			if ent.IsNotFound(err) {
				utils.HandleAPIError(c, err, "ユーザー情報が見つかりませんでした")
			}
			utils.HandleAPIError(c, err, "ユーザー情報取得時にエラーが発生しました")
			return
		}

		// キャッシュにある場合はそれを使う
		cache := cm.middleware.Server.Cache
		cacheKey := fmt.Sprintf("calendars:%s", userid)
		if cacheCalendar, found := cache.CalendarCache.Get(cacheKey); found {
			println("use cache")
			c.Set("calendarList", cacheCalendar)
			c.Next()
			c.Abort()
			return
		}

		calendarList, err := cm.fetchAndCacheCalendars(ctx, entUser)
		if err != nil {
			log.Printf("failed to register calendar list for account: %s, error: %v",email,  err)
			utils.HandleAPIError(c, err, "Googleカレンダーの同期に失敗しました")
			return
		}

		c.Set("calendarList", calendarList)
		c.Next()
	}
}

func(cm *CalendarMiddleware) fetchAndCacheCalendars(ctx context.Context, entUser *ent.User) ([]*customCalendar.CalendarList, error) {

	authManager := cm.middleware.Server.AuthManager
	token, err := authManager.VerifyOAuthToken(ctx, entUser.ID)
	if err != nil {
		log.Printf("failed to verify token for account: %s, error: %v", entUser.Email, err)
		apiErr := utils.GetAPIError(err, "OAuthトークンの認証に失敗しました")
		return nil, apiErr
	}

	calendarService, err := customCalendar.NewCalendar(ctx, token)
	if err != nil {
		log.Printf("failed to create calendar service for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "Googleカレンダー接続に失敗しました")
	}

	calendars, err := calendarService.FetchCalendarList()
	if err != nil {
		log.Printf("failed to fetch calendars for account: %s, error: %v", entUser.Email, err)
		apiErr := utils.HandleGoogleAPIError(err)
		return nil, apiErr
	}

	if err := cm.syncCalendar(ctx, calendars, entUser); err != nil {
		log.Printf("failed to sync calendars for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	// キャッシュに保存
	calendarCache := cm.middleware.Server.Cache.CalendarCache
	cacheKey := fmt.Sprintf("calendars:%s", entUser.ID.String())
	calendarCache.Set(cacheKey, calendars, 5*time.Hour)

	return calendars, nil
}

func(cm *CalendarMiddleware) syncCalendar(ctx context.Context, calendars []*customCalendar.CalendarList, entUser *ent.User) error {
	// トランザクションを開始
	tx, err := cm.middleware.Server.Client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	// トランザクションをデファーで処理
	defer transaction.HandleTransaction(tx, &err)

	calendarRepo := cm.middleware.Server.CalendarRepo
	googleCalendarRepo := cm.middleware.Server.GoogleCalendarRepo

	// 同期対象を集める
	incoming := make(map[string]struct{}, len(calendars))
	// カレンダーがDBに存在するか確認
	for _, cal := range calendars {
		// カレンダーIDをキーにしてマップに追加
		incoming[cal.CalendarID] = struct{}{}
		// カレンダーがDBに存在するか確認
		repoCalOptions := repoCalendar.CalendarQueryOptions{
			WithGoogleCalendarInfo: true,
			GoogleCalendarID: &cal.CalendarID,
		}
		entCalendar, err := calendarRepo.FindByFields(ctx, tx, entUser.ID, repoCalOptions)
		if err != nil {
			if !ent.IsNotFound(err) {
				return fmt.Errorf("failed to find calendar: %w", err)
			}
			
			// カレンダーが存在しない場合は新規作成
			entCalendar, err = calendarRepo.Create(ctx, tx, entUser, nil)
			if err != nil {
				return fmt.Errorf("failed to create calendar: %w", err)
			}
		}

		// GoogleCalendarInfoにカレンダーが存在するか確認
		gCalOptions := googlecalendarinfo.GoogleCalendarInfoQueryOptions{
			GoogleCalendarID : &cal.CalendarID,
		}
		entGoogleCalendar, err := googleCalendarRepo.FindByFields(ctx, tx, gCalOptions)
		if err != nil {
			if !ent.IsNotFound(err) {
				return fmt.Errorf("failed to find google calendar info: %w", err)
			}

			// 存在しない場合は新規作成
			gCalOptions := googlecalendarinfo.GoogleCalendarInfoQueryOptions{
				GoogleCalendarID: &cal.CalendarID,
				Summary: &cal.Summary,
				IsPrimary: &cal.Primary,
			}
			_, err := googleCalendarRepo.Create(ctx, tx, gCalOptions, entCalendar)
			if err != nil {
				return fmt.Errorf("failed to create google calendar info: %w", err)
			}
		}

		// カレンダーが存在する場合は関連付け
		if entGoogleCalendar != nil {
			_, err := googleCalendarRepo.Update(ctx, tx, entGoogleCalendar.ID, googlecalendarinfo.GoogleCalendarInfoQueryOptions{}, entCalendar)
			if err != nil {
				return fmt.Errorf("failed to update google calendar info: %w", err)
			}
		}
	}

	// 既存のカレンダー情報を取得
	dbInfos, err := googleCalendarRepo.ListByUser(ctx, tx, entUser.ID)
	if err != nil {
		return fmt.Errorf("failed to list google calendar info: %w", err)
	}
	// 現在Googleカレンダーに存在しないカレンダー情報を削除
	for _, dbInfo := range dbInfos {
		if _, ok := incoming[dbInfo.GoogleCalendarID]; !ok {
			if err := googleCalendarRepo.SoftDelete(ctx, tx, dbInfo.ID); err != nil {
				return fmt.Errorf("failed to soft delete google calendar info: %w", err)
			}
		}
	}
	// トランザクションをコミット
	return nil
}
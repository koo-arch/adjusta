package calendar

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type SyncUsecase struct {
	userReader             UserReader
	googleTokenProvider    GoogleTokenProvider
	calendarServiceFactory CalendarServiceFactory
	tx                     SyncTransaction
}

func NewSyncUsecase(
	userReader UserReader,
	googleTokenProvider GoogleTokenProvider,
	calendarServiceFactory CalendarServiceFactory,
	tx SyncTransaction,
) *SyncUsecase {
	return &SyncUsecase{
		userReader:             userReader,
		googleTokenProvider:    googleTokenProvider,
		calendarServiceFactory: calendarServiceFactory,
		tx:                     tx,
	}
}

func (uc *SyncUsecase) SyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) ([]*customCalendar.CalendarList, error) {
	entUser, err := uc.userReader.GetByID(ctx, userID)
	if err != nil {
		log.Printf("failed to get user info for account: %s, %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("ユーザー情報が見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError("ユーザー情報取得時にエラーが発生しました")
	}

	token, err := uc.googleTokenProvider.GetToken(ctx, entUser.ID)
	if err != nil {
		log.Printf("failed to verify token for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NormalizeAPIError(err, "OAuthトークンの認証に失敗しました")
	}

	calendarService, err := uc.calendarServiceFactory.New(ctx, token)
	if err != nil {
		log.Printf("failed to create calendar service for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NewInternalError("Googleカレンダー接続に失敗しました")
	}

	calendars, err := calendarService.FetchCalendarList()
	if err != nil {
		log.Printf("failed to fetch calendars for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.FromGoogleAPIError(err)
	}

	if err := uc.syncCalendar(ctx, calendars, entUser); err != nil {
		log.Printf("failed to sync calendars for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return calendars, nil
}

func (uc *SyncUsecase) syncCalendar(ctx context.Context, calendars []*customCalendar.CalendarList, entUser *models.User) error {
	return uc.tx.Do(ctx, func(store SyncStore) error {
		incoming := make(map[string]struct{}, len(calendars))
		for _, cal := range calendars {
			incoming[cal.CalendarID] = struct{}{}

			storedCalendar, err := store.FindCalendarByGoogleCalendarID(ctx, entUser.ID, cal.CalendarID)
			if err != nil {
				if !repoerr.IsNotFound(err) {
					return fmt.Errorf("failed to find calendar: %w", err)
				}

				storedCalendar, err = store.CreateCalendar(ctx, entUser.ID)
				if err != nil {
					return fmt.Errorf("failed to create calendar: %w", err)
				}
			}

			entGoogleCalendar, err := store.FindGoogleCalendarInfoByGoogleCalendarID(ctx, cal.CalendarID)
			if err != nil {
				if !repoerr.IsNotFound(err) {
					return fmt.Errorf("failed to find google calendar info: %w", err)
				}

				_, err := store.CreateGoogleCalendarInfo(ctx, cal.CalendarID, cal.Summary, cal.Primary, storedCalendar.ID)
				if err != nil {
					return fmt.Errorf("failed to create google calendar info: %w", err)
				}
			}

			if entGoogleCalendar != nil {
				if err := store.LinkGoogleCalendarInfoToCalendar(ctx, entGoogleCalendar.ID, storedCalendar.ID); err != nil {
					return fmt.Errorf("failed to update google calendar info: %w", err)
				}
			}
		}

		dbInfos, err := store.ListGoogleCalendarInfosByUser(ctx, entUser.ID)
		if err != nil {
			return fmt.Errorf("failed to list google calendar info: %w", err)
		}

		for _, dbInfo := range dbInfos {
			if _, ok := incoming[dbInfo.GoogleCalendarID]; !ok {
				if err := store.SoftDeleteGoogleCalendarInfo(ctx, dbInfo.ID); err != nil {
					return fmt.Errorf("failed to soft delete google calendar info: %w", err)
				}
			}
		}

		return nil
	})
}

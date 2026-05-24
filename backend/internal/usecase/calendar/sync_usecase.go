package calendar

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	googleOAuth "github.com/koo-arch/adjusta-backend/internal/google/oauth"
	"github.com/koo-arch/adjusta-backend/internal/models"
	internalRepo "github.com/koo-arch/adjusta-backend/internal/repo"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	"github.com/koo-arch/adjusta-backend/utils"
)

type SyncUsecase struct {
	repos              internalRepo.Repositories
	googleTokenManager *googleOAuth.TokenManager
	uow                internalRepo.UnitOfWork
}

func NewSyncUsecase(
	repos internalRepo.Repositories,
	googleTokenManager *googleOAuth.TokenManager,
	uow internalRepo.UnitOfWork,
) *SyncUsecase {
	return &SyncUsecase{
		repos:              repos,
		googleTokenManager: googleTokenManager,
		uow:                uow,
	}
}

func (uc *SyncUsecase) SyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) ([]*customCalendar.CalendarList, error) {
	entUser, err := uc.repos.User.Read(ctx, userID, user.UserQueryOptions{})
	if err != nil {
		log.Printf("failed to get user info for account: %s, %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "ユーザー情報が見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "ユーザー情報取得時にエラーが発生しました")
	}

	token, err := uc.googleTokenManager.GetToken(ctx, entUser.ID)
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

	if err := uc.syncCalendar(ctx, calendars, entUser); err != nil {
		log.Printf("failed to sync calendars for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return calendars, nil
}

func (uc *SyncUsecase) syncCalendar(ctx context.Context, calendars []*customCalendar.CalendarList, entUser *models.User) error {
	return uc.uow.Do(ctx, func(repos internalRepo.Repositories) error {
		incoming := make(map[string]struct{}, len(calendars))
		for _, cal := range calendars {
			incoming[cal.CalendarID] = struct{}{}

			repoCalOptions := repoCalendar.CalendarQueryOptions{
				WithGoogleCalendarInfo: true,
				GoogleCalendarID:       &cal.CalendarID,
			}
			storedCalendar, err := repos.Calendar.FindByFields(ctx, entUser.ID, repoCalOptions)
			if err != nil {
				if !repoerr.IsNotFound(err) {
					return fmt.Errorf("failed to find calendar: %w", err)
				}

				storedCalendar, err = repos.Calendar.Create(ctx, entUser.ID)
				if err != nil {
					return fmt.Errorf("failed to create calendar: %w", err)
				}
			}

			gCalOptions := googlecalendarinfo.GoogleCalendarInfoQueryOptions{
				GoogleCalendarID: &cal.CalendarID,
			}
			entGoogleCalendar, err := repos.GoogleCalendarInfo.FindByFields(ctx, gCalOptions)
			if err != nil {
				if !repoerr.IsNotFound(err) {
					return fmt.Errorf("failed to find google calendar info: %w", err)
				}

				createOptions := googlecalendarinfo.GoogleCalendarInfoQueryOptions{
					GoogleCalendarID: &cal.CalendarID,
					Summary:          &cal.Summary,
					IsPrimary:        &cal.Primary,
				}
				_, err := repos.GoogleCalendarInfo.Create(ctx, createOptions, storedCalendar.ID)
				if err != nil {
					return fmt.Errorf("failed to create google calendar info: %w", err)
				}
			}

			if entGoogleCalendar != nil {
				_, err := repos.GoogleCalendarInfo.Update(ctx, entGoogleCalendar.ID, googlecalendarinfo.GoogleCalendarInfoQueryOptions{}, &storedCalendar.ID)
				if err != nil {
					return fmt.Errorf("failed to update google calendar info: %w", err)
				}
			}
		}

		dbInfos, err := repos.GoogleCalendarInfo.ListByUser(ctx, entUser.ID)
		if err != nil {
			return fmt.Errorf("failed to list google calendar info: %w", err)
		}

		for _, dbInfo := range dbInfos {
			if _, ok := incoming[dbInfo.GoogleCalendarID]; !ok {
				if err := repos.GoogleCalendarInfo.SoftDelete(ctx, dbInfo.ID); err != nil {
					return fmt.Errorf("failed to soft delete google calendar info: %w", err)
				}
			}
		}

		return nil
	})
}

package calendar

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	googleOAuth "github.com/koo-arch/adjusta-backend/internal/google/oauth"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"github.com/koo-arch/adjusta-backend/utils"
)

type SyncUsecase struct {
	client             *ent.Client
	userRepo           user.UserRepository
	calendarRepo       repoCalendar.CalendarRepository
	googleCalendarRepo googlecalendarinfo.GoogleCalendarInfoRepository
	googleTokenManager *googleOAuth.TokenManager
}

func NewSyncUsecase(
	client *ent.Client,
	userRepo user.UserRepository,
	calendarRepo repoCalendar.CalendarRepository,
	googleCalendarRepo googlecalendarinfo.GoogleCalendarInfoRepository,
	googleTokenManager *googleOAuth.TokenManager,
) *SyncUsecase {
	return &SyncUsecase{
		client:             client,
		userRepo:           userRepo,
		calendarRepo:       calendarRepo,
		googleCalendarRepo: googleCalendarRepo,
		googleTokenManager: googleTokenManager,
	}
}

func (uc *SyncUsecase) SyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) ([]*customCalendar.CalendarList, error) {
	entUser, err := uc.userRepo.Read(ctx, nil, userID, user.UserQueryOptions{})
	if err != nil {
		log.Printf("failed to get user info for account: %s, %v", email, err)
		if ent.IsNotFound(err) {
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

func (uc *SyncUsecase) syncCalendar(ctx context.Context, calendars []*customCalendar.CalendarList, entUser *ent.User) error {
	tx, err := uc.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	defer transaction.HandleTransaction(tx, &err)

	incoming := make(map[string]struct{}, len(calendars))
	for _, cal := range calendars {
		incoming[cal.CalendarID] = struct{}{}

		repoCalOptions := repoCalendar.CalendarQueryOptions{
			WithGoogleCalendarInfo: true,
			GoogleCalendarID:       &cal.CalendarID,
		}
		entCalendar, err := uc.calendarRepo.FindByFields(ctx, tx, entUser.ID, repoCalOptions)
		if err != nil {
			if !ent.IsNotFound(err) {
				return fmt.Errorf("failed to find calendar: %w", err)
			}

			entCalendar, err = uc.calendarRepo.Create(ctx, tx, entUser, nil)
			if err != nil {
				return fmt.Errorf("failed to create calendar: %w", err)
			}
		}

		gCalOptions := googlecalendarinfo.GoogleCalendarInfoQueryOptions{
			GoogleCalendarID: &cal.CalendarID,
		}
		entGoogleCalendar, err := uc.googleCalendarRepo.FindByFields(ctx, tx, gCalOptions)
		if err != nil {
			if !ent.IsNotFound(err) {
				return fmt.Errorf("failed to find google calendar info: %w", err)
			}

			createOptions := googlecalendarinfo.GoogleCalendarInfoQueryOptions{
				GoogleCalendarID: &cal.CalendarID,
				Summary:          &cal.Summary,
				IsPrimary:        &cal.Primary,
			}
			_, err := uc.googleCalendarRepo.Create(ctx, tx, createOptions, entCalendar)
			if err != nil {
				return fmt.Errorf("failed to create google calendar info: %w", err)
			}
		}

		if entGoogleCalendar != nil {
			_, err := uc.googleCalendarRepo.Update(ctx, tx, entGoogleCalendar.ID, googlecalendarinfo.GoogleCalendarInfoQueryOptions{}, entCalendar)
			if err != nil {
				return fmt.Errorf("failed to update google calendar info: %w", err)
			}
		}
	}

	dbInfos, err := uc.googleCalendarRepo.ListByUser(ctx, tx, entUser.ID)
	if err != nil {
		return fmt.Errorf("failed to list google calendar info: %w", err)
	}

	for _, dbInfo := range dbInfos {
		if _, ok := incoming[dbInfo.GoogleCalendarID]; !ok {
			if err := uc.googleCalendarRepo.SoftDelete(ctx, tx, dbInfo.ID); err != nil {
				return fmt.Errorf("failed to soft delete google calendar info: %w", err)
			}
		}
	}

	return nil
}

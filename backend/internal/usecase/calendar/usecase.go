package calendar

import (
	"context"
	"log"

	"github.com/google/uuid"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type SyncUsecase struct {
	userRepo               repoUser.UserRepository
	googleTokenProvider    GoogleTokenProvider
	calendarServiceFactory CalendarServiceFactory
	tx                     SyncTransaction
}

func NewSyncUsecase(
	userRepo repoUser.UserRepository,
	googleTokenProvider GoogleTokenProvider,
	calendarServiceFactory CalendarServiceFactory,
	tx SyncTransaction,
) *SyncUsecase {
	return &SyncUsecase{
		userRepo:               userRepo,
		googleTokenProvider:    googleTokenProvider,
		calendarServiceFactory: calendarServiceFactory,
		tx:                     tx,
	}
}

func (uc *SyncUsecase) SyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) ([]*ExternalCalendar, error) {
	entUser, err := uc.userRepo.Read(ctx, userID)
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
		return nil, internalErrors.NormalizeAPIError(err, "Googleカレンダー情報の取得に失敗しました")
	}

	calendars, err = uc.syncCalendar(ctx, calendarService, calendars, entUser)
	if err != nil {
		log.Printf("failed to sync calendars for account: %s, error: %v", entUser.Email, err)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return calendars, nil
}

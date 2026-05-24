package events

import (
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	googleOAuth "github.com/koo-arch/adjusta-backend/internal/google/oauth"
	internalRepo "github.com/koo-arch/adjusta-backend/internal/repo"
)

type Usecase struct {
	googleTokenManager *googleOAuth.TokenManager
	repos              internalRepo.Repositories
	calendarApp        *appCalendar.GoogleCalendarManager
	uow                internalRepo.UnitOfWork
}

func NewUsecase(
	googleTokenManager *googleOAuth.TokenManager,
	repos internalRepo.Repositories,
	calendarApp *appCalendar.GoogleCalendarManager,
	uow internalRepo.UnitOfWork,
) *Usecase {
	return &Usecase{
		googleTokenManager: googleTokenManager,
		repos:              repos,
		calendarApp:        calendarApp,
		uow:                uow,
	}
}

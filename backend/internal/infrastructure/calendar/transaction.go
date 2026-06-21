package calendar

import (
	"context"

	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
)

type calendarSyncTransaction struct {
	uow infraRepository.UnitOfWork
}

func NewCalendarSyncTransaction(uow infraRepository.UnitOfWork) usecaseCalendar.SyncTransaction {
	return &calendarSyncTransaction{uow: uow}
}

func (t *calendarSyncTransaction) Do(ctx context.Context, fn func(repos usecaseCalendar.SyncTxRepositories) error) error {
	return t.uow.Do(ctx, func(repos infraRepository.Repositories) error {
		return fn(usecaseCalendar.SyncTxRepositories{
			Calendar:     repos.Calendar,
			UserCalendar: repos.UserCalendar,
		})
	})
}

package account

import (
	"context"

	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
)

type CalendarSettingsTransaction struct {
	uow infraRepository.UnitOfWork
}

func NewCalendarSettingsTransaction(uow infraRepository.UnitOfWork) *CalendarSettingsTransaction {
	return &CalendarSettingsTransaction{uow: uow}
}

func (t *CalendarSettingsTransaction) DoCalendarSettings(ctx context.Context, fn func(repos usecaseAccount.CalendarSettingsRepositories) error) error {
	return t.uow.Do(ctx, func(repos infraRepository.Repositories) error {
		return fn(usecaseAccount.CalendarSettingsRepositories{
			Calendar:     repos.Calendar,
			UserCalendar: repos.UserCalendar,
		})
	})
}

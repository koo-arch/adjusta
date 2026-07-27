package googlecalendar

import (
	"context"

	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	usecaseGoogleCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/googlecalendar"
)

type syncTransaction struct {
	uow infraRepository.UnitOfWork
}

func NewSyncTransaction(uow infraRepository.UnitOfWork) usecaseGoogleCalendar.Transaction {
	return &syncTransaction{uow: uow}
}

func (t *syncTransaction) Do(ctx context.Context, fn func(repos usecaseGoogleCalendar.Repositories) error) error {
	return t.uow.Do(ctx, func(repos infraRepository.Repositories) error {
		return fn(usecaseGoogleCalendar.Repositories{
			Calendar:     repos.Calendar,
			Event:        repos.Event,
			Message:      repos.OutboxMessage,
			ProposedDate: repos.ProposedDate,
			UserCalendar: repos.UserCalendar,
		})
	})
}

package events

import (
	"context"

	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type eventTransaction struct {
	uow infraRepository.UnitOfWork
}

func NewEventTransaction(uow infraRepository.UnitOfWork) usecaseEvents.EventTransaction {
	return &eventTransaction{
		uow: uow,
	}
}

func (t *eventTransaction) DoEvent(ctx context.Context, fn func(repos usecaseEvents.EventTxRepositories) error) error {
	return t.uow.Do(ctx, func(repos infraRepository.Repositories) error {
		return fn(usecaseEvents.EventTxRepositories{
			Calendar:     repos.Calendar,
			Event:        repos.Event,
			ProposedDate: repos.ProposedDate,
			UserCalendar: repos.UserCalendar,
		})
	})
}

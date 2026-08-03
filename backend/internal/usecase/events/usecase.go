package events

import (
	"context"

	"github.com/google/uuid"
)

type OutboxDispatcher interface {
	Dispatch(ctx context.Context, messageID uuid.UUID) error
}

type Usecase struct {
	repos            EventTxRepositories
	tx               EventTransaction
	googleCalendar   GoogleCalendarGateway
	outboxDispatcher OutboxDispatcher
}

func NewUsecase(
	repos EventTxRepositories,
	tx EventTransaction,
	googleCalendar GoogleCalendarGateway,
	dispatchers ...OutboxDispatcher,
) *Usecase {
	usecase := &Usecase{
		repos:          repos,
		tx:             tx,
		googleCalendar: googleCalendar,
	}
	if len(dispatchers) > 0 {
		usecase.outboxDispatcher = dispatchers[0]
	}
	return usecase
}

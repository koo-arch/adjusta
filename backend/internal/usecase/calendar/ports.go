package calendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/google"
)

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*google.AuthToken, error)
}

type CalendarService interface {
	FetchCalendarList() ([]*CalendarRecord, error)
	CreateCalendar(summary string) (*CalendarRecord, error)
}

type CalendarServiceFactory interface {
	New(ctx context.Context, token *google.AuthToken) (CalendarService, error)
}

type CalendarServiceFactoryFunc func(ctx context.Context, token *google.AuthToken) (CalendarService, error)

func (f CalendarServiceFactoryFunc) New(ctx context.Context, token *google.AuthToken) (CalendarService, error) {
	return f(ctx, token)
}

type SyncTransaction interface {
	Do(ctx context.Context, fn func(repos SyncTxRepositories) error) error
}

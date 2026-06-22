package middlewares

import (
	"context"

	"github.com/google/uuid"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
)

type SessionAuthenticator interface {
	AuthenticateSession(ctx context.Context, sessionToken string) (*repoUser.User, error)
}

type CalendarSyncUsecase interface {
	SyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) ([]*usecaseCalendar.CalendarRecord, error)
}

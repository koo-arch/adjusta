package calendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
)

type UserReader interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*repositorymodel.User, error)
}

type UserReaderFunc func(ctx context.Context, userID uuid.UUID) (*repositorymodel.User, error)

func (f UserReaderFunc) GetByID(ctx context.Context, userID uuid.UUID) (*repositorymodel.User, error) {
	return f(ctx, userID)
}

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*appmodel.GoogleAuthToken, error)
}

type CalendarService interface {
	FetchCalendarList() ([]*customCalendar.CalendarList, error)
}

type CalendarServiceFactory interface {
	New(ctx context.Context, token *appmodel.GoogleAuthToken) (CalendarService, error)
}

type CalendarServiceFactoryFunc func(ctx context.Context, token *appmodel.GoogleAuthToken) (CalendarService, error)

func (f CalendarServiceFactoryFunc) New(ctx context.Context, token *appmodel.GoogleAuthToken) (CalendarService, error) {
	return f(ctx, token)
}

type SyncStore interface {
	FindCalendarByGoogleCalendarID(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repositorymodel.StoredCalendar, error)
	FindAnyCalendarByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*repositorymodel.StoredCalendar, error)
	CreateCalendar(ctx context.Context, googleCalendarID, summary string) (*repositorymodel.StoredCalendar, error)
	UpdateCalendar(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repositorymodel.StoredCalendar, error)
	EnsureUserCalendar(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole) (*repositorymodel.UserCalendar, error)
	ListCalendarsByUser(ctx context.Context, userID uuid.UUID) ([]*repositorymodel.StoredCalendar, error)
	SoftDeleteUserCalendar(ctx context.Context, userID, calendarID uuid.UUID) error
}

type SyncTransaction interface {
	Do(ctx context.Context, fn func(store SyncStore) error) error
}

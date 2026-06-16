package calendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
)

type UserReader interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*repoUser.User, error)
}

type UserReaderFunc func(ctx context.Context, userID uuid.UUID) (*repoUser.User, error)

func (f UserReaderFunc) GetByID(ctx context.Context, userID uuid.UUID) (*repoUser.User, error) {
	return f(ctx, userID)
}

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*appmodel.GoogleAuthToken, error)
}

type CalendarService interface {
	FetchCalendarList() ([]*customCalendar.CalendarList, error)
	CreateCalendar(summary string) (*customCalendar.CalendarList, error)
}

type CalendarServiceFactory interface {
	New(ctx context.Context, token *appmodel.GoogleAuthToken) (CalendarService, error)
}

type CalendarServiceFactoryFunc func(ctx context.Context, token *appmodel.GoogleAuthToken) (CalendarService, error)

func (f CalendarServiceFactoryFunc) New(ctx context.Context, token *appmodel.GoogleAuthToken) (CalendarService, error) {
	return f(ctx, token)
}

type UserCalendarRelationRecord struct {
	CalendarID        uuid.UUID
	GoogleCalendarID  string
	Role              domainvalue.UserCalendarRole
	SyncProposedDates bool
}

type SyncStore interface {
	FindCalendarByGoogleCalendarID(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error)
	FindAnyCalendarByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error)
	CreateCalendar(ctx context.Context, googleCalendarID, summary string) (*repoCalendar.Calendar, error)
	UpdateCalendar(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error)
	EnsureUserCalendarRelation(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole) (*repoUserCalendar.UserCalendar, error)
	ListUserCalendarRelations(ctx context.Context, userID uuid.UUID) ([]*UserCalendarRelationRecord, error)
	SoftDeleteUserCalendarRelation(ctx context.Context, userID, calendarID uuid.UUID) error
}

type SyncTransaction interface {
	Do(ctx context.Context, fn func(store SyncStore) error) error
}

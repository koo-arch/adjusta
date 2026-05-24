package calendar

import (
	"context"

	"github.com/google/uuid"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"golang.org/x/oauth2"
)

type UserReader interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type UserReaderFunc func(ctx context.Context, userID uuid.UUID) (*models.User, error)

func (f UserReaderFunc) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return f(ctx, userID)
}

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*oauth2.Token, error)
}

type CalendarService interface {
	FetchCalendarList() ([]*customCalendar.CalendarList, error)
}

type CalendarServiceFactory interface {
	New(ctx context.Context, token *oauth2.Token) (CalendarService, error)
}

type CalendarServiceFactoryFunc func(ctx context.Context, token *oauth2.Token) (CalendarService, error)

func (f CalendarServiceFactoryFunc) New(ctx context.Context, token *oauth2.Token) (CalendarService, error) {
	return f(ctx, token)
}

type SyncStore interface {
	FindCalendarByGoogleCalendarID(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*models.StoredCalendar, error)
	CreateCalendar(ctx context.Context, userID uuid.UUID) (*models.StoredCalendar, error)
	FindGoogleCalendarInfoByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*models.GoogleCalendarInfo, error)
	CreateGoogleCalendarInfo(ctx context.Context, googleCalendarID, summary string, isPrimary bool, calendarID uuid.UUID) (*models.GoogleCalendarInfo, error)
	LinkGoogleCalendarInfoToCalendar(ctx context.Context, googleCalendarInfoID, calendarID uuid.UUID) error
	ListGoogleCalendarInfosByUser(ctx context.Context, userID uuid.UUID) ([]*models.GoogleCalendarInfo, error)
	SoftDeleteGoogleCalendarInfo(ctx context.Context, id uuid.UUID) error
}

type SyncTransaction interface {
	Do(ctx context.Context, fn func(store SyncStore) error) error
}

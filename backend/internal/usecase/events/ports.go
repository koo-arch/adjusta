package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
)

type GoogleTokenProvider interface {
	GetToken(ctx context.Context, userID uuid.UUID) (*appmodel.GoogleAuthToken, error)
}

type EventTransaction interface {
	DoEvent(ctx context.Context, fn func(repos EventTxRepositories) error) error
}

type GoogleCalendarGateway interface {
	FetchEvents(ctx context.Context, userID uuid.UUID, calendars []*CalendarRecord, startTime, endTime time.Time) (*GoogleEventFetchResult, error)
	UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
}

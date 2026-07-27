package googlecalendar

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Transaction interface {
	Do(ctx context.Context, fn func(repos Repositories) error) error
}

type EventGateway interface {
	UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
	DeleteEvent(ctx context.Context, userID uuid.UUID, calendarID, googleEventID string) error
}

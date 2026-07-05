package events

import (
	"context"

	"github.com/google/uuid"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type GoogleCalendarUsecase interface {
	FetchAllGoogleEvents(ctx context.Context, userID uuid.UUID, email string) ([]*usecaseEvents.FetchedGoogleEvent, error)
}

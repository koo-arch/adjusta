package events

import (
	"context"

	"github.com/google/uuid"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type AgendaUsecase interface {
	FetchUpcomingEvents(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]usecaseEvents.UpcomingEventOutput, error)
	FetchNeedsActionDrafts(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]usecaseEvents.NeedsActionDraftOutput, error)
}

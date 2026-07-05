package events

import (
	"context"

	"github.com/google/uuid"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type DetailUsecase interface {
	FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) (*usecaseEvents.EventDraftDetailOutput, error)
}

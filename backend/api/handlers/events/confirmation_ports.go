package events

import (
	"context"

	"github.com/google/uuid"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type ConfirmationUsecase interface {
	FinalizeProposedDate(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, email string, confirmation usecaseEvents.ConfirmationRequest) error
}

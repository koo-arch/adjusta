package events

import (
	"context"

	"github.com/google/uuid"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type DraftUsecase interface {
	FetchAllDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*usecaseEvents.EventDraftDetailOutput, error)
	FetchDraftedEventsPage(ctx context.Context, userID uuid.UUID, email string, query usecaseEvents.SearchDraftQuery) (*usecaseEvents.EventDraftListOutput, error)
	SearchDraftedEvents(ctx context.Context, userID uuid.UUID, email string, query usecaseEvents.SearchDraftQuery) ([]*usecaseEvents.EventDraftDetailOutput, error)
	SearchDraftedEventsPage(ctx context.Context, userID uuid.UUID, email string, query usecaseEvents.SearchDraftQuery) (*usecaseEvents.EventDraftListOutput, error)
	CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq usecaseEvents.DraftCreationRequest) (*usecaseEvents.EventDraftDetailOutput, error)
	UpdateDraftedEvents(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, email string, eventReq usecaseEvents.DraftUpdateRequest) error
	DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) error
}

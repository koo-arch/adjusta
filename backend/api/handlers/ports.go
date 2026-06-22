package handlers

import (
	"context"

	"github.com/google/uuid"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type AccountProfileUsecase interface {
	FetchGoogleProfile(ctx context.Context, userID uuid.UUID) (*usecaseAccount.GoogleProfile, error)
}

type AuthSessionUsecase interface {
	GoogleLoginURL(state string) string
	CompleteGoogleSignIn(ctx context.Context, code string) (*usecaseAuth.GoogleSignInResult, error)
	Logout(ctx context.Context, sessionToken string) error
}

type EventUsecase interface {
	FetchAllGoogleEvents(ctx context.Context, userID uuid.UUID, email string) ([]*usecaseEvents.FetchedGoogleEvent, error)
	FetchAllDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*usecaseEvents.EventDraftDetailOutput, error)
	SearchDraftedEvents(ctx context.Context, userID uuid.UUID, email string, query usecaseEvents.SearchDraftQuery) ([]*usecaseEvents.EventDraftDetailOutput, error)
	FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) (*usecaseEvents.EventDraftDetailOutput, error)
	FetchUpcomingEvents(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]usecaseEvents.UpcomingEventOutput, error)
	FetchNeedsActionDrafts(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]usecaseEvents.NeedsActionDraftOutput, error)
	CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq usecaseEvents.DraftCreationRequest) (*usecaseEvents.EventDraftDetailOutput, error)
	FinalizeProposedDate(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, email string, confirmation usecaseEvents.ConfirmationRequest) error
	UpdateDraftedEvents(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, email string, eventReq usecaseEvents.DraftUpdateRequest) error
	DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) error
}

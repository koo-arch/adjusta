package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	infraCache "github.com/koo-arch/adjusta-backend/internal/infrastructure/cache"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type SessionAuthenticator interface {
	AuthenticateSession(ctx context.Context, sessionToken string) (*repoUser.User, error)
}

type AccountProfileService interface {
	FetchGoogleProfile(ctx context.Context, userID uuid.UUID) (*usecaseAccount.GoogleProfile, error)
}

type AuthSessionService interface {
	GoogleLoginURL(state string) string
	CompleteGoogleSignIn(ctx context.Context, code string) (*usecaseAuth.GoogleSignInResult, error)
	Logout(ctx context.Context, sessionToken string) error
}

type CalendarSyncService interface {
	SyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) ([]*usecaseCalendar.CalendarRecord, error)
}

type EventService interface {
	FetchAllGoogleEvents(ctx context.Context, userID uuid.UUID, email string) ([]*appmodel.GoogleEvent, error)
	FetchAllDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*appmodel.EventDraftDetail, error)
	SearchDraftedEvents(ctx context.Context, userID uuid.UUID, email string, query usecaseEvents.SearchDraftQuery) ([]*appmodel.EventDraftDetail, error)
	FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) (*appmodel.EventDraftDetail, error)
	FetchUpcomingEvents(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]appmodel.UpcomingEvent, error)
	FetchNeedsActionDrafts(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]appmodel.NeedsActionDraft, error)
	CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *appmodel.EventDraftCreation) (*appmodel.EventDraftDetail, error)
	FinalizeProposedDate(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, email string, confirmation usecaseEvents.ConfirmationRequest) error
	UpdateDraftedEvents(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, email string, eventReq *appmodel.EventDraftUpdate) error
	DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) error
}

type Dependencies struct {
	Cache                 *infraCache.Cache
	SessionAuthenticator  SessionAuthenticator
	AccountProfileUsecase AccountProfileService
	AuthSessionUsecase    AuthSessionService
	CalendarSyncUsecase   CalendarSyncService
	EventUsecase          EventService
}

type Server struct {
	Cache                 *infraCache.Cache
	SessionAuthenticator  SessionAuthenticator
	AccountProfileUsecase AccountProfileService
	AuthSessionUsecase    AuthSessionService
	CalendarSyncUsecase   CalendarSyncService
	EventUsecase          EventService
}

func NewServer(deps Dependencies) *Server {
	return &Server{
		Cache:                 deps.Cache,
		SessionAuthenticator:  deps.SessionAuthenticator,
		AccountProfileUsecase: deps.AccountProfileUsecase,
		AuthSessionUsecase:    deps.AuthSessionUsecase,
		CalendarSyncUsecase:   deps.CalendarSyncUsecase,
		EventUsecase:          deps.EventUsecase,
	}
}

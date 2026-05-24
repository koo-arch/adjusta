package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/cache"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type SessionAuthenticator interface {
	AuthenticateSession(ctx context.Context, sessionToken string) (*models.User, error)
}

type AccountProfileService interface {
	FetchGoogleProfile(ctx context.Context, userID uuid.UUID) (*usecaseAccount.GoogleProfile, error)
}

type AuthSessionService interface {
	CompleteGoogleSignIn(ctx context.Context, code string) (*usecaseAuth.GoogleSignInResult, error)
	Logout(ctx context.Context, sessionToken string) error
}

type CalendarSyncService interface {
	SyncGoogleCalendars(ctx context.Context, userID uuid.UUID, email string) ([]*customCalendar.CalendarList, error)
}

type EventService interface {
	FetchAllGoogleEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.GoogleEvent, error)
	FetchAllDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.EventDraftDetail, error)
	SearchDraftedEvents(ctx context.Context, userID uuid.UUID, email string, query usecaseEvents.SearchDraftQuery) ([]*models.EventDraftDetail, error)
	FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email, slug string) (*models.EventDraftDetail, error)
	FetchUpcomingEvents(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]models.UpcomingEvent, error)
	FetchNeedsActionDrafts(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]models.NeedsActionDraft, error)
	CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftCreation) (*models.EventDraftDetail, error)
	FinalizeProposedDate(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *models.ConfirmEvent) error
	UpdateDraftedEvents(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *models.EventDraftUpdate) error
	DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftDetail) error
}

type Dependencies struct {
	Cache                 *cache.Cache
	SessionAuthenticator  SessionAuthenticator
	AccountProfileUsecase AccountProfileService
	AuthSessionUsecase    AuthSessionService
	CalendarSyncUsecase   CalendarSyncService
	EventUsecase          EventService
}

type Server struct {
	Cache                 *cache.Cache
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

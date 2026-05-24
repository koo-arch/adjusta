package api

import (
	"github.com/koo-arch/adjusta-backend/cache"
	"github.com/koo-arch/adjusta-backend/ent"
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	googleOAuth "github.com/koo-arch/adjusta-backend/internal/google/oauth"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/repo/account"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	repoSession "github.com/koo-arch/adjusta-backend/internal/repo/session"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type Server struct {
	Client                *ent.Client
	Cache                 *cache.Cache
	UserRepo              user.UserRepository
	AccountRepo           repoAccount.AccountRepository
	SessionRepo           repoSession.SessionRepository
	CalendarRepo          dbCalendar.CalendarRepository
	GoogleCalendarRepo    googlecalendarinfo.GoogleCalendarInfoRepository
	EventRepo             event.EventRepository
	DateRepo              proposeddate.ProposedDateRepository
	AuthManager           *auth.AuthManager
	GoogleTokenManager    *googleOAuth.TokenManager
	AccountProfileUsecase *usecaseAccount.ProfileUsecase
	AuthSessionUsecase    *usecaseAuth.SessionUsecase
	CalendarSyncUsecase   *usecaseCalendar.SyncUsecase
	EventUsecase          *usecaseEvents.Usecase
}

func NewServer(client *ent.Client) *Server {
	cache := cache.NewCache()

	userRepo := user.NewUserRepository(client)
	accountRepo := repoAccount.NewAccountRepository(client)
	sessionRepo := repoSession.NewSessionRepository(client)
	calendarRepo := dbCalendar.NewCalendarRepository(client)
	googleCalendarRepo := googlecalendarinfo.NewGoogleCalendarInfoRepository(client)
	eventRepo := event.NewEventRepository(client)
	dateRepo := proposeddate.NewProposedDateRepository(client)
	calendarApp := appCalendar.NewGoogleCalendarManager(client) // Google Calendar API manager
	authManager := auth.NewAuthManager(client, userRepo, accountRepo, sessionRepo)
	googleTokenManager := googleOAuth.NewTokenManager(accountRepo)
	accountProfileUsecase := usecaseAccount.NewProfileUsecase(googleTokenManager)
	authSessionUsecase := usecaseAuth.NewSessionUsecase(authManager)
	calendarSyncUsecase := usecaseCalendar.NewSyncUsecase(client, userRepo, calendarRepo, googleCalendarRepo, googleTokenManager)
	eventUsecase := usecaseEvents.NewUsecase(client, googleTokenManager, calendarRepo, googleCalendarRepo, eventRepo, dateRepo, calendarApp)

	return &Server{
		Client:                client,
		Cache:                 cache,
		UserRepo:              userRepo,
		AccountRepo:           accountRepo,
		SessionRepo:           sessionRepo,
		CalendarRepo:          calendarRepo,
		GoogleCalendarRepo:    googleCalendarRepo,
		EventRepo:             eventRepo,
		DateRepo:              dateRepo,
		AuthManager:           authManager,
		GoogleTokenManager:    googleTokenManager,
		AccountProfileUsecase: accountProfileUsecase,
		AuthSessionUsecase:    authSessionUsecase,
		CalendarSyncUsecase:   calendarSyncUsecase,
		EventUsecase:          eventUsecase,
	}
}

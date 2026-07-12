package app

import (
	"time"

	"github.com/gin-contrib/sessions"
	apiCookie "github.com/koo-arch/adjusta-backend/api/cookie"
	accountHandlers "github.com/koo-arch/adjusta-backend/api/handlers/account"
	eventHandlers "github.com/koo-arch/adjusta-backend/api/handlers/events"
	oauthHandlers "github.com/koo-arch/adjusta-backend/api/handlers/oauth"
	userHandlers "github.com/koo-arch/adjusta-backend/api/handlers/user"
	"github.com/koo-arch/adjusta-backend/api/middlewares"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
	"github.com/koo-arch/adjusta-backend/internal/config"
	infraAccount "github.com/koo-arch/adjusta-backend/internal/infrastructure/account"
	infraAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/auth"
	infraCache "github.com/koo-arch/adjusta-backend/internal/infrastructure/cache"
	infraCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/calendar"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"
	infraEvents "github.com/koo-arch/adjusta-backend/internal/infrastructure/events"
	infraGoogleCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/googlecalendar"
	infraGoogleOAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/googleoauth"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	"github.com/koo-arch/adjusta-backend/internal/usecase/account/calendarsetting"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type dependencies struct {
	accountHandler     *accountHandlers.Handler
	userHandler        *userHandlers.Handler
	oauthHandler       *oauthHandlers.Handler
	eventHandler       *eventHandlers.Handler
	authMiddleware     *middlewares.AuthMiddleware
	calendarMiddleware *middlewares.CalendarMiddleware
	sessionMiddleware  *middlewares.SessionMiddleware
	cookieOptions      sessions.Options
}

func buildDependencies(client *ent.Client, cfg config.Config) *dependencies {
	calendarCache := infraCache.NewCalendarCache(5*time.Minute, 10*time.Minute)
	repos := infraRepository.NewRepositories(client)
	uow := infraRepository.NewUnitOfWork(client)
	calendarApp := infraGoogleCalendar.NewGoogleCalendarManager()
	googleOAuthClient := infraGoogleOAuth.NewClient(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURI)
	cookieManager := apiCookie.NewManager(cfg.Domain, !cfg.IsDevelopment())
	cookieSessionStore := sessionctx.NewCookieSessionStore(cookieManager)
	cookieOptions := cookieManager.Options()
	sessionLifetime := time.Duration(cookieOptions.MaxAge) * time.Second

	authenticator := usecaseAuth.NewAuthenticator(
		usecaseAuth.AuthRepositories{
			User:    repos.User,
			Account: repos.Account,
			Session: repos.Session,
		},
		infraAuth.NewAuthTransaction(uow),
		sessionLifetime,
	)
	googleTokenManager := infraGoogleOAuth.NewTokenManager(repos.Account, googleOAuthClient)
	accountProfileUsecase := usecaseAccount.NewProfileUsecase(
		googleTokenManager,
		infraAuth.NewGoogleUserInfoFetcher(googleOAuthClient),
	)
	oauthUsecase := usecaseAuth.NewOAuthUsecase(
		authenticator,
		infraAuth.NewGoogleOAuthGateway(googleOAuthClient),
		infraAuth.NewGoogleUserInfoFetcher(googleOAuthClient),
	)
	calendarSyncUsecase := usecaseCalendar.NewSyncUsecase(
		repos.User,
		googleTokenManager,
		infraGoogleCalendar.NewCalendarServiceFactory(googleOAuthClient),
		infraCalendar.NewCalendarSyncTransaction(uow),
		calendarCache,
	)
	calendarSettingsUsecase := calendarsetting.NewUsecase(
		calendarsetting.CalendarSettingsRepositories{
			Calendar:     repos.Calendar,
			UserCalendar: repos.UserCalendar,
		},
		infraAccount.NewCalendarSettingsTransaction(uow),
		calendarSyncUsecase,
	)
	eventUsecase := usecaseEvents.NewUsecase(
		usecaseEvents.EventTxRepositories{
			Calendar:     repos.Calendar,
			Event:        repos.Event,
			ProposedDate: repos.ProposedDate,
			UserCalendar: repos.UserCalendar,
		},
		infraEvents.NewEventTransaction(uow),
		infraGoogleCalendar.NewEventGateway(googleTokenManager, calendarApp, googleOAuthClient),
	)

	return &dependencies{
		accountHandler:     accountHandlers.NewHandler(accountProfileUsecase, calendarSettingsUsecase),
		userHandler:        userHandlers.NewHandler(accountProfileUsecase),
		oauthHandler:       oauthHandlers.NewHandler(oauthUsecase, cfg.RedirectURLAfterLogin, cookieSessionStore),
		eventHandler:       eventHandlers.NewHandler(eventUsecase, eventUsecase, eventUsecase, eventUsecase, eventUsecase),
		authMiddleware:     middlewares.NewAuthMiddleware(authenticator, cookieSessionStore),
		calendarMiddleware: middlewares.NewCalendarMiddleware(calendarSyncUsecase),
		sessionMiddleware:  middlewares.NewSessionMiddleware(cookieSessionStore),
		cookieOptions:      cookieOptions,
	}
}

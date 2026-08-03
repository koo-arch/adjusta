package app

import (
	"context"
	"time"

	"github.com/gin-contrib/sessions"
	apiCookie "github.com/koo-arch/adjusta-backend/api/cookie"
	accountHandlers "github.com/koo-arch/adjusta-backend/api/handlers/account"
	eventHandlers "github.com/koo-arch/adjusta-backend/api/handlers/events"
	internalTaskHandlers "github.com/koo-arch/adjusta-backend/api/handlers/internaltasks"
	oauthHandlers "github.com/koo-arch/adjusta-backend/api/handlers/oauth"
	userHandlers "github.com/koo-arch/adjusta-backend/api/handlers/user"
	"github.com/koo-arch/adjusta-backend/api/middlewares"
	"github.com/koo-arch/adjusta-backend/api/sessionctx"
	"github.com/koo-arch/adjusta-backend/internal/config"
	infraAccount "github.com/koo-arch/adjusta-backend/internal/infrastructure/account"
	infraAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/auth"
	infraCache "github.com/koo-arch/adjusta-backend/internal/infrastructure/cache"
	infraCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/calendar"
	infraCloudTasks "github.com/koo-arch/adjusta-backend/internal/infrastructure/cloudtasks"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"
	infraEvents "github.com/koo-arch/adjusta-backend/internal/infrastructure/events"
	infraGoogleCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/googlecalendar"
	infraGoogleOAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/googleoauth"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	infraTaskAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/taskauth"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	"github.com/koo-arch/adjusta-backend/internal/usecase/account/calendarsetting"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
	usecaseGoogleCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/googlecalendar"
	usecaseOutbox "github.com/koo-arch/adjusta-backend/internal/usecase/outbox"
)

type dependencies struct {
	accountHandler      *accountHandlers.Handler
	userHandler         *userHandlers.Handler
	oauthHandler        *oauthHandlers.Handler
	eventHandler        *eventHandlers.Handler
	authMiddleware      *middlewares.AuthMiddleware
	calendarMiddleware  *middlewares.CalendarMiddleware
	sessionMiddleware   *middlewares.SessionMiddleware
	internalTaskHandler *internalTaskHandlers.Handler
	cookieOptions       sessions.Options
}

func buildDependencies(client *ent.Client, cfg config.Config) (*dependencies, error) {
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

	var outboxDispatcher *usecaseOutbox.Dispatcher
	var eventDispatchers []usecaseEvents.OutboxDispatcher
	if cfg.CloudTasksEnabled() {
		publisher, err := infraCloudTasks.NewPublisher(context.Background(), infraCloudTasks.PublisherConfig{
			ProjectID:             cfg.CloudTasksProjectID,
			Location:              cfg.CloudTasksLocation,
			Queue:                 cfg.CloudTasksQueue,
			HandlerURL:            cfg.CloudTasksHandlerURL,
			OIDCAudience:          cfg.CloudTasksOIDCAudience,
			InvokerServiceAccount: cfg.CloudTasksInvokerServiceAccount,
		})
		if err != nil {
			return nil, err
		}
		outboxDispatcher = usecaseOutbox.NewDispatcher(repos.OutboxMessage, publisher)
		eventDispatchers = append(eventDispatchers, outboxDispatcher)
	}

	googleCalendarGateway := infraGoogleCalendar.NewEventGateway(googleTokenManager, calendarApp, googleOAuthClient)
	eventUsecase := usecaseEvents.NewUsecase(
		usecaseEvents.EventTxRepositories{
			Calendar:      repos.Calendar,
			Event:         repos.Event,
			OutboxMessage: repos.OutboxMessage,
			ProposedDate:  repos.ProposedDate,
			UserCalendar:  repos.UserCalendar,
		},
		infraEvents.NewEventTransaction(uow),
		googleCalendarGateway,
		eventDispatchers...,
	)

	var internalTaskHandler *internalTaskHandlers.Handler
	if outboxDispatcher != nil {
		googleCalendarSyncUsecase := usecaseGoogleCalendar.NewSyncUsecase(
			usecaseGoogleCalendar.Repositories{
				Calendar:     repos.Calendar,
				Event:        repos.Event,
				Message:      repos.OutboxMessage,
				ProposedDate: repos.ProposedDate,
				UserCalendar: repos.UserCalendar,
			},
			infraGoogleCalendar.NewSyncTransaction(uow),
			googleCalendarGateway,
		)
		internalTaskHandler = internalTaskHandlers.NewHandler(
			googleCalendarSyncUsecase,
			outboxDispatcher,
			infraTaskAuth.NewOIDCVerifier(cfg.CloudTasksOIDCAudience, cfg.CloudTasksInvokerServiceAccount),
		)
	}

	return &dependencies{
		accountHandler:      accountHandlers.NewHandler(accountProfileUsecase, calendarSettingsUsecase),
		userHandler:         userHandlers.NewHandler(accountProfileUsecase),
		oauthHandler:        oauthHandlers.NewHandler(oauthUsecase, cfg.RedirectURLAfterLogin, cookieSessionStore),
		eventHandler:        eventHandlers.NewHandler(eventUsecase, eventUsecase, eventUsecase, eventUsecase, eventUsecase),
		authMiddleware:      middlewares.NewAuthMiddleware(authenticator, cookieSessionStore),
		calendarMiddleware:  middlewares.NewCalendarMiddleware(calendarSyncUsecase),
		sessionMiddleware:   middlewares.NewSessionMiddleware(cookieSessionStore),
		internalTaskHandler: internalTaskHandler,
		cookieOptions:       cookieOptions,
	}, nil
}

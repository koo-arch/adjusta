package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	apiCookie "github.com/koo-arch/adjusta-backend/api/cookie"
	accountHandlers "github.com/koo-arch/adjusta-backend/api/handlers/account"
	eventHandlers "github.com/koo-arch/adjusta-backend/api/handlers/events"
	oauthHandlers "github.com/koo-arch/adjusta-backend/api/handlers/oauth"
	userHandlers "github.com/koo-arch/adjusta-backend/api/handlers/user"
	"github.com/koo-arch/adjusta-backend/api/middlewares"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/config"
	infraAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/auth"
	infraCache "github.com/koo-arch/adjusta-backend/internal/infrastructure/cache"
	infraCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/calendar"
	infraEvents "github.com/koo-arch/adjusta-backend/internal/infrastructure/events"
	infraGoogleCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/googlecalendar"
	infraGoogleOAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/googleoauth"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"

	_ "github.com/koo-arch/adjusta-backend/ent/runtime"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()

	// DB接続
	databaseURL := cfg.DatabaseURL
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	client, err := ent.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("failed closing connection to postgres: %v", err)
		}
	}()

	// データベースのスキーマを更新
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	calendarCache := infraCache.NewCalendarCache(5*time.Minute, 10*time.Minute)
	repos := infraRepository.NewRepositories(client)
	uow := infraRepository.NewUnitOfWork(client)
	calendarApp := infraGoogleCalendar.NewGoogleCalendarManager()
	googleOAuthClient := infraGoogleOAuth.NewClient(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURI)
	cookieManager := apiCookie.NewManager(cfg.Domain, cfg.GoEnv != "development")
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

	//Ginフレームワークのデフォルトの設定を使用してルータを作成
	router := gin.Default()

	// CORSの設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(cookieOptions)
	router.Use(sessions.Sessions(apiCookie.SessionCookieName, store))

	accountHandler := accountHandlers.NewHandler(accountProfileUsecase)
	userHandler := userHandlers.NewHandler(accountProfileUsecase)
	oauthHandler := oauthHandlers.NewHandler(oauthUsecase, cfg.RedirectURLAfterLogin, cookieManager)
	eventHandler := eventHandlers.NewHandler(
		eventUsecase,
		eventUsecase,
		eventUsecase,
		eventUsecase,
		eventUsecase,
	)

	authMiddleware := middlewares.NewAuthMiddleware(authenticator, cookieManager)
	calendarMiddleware := middlewares.NewCalendarMiddleware(calendarSyncUsecase)
	sessionMiddleware := middlewares.NewSessionMiddleware(cookieManager)

	// ルートハンドラの定義
	router.GET("/auth/google/login", oauthHandler.GoogleLoginHandler)
	router.GET("/auth/google/callback", oauthHandler.GoogleCallbackHandler())
	router.GET("/auth/logout", oauthHandler.LogoutHandler)

	// 認証が必要なAPIグループ
	auth := router.Group("/api")
	auth.Use(sessionMiddleware.SessionRenewal(), authMiddleware.AuthUser())
	{
		auth.GET("/users/me", userHandler.GetCurrentUserHandler())
		auth.GET("/account/list", accountHandler.FetchAccountsHandler())
		calendar := auth.Group("/calendar").Use(calendarMiddleware.SyncGoogleCalendars())
		{
			calendar.GET("/list", eventHandler.FetchEventListHandler())
			calendar.GET("/event/draft/list", eventHandler.FetchAllEventDraftListHandler())
			calendar.GET("/event/draft/:id", eventHandler.FetchEventDraftDetailHandler())
			calendar.POST("/event/draft", eventHandler.CreateEventDraftHandler())
			calendar.PATCH("/event/confirm/:id", eventHandler.EventFinalizeHandler())
			calendar.PUT("/event/draft/:id", eventHandler.UpdateEventDraftHandler())
			calendar.DELETE("/event/draft/:id", eventHandler.DeleteEventDraftHandler())
		}

		auth.GET("/event/draft/search", eventHandler.SearchEventDraftHandler())
		auth.GET("/event/confirmed/upcoming", eventHandler.FetchUpcomingEventsHandler())
		auth.GET("/event/draft/needs-action", eventHandler.FetchNeedsActionDraftsHandler())
	}

	// サーバー起動
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}

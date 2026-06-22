package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/handlers"
	"github.com/koo-arch/adjusta-backend/api/middlewares"
	"github.com/koo-arch/adjusta-backend/ent"
	infraAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/auth"
	infraCache "github.com/koo-arch/adjusta-backend/internal/infrastructure/cache"
	infraCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/calendar"
	infraConfigs "github.com/koo-arch/adjusta-backend/internal/infrastructure/configs"
	infraCookie "github.com/koo-arch/adjusta-backend/internal/infrastructure/cookie"
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
	// 環境変数の読み込み
	infraConfigs.LoadEnv()

	// DB接続
	databaseURL := infraConfigs.GetEnv("DATABASE_URL")
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

	cacheStore := infraCache.NewCache()
	repos := infraRepository.NewRepositories(client)
	uow := infraRepository.NewUnitOfWork(client)
	calendarApp := infraGoogleCalendar.NewGoogleCalendarManager()
	sessionLifetime := time.Duration(infraCookie.DefaultCookieOptions().MaxAge) * time.Second
	authService := usecaseAuth.NewAuthService(
		usecaseAuth.AuthRepositories{
			User:    repos.User,
			Account: repos.Account,
			Session: repos.Session,
		},
		infraAuth.NewAuthTransaction(uow),
		sessionLifetime,
	)
	googleTokenManager := infraGoogleOAuth.NewTokenManager(repos.Account)
	accountProfileUsecase := usecaseAccount.NewProfileUsecase(
		googleTokenManager,
		infraAuth.NewGoogleUserInfoFetcher(),
	)
	authSessionUsecase := usecaseAuth.NewSessionUsecase(
		authService,
		infraAuth.NewGoogleOAuthGateway(),
		infraAuth.NewGoogleUserInfoFetcher(),
	)
	calendarSyncUsecase := usecaseCalendar.NewSyncUsecase(
		repos.User,
		googleTokenManager,
		infraGoogleCalendar.NewCalendarServiceFactory(),
		infraCalendar.NewCalendarSyncTransaction(uow),
	)
	eventUsecase := usecaseEvents.NewUsecase(
		usecaseEvents.EventTxRepositories{
			Calendar:     repos.Calendar,
			Event:        repos.Event,
			ProposedDate: repos.ProposedDate,
			UserCalendar: repos.UserCalendar,
		},
		infraEvents.NewEventTransaction(uow),
		infraGoogleCalendar.NewEventGateway(googleTokenManager, calendarApp),
	)

	//Ginフレームワークのデフォルトの設定を使用してルータを作成
	router := gin.Default()

	// CORSの設定
	allowedOrigins := strings.Split(infraConfigs.GetEnv("CORS_ALLOW_ORIGINS"), ",")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte(infraConfigs.GetEnv("SESSION_SECRET")))
	store.Options(infraCookie.DefaultCookieOptions())
	router.Use(sessions.Sessions(infraCookie.SessionCookieName, store))

	accountHandler := handlers.NewAccountHandler(accountProfileUsecase)
	userHandler := handlers.NewUserHandler(accountProfileUsecase)
	oauthHandler := handlers.NewOauthHandler(authSessionUsecase)
	eventHandler := handlers.NewEventHandler(eventUsecase)

	authMiddleware := middlewares.NewAuthMiddleware(authService)
	calendarMiddleware := middlewares.NewCalendarMiddleware(cacheStore, calendarSyncUsecase)
	sessionMiddleware := middlewares.NewSessionMiddleware()

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
	port := infraConfigs.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}

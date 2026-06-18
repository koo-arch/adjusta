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
	"github.com/koo-arch/adjusta-backend/api"
	"github.com/koo-arch/adjusta-backend/api/handlers"
	"github.com/koo-arch/adjusta-backend/api/middlewares"
	"github.com/koo-arch/adjusta-backend/cache"
	"github.com/koo-arch/adjusta-backend/configs"
	opCookie "github.com/koo-arch/adjusta-backend/cookie"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	googleOAuth "github.com/koo-arch/adjusta-backend/internal/google/oauth"
	googleUserInfo "github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	infraAuth "github.com/koo-arch/adjusta-backend/internal/infrastructure/auth"
	infraCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/calendar"
	infraEvents "github.com/koo-arch/adjusta-backend/internal/infrastructure/events"
	infraGoogleCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/googlecalendar"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
	"golang.org/x/oauth2"

	_ "github.com/koo-arch/adjusta-backend/ent/runtime"
	_ "github.com/lib/pq"
)

func main() {
	// 環境変数の読み込み
	configs.LoadEnv()

	// DB接続
	databaseURL := configs.GetEnv("DATABASE_URL")
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

	cacheStore := cache.NewCache()
	repos := infraRepository.NewRepositories(client)
	uow := infraRepository.NewUnitOfWork(client, repos)
	calendarApp := infraGoogleCalendar.NewGoogleCalendarManager()
	sessionLifetime := time.Duration(opCookie.DefaultCookieOptions().MaxAge) * time.Second
	authService := usecaseAuth.NewAuthService(
		infraAuth.NewAuthReader(repos.User, repos.Account),
		infraAuth.NewAuthTransaction(uow),
		infraAuth.NewAuthSessionStore(repos.Session),
		sessionLifetime,
	)
	googleTokenManager := googleOAuth.NewTokenManager(repos.Account)
	accountProfileUsecase := usecaseAccount.NewProfileUsecase(
		googleTokenManager,
		usecaseAccount.UserInfoFetcherFunc(fetchGoogleUserProfile),
	)
	authSessionUsecase := usecaseAuth.NewSessionUsecase(
		authService,
		usecaseAuth.OAuthGatewayFuncs{
			AuthCodeURLFn: buildGoogleAuthURL,
			ExchangeFn:    exchangeGoogleAuthToken,
		},
		usecaseAuth.UserInfoFetcherFunc(fetchGoogleUserProfile),
	)
	calendarSyncUsecase := usecaseCalendar.NewSyncUsecase(
		infraCalendar.NewCalendarSyncUserReader(repos.User),
		googleTokenManager,
		infraGoogleCalendar.NewCalendarServiceFactory(),
		infraCalendar.NewCalendarSyncTransaction(uow),
	)
	eventUsecase := usecaseEvents.NewUsecase(
		infraEvents.NewEventReader(repos),
		infraEvents.NewEventTransaction(uow),
		infraGoogleCalendar.NewEventGateway(googleTokenManager, calendarApp),
	)

	server := api.NewServer(api.Dependencies{
		Cache:                 cacheStore,
		SessionAuthenticator:  authService,
		AccountProfileUsecase: accountProfileUsecase,
		AuthSessionUsecase:    authSessionUsecase,
		CalendarSyncUsecase:   calendarSyncUsecase,
		EventUsecase:          eventUsecase,
	})

	//Ginフレームワークのデフォルトの設定を使用してルータを作成
	router := gin.Default()

	// CORSの設定
	allowedOrigins := strings.Split(configs.GetEnv("CORS_ALLOW_ORIGINS"), ",")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte(configs.GetEnv("SESSION_SECRET")))
	store.Options(opCookie.DefaultCookieOptions())
	router.Use(sessions.Sessions("session", store))

	handler := handlers.NewHandler(server)
	accountHandler := handlers.NewAccountHandler(handler)
	userHandler := handlers.NewUserHandler(handler)
	oauthHandler := handlers.NewOauthHandler(handler)
	calendarHandler := handlers.NewCalendarHandler(handler)

	middleware := middlewares.NewMiddleware(server)
	authMiddleware := middlewares.NewAuthMiddleware(middleware)
	calendarMiddleware := middlewares.NewCalendarMiddleware(middleware)
	sessionMiddleware := middlewares.NewSessionMiddleware(middleware)

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
			calendar.GET("/list", calendarHandler.FetchEventListHandler())
			calendar.GET("/event/draft/list", calendarHandler.FetchAllEventDraftListHandler())
			calendar.GET("/event/draft/:id", calendarHandler.FetchEventDraftDetailHandler())
			calendar.POST("/event/draft", calendarHandler.CreateEventDraftHandler())
			calendar.PATCH("/event/confirm/:id", calendarHandler.EventFinalizeHandler())
			calendar.PUT("/event/draft/:id", calendarHandler.UpdateEventDraftHandler())
			calendar.DELETE("/event/draft/:id", calendarHandler.DeleteEventDraftHandler())
		}

		auth.GET("/event/draft/search", calendarHandler.SearchEventDraftHandler())
		auth.GET("/event/confirmed/upcoming", calendarHandler.FetchUpcomingEventsHandler())
		auth.GET("/event/draft/needs-action", calendarHandler.FetchNeedsActionDraftsHandler())
	}

	// サーバー起動
	port := configs.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}

func exchangeGoogleAuthToken(ctx context.Context, code string) (*appmodel.GoogleAuthToken, error) {
	token, err := googleOAuth.GetGoogleAuthConfig().Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	scope := token.Extra("scope")
	scopeValue, ok := scope.(string)
	var tokenScope *string
	if ok && scopeValue != "" {
		tokenScope = &scopeValue
	}

	return &appmodel.GoogleAuthToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		Scope:        tokenScope,
	}, nil
}

func buildGoogleAuthURL(state string) string {
	return googleOAuth.GetGoogleAuthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func fetchGoogleUserProfile(ctx context.Context, token *appmodel.GoogleAuthToken) (*appmodel.GoogleUserProfile, error) {
	userInfo, err := googleUserInfo.FetchGoogleUserInfo(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
	if err != nil {
		return nil, err
	}

	return &appmodel.GoogleUserProfile{
		GoogleID: userInfo.GoogleID,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Picture:  userInfo.Picture,
	}, nil
}

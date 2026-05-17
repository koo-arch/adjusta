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
	"github.com/koo-arch/adjusta-backend/configs"
	"github.com/koo-arch/adjusta-backend/ent"
	opCookie "github.com/koo-arch/adjusta-backend/cookie"

	_ "github.com/koo-arch/adjusta-backend/ent/runtime"
	"github.com/koo-arch/adjusta-backend/scheduler"
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

	server := api.NewServer(client)

	// cache
	cache := server.Cache

	// JWTキーの起動時生成
	keyManager := server.KeyManager
	if err := keyManager.InitializeJWTKeys(ctx); err != nil {
		log.Fatalf("failed to initialize JWT")
	}

	// スケジューラーのセットアップ
	s := scheduler.NewScheduler(client, cache)
	s.SetupJobs(ctx)
	s.Start()
	defer s.Stop()

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
			calendar.GET("/event/draft/:slug", calendarHandler.FetchEventDraftDetailHandler())
			calendar.POST("/event/draft", calendarHandler.CreateEventDraftHandler())
			calendar.PATCH("/event/confirm/:slug", calendarHandler.EventFinalizeHandler())
			calendar.PUT("/event/draft/:slug", calendarHandler.UpdateEventDraftHandler())
			calendar.DELETE("/event/draft/:slug", calendarHandler.DeleteEventDraftHandler())
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

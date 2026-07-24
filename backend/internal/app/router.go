package app

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	apiCookie "github.com/koo-arch/adjusta-backend/api/cookie"
	"github.com/koo-arch/adjusta-backend/internal/config"
)

func newRouter(cfg config.Config, deps *dependencies) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(deps.cookieOptions)
	router.Use(sessions.Sessions(apiCookie.SessionCookieName, store))

	registerRoutes(router, deps)

	return router
}

func registerRoutes(router *gin.Engine, deps *dependencies) {
	router.GET("/auth/google/login", deps.oauthHandler.GoogleLoginHandler)
	router.GET("/auth/google/callback", deps.oauthHandler.GoogleCallbackHandler())
	router.GET("/auth/logout", deps.oauthHandler.LogoutHandler)

	auth := router.Group("/api")
	auth.Use(deps.sessionMiddleware.SessionRenewal(), deps.authMiddleware.AuthUser())
	{
		auth.GET("/auth/google/reauthorize", deps.oauthHandler.GoogleReauthorizationHandler)
		auth.GET("/users/me", deps.userHandler.GetCurrentUserHandler())
		auth.GET("/account/list", deps.accountHandler.FetchAccountsHandler())
		auth.GET("/user-calendars", deps.accountHandler.ListCalendarSettingsHandler())
		auth.PATCH("/user-calendars/:id", deps.accountHandler.UpdateCalendarSettingHandler())
		auth.GET("/calendar-settings/candidate-sync", deps.accountHandler.GetCandidateSyncSettingHandler())
		auth.PUT("/calendar-settings/candidate-sync", deps.accountHandler.SetCandidateSyncSettingHandler())
		calendar := auth.Group("/calendar").Use(deps.calendarMiddleware.SyncGoogleCalendars())
		{
			calendar.GET("/list", deps.eventHandler.FetchEventListHandler())
			calendar.GET("/event/draft/list", deps.eventHandler.FetchAllEventDraftListHandler())
			calendar.GET("/event/draft/:id", deps.eventHandler.FetchEventDraftDetailHandler())
			calendar.POST("/event/draft", deps.eventHandler.CreateEventDraftHandler())
			calendar.PATCH("/event/confirm/:id", deps.eventHandler.EventFinalizeHandler())
			calendar.PUT("/event/draft/:id", deps.eventHandler.UpdateEventDraftHandler())
			calendar.DELETE("/event/draft/:id", deps.eventHandler.DeleteEventDraftHandler())
		}

		auth.GET("/event/draft/search", deps.eventHandler.SearchEventDraftHandler())
		auth.GET("/event/confirmed/upcoming", deps.eventHandler.FetchUpcomingEventsHandler())
		auth.GET("/event/draft/needs-action", deps.eventHandler.FetchNeedsActionDraftsHandler())
	}
}

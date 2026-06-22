package middlewares

import (
	"github.com/koo-arch/adjusta-backend/api"
	infraCache "github.com/koo-arch/adjusta-backend/internal/infrastructure/cache"
)

type Dependencies struct {
	Cache                *infraCache.Cache
	SessionAuthenticator api.SessionAuthenticator
	CalendarSyncService  api.CalendarSyncService
}

type Middleware struct {
	cache                *infraCache.Cache
	sessionAuthenticator api.SessionAuthenticator
	calendarSyncService  api.CalendarSyncService
}

func NewMiddleware(deps Dependencies) *Middleware {
	return &Middleware{
		cache:                deps.Cache,
		sessionAuthenticator: deps.SessionAuthenticator,
		calendarSyncService:  deps.CalendarSyncService,
	}
}

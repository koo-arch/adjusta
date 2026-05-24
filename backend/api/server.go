package api

import (
	"github.com/koo-arch/adjusta-backend/cache"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
	usecaseAuth "github.com/koo-arch/adjusta-backend/internal/usecase/auth"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type Dependencies struct {
	Cache                 *cache.Cache
	AuthManager           *auth.AuthManager
	AccountProfileUsecase *usecaseAccount.ProfileUsecase
	AuthSessionUsecase    *usecaseAuth.SessionUsecase
	CalendarSyncUsecase   *usecaseCalendar.SyncUsecase
	EventUsecase          *usecaseEvents.Usecase
}

type Server struct {
	Cache                 *cache.Cache
	AuthManager           *auth.AuthManager
	AccountProfileUsecase *usecaseAccount.ProfileUsecase
	AuthSessionUsecase    *usecaseAuth.SessionUsecase
	CalendarSyncUsecase   *usecaseCalendar.SyncUsecase
	EventUsecase          *usecaseEvents.Usecase
}

func NewServer(deps Dependencies) *Server {
	return &Server{
		Cache:                 deps.Cache,
		AuthManager:           deps.AuthManager,
		AccountProfileUsecase: deps.AccountProfileUsecase,
		AuthSessionUsecase:    deps.AuthSessionUsecase,
		CalendarSyncUsecase:   deps.CalendarSyncUsecase,
		EventUsecase:          deps.EventUsecase,
	}
}

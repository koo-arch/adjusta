package api

import (
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/cache"
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	appEvents "github.com/koo-arch/adjusta-backend/internal/apps/events"
	"github.com/koo-arch/adjusta-backend/internal/apps/events/event_operations"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/repo/oauthtoken"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
)

type Server struct {
	Client 	      			 *ent.Client
	Cache 	   			 	 *cache.Cache
	UserRepo      			 user.UserRepository
	OAuthRepo   			 oauthtoken.OAuthTokenRepository
	CalendarRepo  			 dbCalendar.CalendarRepository
	GoogleCalendarRepo   	 googlecalendarinfo.GoogleCalendarInfoRepository
	EventRepo     			 event.EventRepository
	DateRepo      			 proposeddate.ProposedDateRepository
	AuthManager   			 *auth.AuthManager
	JWTManager    			 *auth.JWTManager
	KeyManager    			 *auth.KeyManager
	EventManager 			 *appEvents.EventManager
	EventCreationManager     *event_operations.EventCreationManager
	EventFetchingManager     *event_operations.EventFetchingManager
	EventUpdateManager       *event_operations.EventUpdateManager
	EventDeleteManager 	 	 *event_operations.EventDeleteManager
}

func NewServer(client *ent.Client) *Server {
	cache := cache.NewCache()

	userRepo := user.NewUserRepository(client)
	oauthRepo := oauthtoken.NewOAuthTokenRepository(client)
	calendarRepo := dbCalendar.NewCalendarRepository(client)
	googleCalendarRepo := googlecalendarinfo.NewGoogleCalendarInfoRepository(client)
	eventRepo := event.NewEventRepository(client)
	dateRepo := proposeddate.NewProposedDateRepository(client)
	calendarApp := appCalendar.NewGoogleCalendarManager(client) // Google Calendar API manager
	keyManager := auth.NewKeyManager(client, cache)
	jwtManager := auth.NewJWTManager(client, keyManager)
	authManager := auth.NewAuthManager(client, userRepo, oauthRepo)

	eventManager := appEvents.NewEventManager(client, authManager, calendarRepo, googleCalendarRepo, eventRepo, dateRepo, calendarApp)

	return &Server{
		Client:        			  client,
		Cache:         			  cache,
		UserRepo:      			  userRepo,
		OAuthRepo:     			  oauthRepo,
		CalendarRepo:  			  calendarRepo,
		GoogleCalendarRepo:      googleCalendarRepo,
		EventRepo:     			  eventRepo,
		DateRepo:      			  dateRepo,
		AuthManager:   			  authManager,
		JWTManager:    			  jwtManager,
		KeyManager:    			  keyManager,
		EventManager:  			  eventManager,
		EventCreationManager:     event_operations.NewEventCreationManager(eventManager),
		EventFetchingManager:     event_operations.NewEventFetchingManager(eventManager),
		EventUpdateManager:       event_operations.NewEventUpdateManager(eventManager),
		EventDeleteManager:       event_operations.NewEventDeleteManager(eventManager),
	}
}
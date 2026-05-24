package events

import (
	"github.com/koo-arch/adjusta-backend/ent"
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	googleOAuth "github.com/koo-arch/adjusta-backend/internal/google/oauth"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
)

type Usecase struct {
	client             *ent.Client
	googleTokenManager *googleOAuth.TokenManager
	calendarRepo       repoCalendar.CalendarRepository
	googleCalendarRepo googlecalendarinfo.GoogleCalendarInfoRepository
	eventRepo          event.EventRepository
	dateRepo           proposeddate.ProposedDateRepository
	calendarApp        *appCalendar.GoogleCalendarManager
}

func NewUsecase(
	client *ent.Client,
	googleTokenManager *googleOAuth.TokenManager,
	calendarRepo repoCalendar.CalendarRepository,
	googleCalendarRepo googlecalendarinfo.GoogleCalendarInfoRepository,
	eventRepo event.EventRepository,
	dateRepo proposeddate.ProposedDateRepository,
	calendarApp *appCalendar.GoogleCalendarManager,
) *Usecase {
	return &Usecase{
		client:             client,
		googleTokenManager: googleTokenManager,
		calendarRepo:       calendarRepo,
		googleCalendarRepo: googleCalendarRepo,
		eventRepo:          eventRepo,
		dateRepo:           dateRepo,
		calendarApp:        calendarApp,
	}
}

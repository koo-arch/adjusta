package calendar

import (
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
)

type SyncTxRepositories struct {
	Calendar     repoCalendar.CalendarRepository
	UserCalendar repoUserCalendar.UserCalendarRepository
}

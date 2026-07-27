package events

import (
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoOutboxMessage "github.com/koo-arch/adjusta-backend/internal/domain/outboxmessage"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
)

type EventTxRepositories struct {
	Calendar      repoCalendar.CalendarRepository
	Event         repoEvent.EventRepository
	OutboxMessage repoOutboxMessage.Repository
	ProposedDate  repoProposedDate.ProposedDateRepository
	UserCalendar  repoUserCalendar.UserCalendarRepository
}

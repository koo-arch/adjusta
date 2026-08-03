package repository

import (
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoOutboxMessage "github.com/koo-arch/adjusta-backend/internal/domain/outboxmessage"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"
	infraAccount "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/account"
	infraCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/calendar"
	infraEvent "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/event"
	infraOutboxMessage "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/outboxmessage"
	infraProposedDate "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/proposeddate"
	infraSession "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/session"
	infraUser "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/user"
	infraUserCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/usercalendar"
)

type Repositories struct {
	Account       repoAccount.AccountRepository
	Calendar      repoCalendar.CalendarRepository
	Event         repoEvent.EventRepository
	OutboxMessage repoOutboxMessage.Repository
	ProposedDate  repoProposedDate.ProposedDateRepository
	Session       repoSession.SessionRepository
	User          repoUser.UserRepository
	UserCalendar  repoUserCalendar.UserCalendarRepository
}

func NewRepositories(client *ent.Client) Repositories {
	return Repositories{
		Account:       infraAccount.NewAccountRepository(client),
		Calendar:      infraCalendar.NewCalendarRepository(client),
		Event:         infraEvent.NewEventRepository(client),
		OutboxMessage: infraOutboxMessage.NewRepository(client),
		ProposedDate:  infraProposedDate.NewProposedDateRepository(client),
		Session:       infraSession.NewSessionRepository(client),
		User:          infraUser.NewUserRepository(client),
		UserCalendar:  infraUserCalendar.NewUserCalendarRepository(client),
	}
}

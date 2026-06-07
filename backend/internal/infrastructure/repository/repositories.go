package repository

import (
	"github.com/koo-arch/adjusta-backend/ent"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/domain/account"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	infraAccount "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/account"
	infraCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/calendar"
	infraEvent "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/event"
	infraProposedDate "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/proposeddate"
	infraSession "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/session"
	infraUser "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/user"
	infraUserCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type Repositories struct {
	Account      repoAccount.AccountRepository
	Calendar     repoCalendar.CalendarRepository
	Event        repoEvent.EventRepository
	ProposedDate repoProposedDate.ProposedDateRepository
	Session      repoSession.SessionRepository
	User         repoUser.UserRepository
	UserCalendar repoUserCalendar.UserCalendarRepository
}

func NewRepositories(client *ent.Client) Repositories {
	return Repositories{
		Account:      infraAccount.NewAccountRepository(client),
		Calendar:     infraCalendar.NewCalendarRepository(client),
		Event:        infraEvent.NewEventRepository(client),
		ProposedDate: infraProposedDate.NewProposedDateRepository(client),
		Session:      infraSession.NewSessionRepository(client),
		User:         infraUser.NewUserRepository(client),
		UserCalendar: infraUserCalendar.NewUserCalendarRepository(client),
	}
}

func (r Repositories) WithTx(tx transaction.Tx) Repositories {
	return Repositories{
		Account:      r.Account.WithTx(tx),
		Calendar:     r.Calendar.WithTx(tx),
		Event:        r.Event.WithTx(tx),
		ProposedDate: r.ProposedDate.WithTx(tx),
		Session:      r.Session.WithTx(tx),
		User:         r.User.WithTx(tx),
		UserCalendar: r.UserCalendar.WithTx(tx),
	}
}

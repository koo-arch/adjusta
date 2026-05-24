package repo

import (
	"github.com/koo-arch/adjusta-backend/ent"
	repoAccount "github.com/koo-arch/adjusta-backend/internal/repo/account"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	repoEvent "github.com/koo-arch/adjusta-backend/internal/repo/event"
	repoGoogleCalendarInfo "github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	repoSession "github.com/koo-arch/adjusta-backend/internal/repo/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/repo/user"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type Repositories struct {
	Account            repoAccount.AccountRepository
	Calendar           repoCalendar.CalendarRepository
	Event              repoEvent.EventRepository
	GoogleCalendarInfo repoGoogleCalendarInfo.GoogleCalendarInfoRepository
	ProposedDate       repoProposedDate.ProposedDateRepository
	Session            repoSession.SessionRepository
	User               repoUser.UserRepository
}

func NewRepositories(client *ent.Client) Repositories {
	return Repositories{
		Account:            repoAccount.NewAccountRepository(client),
		Calendar:           repoCalendar.NewCalendarRepository(client),
		Event:              repoEvent.NewEventRepository(client),
		GoogleCalendarInfo: repoGoogleCalendarInfo.NewGoogleCalendarInfoRepository(client),
		ProposedDate:       repoProposedDate.NewProposedDateRepository(client),
		Session:            repoSession.NewSessionRepository(client),
		User:               repoUser.NewUserRepository(client),
	}
}

func (r Repositories) WithTx(tx transaction.Tx) Repositories {
	return Repositories{
		Account:            r.Account.WithTx(tx),
		Calendar:           r.Calendar.WithTx(tx),
		Event:              r.Event.WithTx(tx),
		GoogleCalendarInfo: r.GoogleCalendarInfo.WithTx(tx),
		ProposedDate:       r.ProposedDate.WithTx(tx),
		Session:            r.Session.WithTx(tx),
		User:               r.User.WithTx(tx),
	}
}

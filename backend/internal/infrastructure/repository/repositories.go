package repository

import (
	"github.com/koo-arch/adjusta-backend/ent"
	infraAccount "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/account"
	infraCalendar "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/calendar"
	infraEvent "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/event"
	infraGoogleCalendarInfo "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/googlecalendarinfo"
	infraProposedDate "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/proposeddate"
	infraSession "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/session"
	infraUser "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/user"
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
		Account:            infraAccount.NewAccountRepository(client),
		Calendar:           infraCalendar.NewCalendarRepository(client),
		Event:              infraEvent.NewEventRepository(client),
		GoogleCalendarInfo: infraGoogleCalendarInfo.NewGoogleCalendarInfoRepository(client),
		ProposedDate:       infraProposedDate.NewProposedDateRepository(client),
		Session:            infraSession.NewSessionRepository(client),
		User:               infraUser.NewUserRepository(client),
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

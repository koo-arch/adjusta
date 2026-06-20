package events

type Usecase struct {
	repos          EventRepositories
	tx             EventTransaction
	googleCalendar GoogleCalendarGateway
}

func NewUsecase(
	repos EventRepositories,
	tx EventTransaction,
	googleCalendar GoogleCalendarGateway,
) *Usecase {
	return &Usecase{
		repos:          repos,
		tx:             tx,
		googleCalendar: googleCalendar,
	}
}

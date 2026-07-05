package events

type Usecase struct {
	repos          EventTxRepositories
	tx             EventTransaction
	googleCalendar GoogleCalendarGateway
}

func NewUsecase(
	repos EventTxRepositories,
	tx EventTransaction,
	googleCalendar GoogleCalendarGateway,
) *Usecase {
	return &Usecase{
		repos:          repos,
		tx:             tx,
		googleCalendar: googleCalendar,
	}
}

package events

type Usecase struct {
	reader         EventReader
	tx             EventTransaction
	googleCalendar GoogleCalendarGateway
}

func NewUsecase(
	reader EventReader,
	tx EventTransaction,
	googleCalendar GoogleCalendarGateway,
) *Usecase {
	return &Usecase{
		reader:         reader,
		tx:             tx,
		googleCalendar: googleCalendar,
	}
}

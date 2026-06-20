package googlecalendar

type FetchedEvent struct {
	ID          string
	Summary     string
	Description string
	Location    string
	ColorID     string
	Start       string
	End         string
}

package value

type EventStatus string

const (
	StatusDraft     EventStatus = "draft"
	StatusActive    EventStatus = "active"
	StatusConfirmed EventStatus = "confirmed"
	StatusCancelled EventStatus = "cancelled"
)

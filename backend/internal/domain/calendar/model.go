package calendar

import "github.com/google/uuid"

type Calendar struct {
	ID               uuid.UUID
	GoogleCalendarID string
	Summary          string
	Description      *string
	Timezone         *string
}

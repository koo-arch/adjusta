package repositorymodel

import "github.com/google/uuid"

type StoredCalendar struct {
	ID               uuid.UUID
	GoogleCalendarID string
	Summary          string
	Description      *string
	Timezone         *string
}

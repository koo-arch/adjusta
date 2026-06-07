package repositorymodel

import (
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

type UserCalendar struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	CalendarID        uuid.UUID
	Role              domainvalue.UserCalendarRole
	IsVisible         bool
	SyncProposedDates bool
}

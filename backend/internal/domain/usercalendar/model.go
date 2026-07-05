package usercalendar

import (
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type UserCalendar struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	CalendarID        uuid.UUID
	Role              value.UserCalendarRole
	IsVisible         bool
	SyncProposedDates bool
}

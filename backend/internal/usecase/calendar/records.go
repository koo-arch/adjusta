package calendar

import (
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type CalendarRecord struct {
	CalendarID string
	Summary    string
	Primary    bool
}

type UserCalendarRelationRecord struct {
	CalendarID        uuid.UUID
	GoogleCalendarID  string
	Role              value.UserCalendarRole
	SyncProposedDates bool
}

package usercalendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

type UserCalendarQueryOptions struct {
	Role              *domainvalue.UserCalendarRole
	IsVisible         *bool
	SyncProposedDates *bool
}

type UserCalendarRepository interface {
	FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*UserCalendar, error)
	Ensure(ctx context.Context, userID, calendarID uuid.UUID, opt UserCalendarQueryOptions) (*UserCalendar, error)
	SoftDeleteByUserAndCalendar(ctx context.Context, userID, calendarID uuid.UUID) error
}

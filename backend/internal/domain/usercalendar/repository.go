package usercalendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type UserCalendarQueryOptions struct {
	Role              *value.UserCalendarRole
	IsVisible         *bool
	SyncProposedDates *bool
}

type UserCalendarRepository interface {
	FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*UserCalendar, error)
	Ensure(ctx context.Context, userID, calendarID uuid.UUID, opt UserCalendarQueryOptions) (*UserCalendar, error)
	SoftDeleteByUserAndCalendar(ctx context.Context, userID, calendarID uuid.UUID) error
}

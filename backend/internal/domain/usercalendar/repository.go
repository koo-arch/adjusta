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
	FindByIDAndUser(ctx context.Context, userID, id uuid.UUID) (*UserCalendar, error)
	FindByRole(ctx context.Context, userID uuid.UUID, role value.UserCalendarRole) (*UserCalendar, error)
	FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*UserCalendar, error)
	Ensure(ctx context.Context, userID, calendarID uuid.UUID, opt UserCalendarQueryOptions) (*UserCalendar, error)
	Update(ctx context.Context, userID, id uuid.UUID, opt UserCalendarQueryOptions) (*UserCalendar, error)
	SoftDeleteByUserAndCalendar(ctx context.Context, userID, calendarID uuid.UUID) error
}

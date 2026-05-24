package calendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type CalendarQueryOptions struct {
	GoogleCalendarID       *string
	Summary                *string
	IsPrimary              *bool
	WithGoogleCalendarInfo bool `json:"with_google_calendar_info"`
	WithEvents             bool `json:"with_events"`
	WithProposedDates      bool `json:"with_proposed_dates"`
	EventOffset            int
	EventLimit             int
	ProposedDateOffset     int
	ProposedDateLimit      int
}

type CalendarRepository interface {
	WithTx(tx transaction.Tx) CalendarRepository
	Read(ctx context.Context, id uuid.UUID, opt CalendarQueryOptions) (*models.StoredCalendar, error)
	FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*models.StoredCalendar, error)
	FindByFields(ctx context.Context, userID uuid.UUID, opt CalendarQueryOptions) (*models.StoredCalendar, error)
	FilterByFields(ctx context.Context, userID uuid.UUID, opt CalendarQueryOptions) ([]*models.StoredCalendar, error)
	Create(ctx context.Context, userID uuid.UUID) (*models.StoredCalendar, error)
	Update(ctx context.Context, id uuid.UUID) (*models.StoredCalendar, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}

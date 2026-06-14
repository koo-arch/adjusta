package calendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type CalendarQueryOptions struct {
	GoogleCalendarID   *string
	Summary            *string
	Role               *domainvalue.UserCalendarRole
	WithEvents         bool `json:"with_events"`
	WithProposedDates  bool `json:"with_proposed_dates"`
	EventOffset        int
	EventLimit         int
	ProposedDateOffset int
	ProposedDateLimit  int
}

type CalendarMutationOptions struct {
	GoogleCalendarID *string
	Summary          *string
	Description      *string
	Timezone         *string
}

type CalendarRepository interface {
	WithTx(tx transaction.Tx) CalendarRepository
	Read(ctx context.Context, id uuid.UUID, opt CalendarQueryOptions) (*Calendar, error)
	FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*Calendar, error)
	FindByFields(ctx context.Context, userID uuid.UUID, opt CalendarQueryOptions) (*Calendar, error)
	FindByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*Calendar, error)
	FilterByFields(ctx context.Context, userID uuid.UUID, opt CalendarQueryOptions) ([]*Calendar, error)
	Create(ctx context.Context, opt CalendarMutationOptions) (*Calendar, error)
	Update(ctx context.Context, id uuid.UUID, opt CalendarMutationOptions) (*Calendar, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}

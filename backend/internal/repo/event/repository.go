package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"google.golang.org/api/calendar/v3"
)

type EventQueryOptions struct {
	Summary              *string
	Location             *string
	Description          *string
	Status               *domainvalue.EventStatus
	ConfirmedDateID      *uuid.UUID
	GoogleEventID        *string
	Slug                 *string
	WithProposedDates    bool
	EventOffset          int
	EventLimit           int
	ProposedDateOffset   int
	ProposedDateLimit    int
	ProposedDateStartGTE *time.Time
	ProposedDateStartLTE *time.Time
	ProposedDateEndGTE   *time.Time
	ProposedDateEndLTE   *time.Time
	SortBy               string
	SortOrder            string
}

type EventRepository interface {
	WithTx(tx transaction.Tx) EventRepository
	Read(ctx context.Context, id uuid.UUID, opt EventQueryOptions) (*repositorymodel.StoredEvent, error)
	FilterByCalendarID(ctx context.Context, calendarID uuid.UUID, opt EventQueryOptions) ([]*repositorymodel.StoredEvent, error)
	FindBySlugAndUser(ctx context.Context, userID uuid.UUID, slug string, opt EventQueryOptions) (*repositorymodel.StoredEvent, error)
	Create(ctx context.Context, googleEvent *calendar.Event, calendarID uuid.UUID) (*repositorymodel.StoredEvent, error)
	Update(ctx context.Context, id uuid.UUID, opt EventQueryOptions) (*repositorymodel.StoredEvent, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	SearchEvents(ctx context.Context, id, calendarID uuid.UUID, opt EventQueryOptions) ([]*repositorymodel.StoredEvent, error)
}

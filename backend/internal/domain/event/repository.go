package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type EventCreateOptions struct {
	Title       string
	Location    string
	Description string
}

type EventQueryOptions struct {
	Title                  *string
	Location               *string
	Description            *string
	Status                 *domainvalue.EventStatus
	SyncStatus             *domainvalue.SyncStatus
	ConfirmedDateID        *uuid.UUID
	ConfirmedGoogleEventID *string
	LastSyncedAt           *time.Time
	ClearLastSyncedAt      bool
	LastSyncError          *string
	ClearLastSyncError     bool
	Slug                   *string
	WithProposedDates      bool
	EventOffset            int
	EventLimit             int
	ProposedDateOffset     int
	ProposedDateLimit      int
	ProposedDateStartGTE   *time.Time
	ProposedDateStartLTE   *time.Time
	ProposedDateEndGTE     *time.Time
	ProposedDateEndLTE     *time.Time
	SortBy                 string
	SortOrder              string
}

type EventRepository interface {
	WithTx(tx transaction.Tx) EventRepository
	Read(ctx context.Context, id uuid.UUID, opt EventQueryOptions) (*Event, error)
	FilterByCalendarID(ctx context.Context, calendarID uuid.UUID, opt EventQueryOptions) ([]*Event, error)
	FindBySlugAndUser(ctx context.Context, userID uuid.UUID, slug string, opt EventQueryOptions) (*Event, error)
	Create(ctx context.Context, userID uuid.UUID, opt EventCreateOptions, primaryCalendarID uuid.UUID) (*Event, error)
	Update(ctx context.Context, id uuid.UUID, opt EventQueryOptions) (*Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	SearchEvents(ctx context.Context, id, calendarID uuid.UUID, opt EventQueryOptions) ([]*Event, error)
}

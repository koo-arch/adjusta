package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"google.golang.org/api/calendar/v3"
)

type EventQueryOptions struct {
	Summary              *string
	Location             *string
	Description          *string
	Status               *models.EventStatus
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
	Read(ctx context.Context, id uuid.UUID, opt EventQueryOptions) (*models.StoredEvent, error)
	FilterByCalendarID(ctx context.Context, calendarID uuid.UUID, opt EventQueryOptions) ([]*models.StoredEvent, error)
	FindBySlugAndUser(ctx context.Context, userID uuid.UUID, slug string, opt EventQueryOptions) (*models.StoredEvent, error)
	Create(ctx context.Context, googleEvent *calendar.Event, calendarID uuid.UUID) (*models.StoredEvent, error)
	Update(ctx context.Context, id uuid.UUID, opt EventQueryOptions) (*models.StoredEvent, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	SearchEvents(ctx context.Context, id, calendarID uuid.UUID, opt EventQueryOptions) ([]*models.StoredEvent, error)
}

package proposeddate

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type ProposedDateCreateOptions struct {
	StartTime     time.Time
	EndTime       time.Time
	Priority      int
	GoogleEventID *string
	Status        *value.ProposedDateStatus
	SyncStatus    *value.SyncStatus
	LastSyncedAt  *time.Time
	LastSyncError *string
}

type ProposedDateUpdateOptions struct {
	GoogleEventID      *string
	StartTime          *time.Time
	EndTime            *time.Time
	Priority           *int
	Status             *value.ProposedDateStatus
	SyncStatus         *value.SyncStatus
	LastSyncedAt       *time.Time
	ClearLastSyncedAt  bool
	LastSyncError      *string
	ClearLastSyncError bool
}

type ProposedDateRepository interface {
	Read(ctx context.Context, id uuid.UUID) (*ProposedDate, error)
	FilterByEventID(ctx context.Context, eventID uuid.UUID) ([]*ProposedDate, error)
	Create(ctx context.Context, opt ProposedDateCreateOptions, eventID uuid.UUID) (*ProposedDate, error)
	Update(ctx context.Context, id uuid.UUID, opt ProposedDateUpdateOptions) (*ProposedDate, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	CreateBulk(ctx context.Context, opts []ProposedDateCreateOptions, eventID uuid.UUID) ([]*ProposedDate, error)
}

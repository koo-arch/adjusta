package proposeddate

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type ProposedDateQueryOptions struct {
	GoogleEventID      *string
	StartTime          *time.Time
	EndTime            *time.Time
	Priority           *int
	Status             *domainvalue.ProposedDateStatus
	SyncStatus         *domainvalue.SyncStatus
	LastSyncedAt       *time.Time
	ClearLastSyncedAt  bool
	LastSyncError      *string
	ClearLastSyncError bool
}

type ProposedDateRepository interface {
	WithTx(tx transaction.Tx) ProposedDateRepository
	Read(ctx context.Context, id uuid.UUID) (*repositorymodel.StoredProposedDate, error)
	FilterByEventID(ctx context.Context, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error)
	ExclusionEventID(ctx context.Context, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error)
	Create(ctx context.Context, opt ProposedDateQueryOptions, eventID uuid.UUID) (*repositorymodel.StoredProposedDate, error)
	Update(ctx context.Context, id uuid.UUID, opt ProposedDateQueryOptions) (*repositorymodel.StoredProposedDate, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	CreateBulk(ctx context.Context, selectedDates []appmodel.SelectedDate, eventID uuid.UUID) ([]*repositorymodel.StoredProposedDate, error)
}

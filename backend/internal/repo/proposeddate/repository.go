package proposeddate

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type ProposedDateQueryOptions struct {
	StartTime *time.Time
	EndTime   *time.Time
	Priority  *int
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
	DecrementPriorityExceptID(ctx context.Context, eventID, excludeID uuid.UUID) error
	ReorderPriority(ctx context.Context, eventID uuid.UUID) error
}

package googlecalendarinfo

import (
	"context"

	"github.com/google/uuid"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type GoogleCalendarInfoQueryOptions struct {
	GoogleCalendarID *string
	Summary          *string
	IsPrimary        *bool
}

type GoogleCalendarInfoRepository interface {
	WithTx(tx transaction.Tx) GoogleCalendarInfoRepository
	Read(ctx context.Context, id uuid.UUID) (*repositorymodel.GoogleCalendarInfo, error)
	FindByFields(ctx context.Context, opt GoogleCalendarInfoQueryOptions) (*repositorymodel.GoogleCalendarInfo, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*repositorymodel.GoogleCalendarInfo, error)
	Create(ctx context.Context, opt GoogleCalendarInfoQueryOptions, calendarID uuid.UUID) (*repositorymodel.GoogleCalendarInfo, error)
	Update(ctx context.Context, id uuid.UUID, opt GoogleCalendarInfoQueryOptions, calendarID *uuid.UUID) (*repositorymodel.GoogleCalendarInfo, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}

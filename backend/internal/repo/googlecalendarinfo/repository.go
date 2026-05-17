package googlecalendarinfo

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
)

type GoogleCalendarInfoQueryOptions struct {
	GoogleCalendarID *string
	Summary    *string
	IsPrimary  *bool
}

type GoogleCalendarInfoRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.GoogleCalendarInfo, error)
	FindByFields(ctx context.Context, tx *ent.Tx, opt GoogleCalendarInfoQueryOptions) (*ent.GoogleCalendarInfo, error)
	ListByUser(ctx context.Context, tx *ent.Tx, userID uuid.UUID) ([]*ent.GoogleCalendarInfo, error)
	Create(ctx context.Context, tx *ent.Tx, opt GoogleCalendarInfoQueryOptions, entCalendar *ent.Calendar) (*ent.GoogleCalendarInfo, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt GoogleCalendarInfoQueryOptions, entCalendar *ent.Calendar) (*ent.GoogleCalendarInfo, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
}
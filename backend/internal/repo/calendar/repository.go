package calendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
)

type CalendarQueryOptions struct {
	GoogleCalendarID *string
    Summary    *string
    IsPrimary  *bool
	WithGoogleCalendarInfo bool `json:"with_google_calendar_info"`
	WithEvents bool `json:"with_events"`
	WithProposedDates bool `json:"with_proposed_dates"`
	EventOffset int
	EventLimit int
	ProposedDateOffset int
	ProposedDateLimit int
}

type CalendarRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt CalendarQueryOptions) (*ent.Calendar, error)
	FilterByUserID(ctx context.Context, tx *ent.Tx, userID uuid.UUID) ([]*ent.Calendar, error)
	FindByFields(ctx context.Context, tx *ent.Tx, userID uuid.UUID, opt CalendarQueryOptions) (*ent.Calendar, error)
	FilterByFields(ctx context.Context, tx *ent.Tx, userID uuid.UUID, opt CalendarQueryOptions) ([]*ent.Calendar, error)
	Create(ctx context.Context, tx *ent.Tx, entUser *ent.User, entGoogleCalendar *ent.GoogleCalendarInfo) (*ent.Calendar, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.Calendar, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	SoftDeleteWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	RestoreWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
}
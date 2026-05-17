package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type EventQueryOptions struct {
	Summary		   *string
	Location	   *string
	Description	   *string
	Status 		     *models.EventStatus
	ConfirmedDateID *uuid.UUID
	GoogleEventID   *string
	Slug 		     *string
	WithProposedDates bool
	EventOffset       int
	EventLimit        int
	ProposedDateOffset int
	ProposedDateLimit  int
	ProposedDateStartGTE *time.Time
	ProposedDateStartLTE  *time.Time
	ProposedDateEndGTE  *time.Time
	ProposedDateEndLTE  *time.Time
	SortBy 		     	string
	SortOrder 	    	string
}

type EventRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt EventQueryOptions) (*ent.Event, error)
	FilterByCalendarID(ctx context.Context, tx *ent.Tx, calendarID uuid.UUID, opt EventQueryOptions) ([]*ent.Event, error)
	FindBySlugAndUser(ctx context.Context, tx *ent.Tx, userID uuid.UUID, slug string, opt EventQueryOptions) (*ent.Event, error)
	Create(ctx context.Context, tx *ent.Tx, googleEvent *calendar.Event, entCalendar *ent.Calendar) (*ent.Event, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt EventQueryOptions) (*ent.Event, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	SoftDeleteWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	RestoreWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	SearchEvents(ctx context.Context, tx *ent.Tx, id, calendarID uuid.UUID, opt EventQueryOptions) ([]*ent.Event, error)
}
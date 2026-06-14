package event

import (
	"time"

	"github.com/google/uuid"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

type Event struct {
	UserID                 uuid.UUID
	ID                     uuid.UUID
	PrimaryCalendarID      uuid.UUID
	Title                  string
	Location               string
	Description            string
	Status                 domainvalue.EventStatus
	ConfirmedDateID        uuid.UUID
	GoogleEventID          string
	ConfirmedGoogleEventID *string
	SyncStatus             domainvalue.SyncStatus
	LastSyncedAt           *time.Time
	LastSyncError          *string
	Slug                   string
	ProposedDates          []*repoProposedDate.ProposedDate
}

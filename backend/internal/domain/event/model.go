package event

import (
	"time"

	"github.com/google/uuid"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type Event struct {
	UserID                 uuid.UUID
	ID                     uuid.UUID
	PrimaryCalendarID      uuid.UUID
	Title                  string
	Location               string
	Description            string
	Status                 value.EventStatus
	ConfirmedDateID        uuid.UUID
	ConfirmedGoogleEventID *string
	SyncStatus             value.SyncStatus
	LastSyncedAt           *time.Time
	LastSyncError          *string
	ProposedDates          []*repoProposedDate.ProposedDate
}

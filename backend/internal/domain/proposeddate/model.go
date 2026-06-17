package proposeddate

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

type ProposedDate struct {
	ID            uuid.UUID
	EventID       uuid.UUID
	GoogleEventID *string
	StartTime     time.Time
	EndTime       time.Time
	Priority      int
	Status        domainvalue.ProposedDateStatus
	SyncStatus    domainvalue.SyncStatus
	LastSyncedAt  *time.Time
	LastSyncError *string
}

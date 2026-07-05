package proposeddate

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type ProposedDate struct {
	ID            uuid.UUID
	EventID       uuid.UUID
	GoogleEventID *string
	StartTime     time.Time
	EndTime       time.Time
	Priority      int
	Status        value.ProposedDateStatus
	SyncStatus    value.SyncStatus
	LastSyncedAt  *time.Time
	LastSyncError *string
}

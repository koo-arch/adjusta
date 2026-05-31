package repositorymodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

type StoredProposedDate struct {
	ID        uuid.UUID
	EventID   *uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	Priority  int
}

type StoredEvent struct {
	ID              uuid.UUID
	Summary         string
	Location        string
	Description     string
	Status          domainvalue.EventStatus
	ConfirmedDateID uuid.UUID
	GoogleEventID   string
	Slug            string
	ProposedDates   []*StoredProposedDate
}

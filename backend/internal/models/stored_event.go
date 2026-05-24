package models

import (
	"time"

	"github.com/google/uuid"
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
	Status          EventStatus
	ConfirmedDateID uuid.UUID
	GoogleEventID   string
	Slug            string
	ProposedDates   []*StoredProposedDate
}

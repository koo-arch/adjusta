package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type SelectedDate struct {
	Start    time.Time
	End      time.Time
	Priority int
}

type DraftCreationRequest struct {
	Title         string
	Location      string
	Description   string
	SelectedDates []SelectedDate
}

type ProposedDateRequest struct {
	ID            *uuid.UUID
	GoogleEventID *string
	Start         *time.Time
	End           *time.Time
	Priority      int
}

type DraftUpdateRequest struct {
	Title         string
	Location      string
	Description   string
	Status        value.EventStatus
	ProposedDates []ProposedDateRequest
}

type ConfirmationRequest struct {
	ID            *uuid.UUID
	GoogleEventID string
	Start         *time.Time
	End           *time.Time
	Priority      int
}

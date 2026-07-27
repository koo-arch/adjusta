package outboxmessage

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID                uuid.UUID
	MessageType       string
	AggregateType     string
	AggregateID       uuid.UUID
	Payload           json.RawMessage
	DispatchAttempts  int
	AvailableAt       time.Time
	DispatchedAt      *time.Time
	ProcessedAt       *time.Time
	LastDispatchError *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

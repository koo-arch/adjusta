package outboxmessage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateOptions struct {
	MessageType   string
	AggregateType string
	AggregateID   uuid.UUID
	Payload       json.RawMessage
	AvailableAt   time.Time
}

type Repository interface {
	Create(ctx context.Context, opt CreateOptions) (*Message, error)
	Read(ctx context.Context, id uuid.UUID) (*Message, error)
	ListDispatchable(ctx context.Context, now time.Time, limit int) ([]*Message, error)
	MarkDispatched(ctx context.Context, id uuid.UUID, dispatchedAt time.Time) error
	RecordDispatchFailure(ctx context.Context, id uuid.UUID, message string) error
	MarkProcessed(ctx context.Context, id uuid.UUID, processedAt time.Time) error
}

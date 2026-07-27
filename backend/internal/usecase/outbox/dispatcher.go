package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/outboxmessage"
)

type Publisher interface {
	Publish(ctx context.Context, messageID uuid.UUID) error
}

type Dispatcher struct {
	messages  outboxmessage.Repository
	publisher Publisher
}

func NewDispatcher(messages outboxmessage.Repository, publisher Publisher) *Dispatcher {
	return &Dispatcher{messages: messages, publisher: publisher}
}

func (d *Dispatcher) DispatchPending(ctx context.Context, limit int) error {
	messages, err := d.messages.ListDispatchable(ctx, time.Now(), limit)
	if err != nil {
		return err
	}

	var dispatchErrors []error
	for _, message := range messages {
		if message == nil {
			continue
		}
		if err := d.Dispatch(ctx, message.ID); err != nil {
			dispatchErrors = append(dispatchErrors, err)
		}
	}
	return errors.Join(dispatchErrors...)
}

func (d *Dispatcher) Dispatch(ctx context.Context, messageID uuid.UUID) error {
	if err := d.publisher.Publish(ctx, messageID); err != nil {
		if recordErr := d.messages.RecordDispatchFailure(ctx, messageID, err.Error()); recordErr != nil {
			return errors.Join(
				fmt.Errorf("dispatch message %s: %w", messageID, err),
				fmt.Errorf("record dispatch failure for %s: %w", messageID, recordErr),
			)
		}
		return fmt.Errorf("dispatch message %s: %w", messageID, err)
	}
	if err := d.messages.MarkDispatched(ctx, messageID, time.Now()); err != nil {
		return fmt.Errorf("mark message %s dispatched: %w", messageID, err)
	}
	return nil
}

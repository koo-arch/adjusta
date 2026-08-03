package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/outboxmessage"
)

func enqueueGoogleCalendarSync(ctx context.Context, repos EventTxRepositories, eventID uuid.UUID) (uuid.UUID, error) {
	message, err := repos.OutboxMessage.Create(ctx, outboxmessage.CreateOptions{
		MessageType:   outboxmessage.MessageTypeCalendarEventSync,
		AggregateType: outboxmessage.AggregateTypeEvent,
		AggregateID:   eventID,
		Payload:       json.RawMessage(`{}`),
		AvailableAt:   time.Now(),
	})
	if err != nil {
		return uuid.Nil, err
	}
	return message.ID, nil
}

func (uc *Usecase) dispatchOutboxMessage(ctx context.Context, messageID uuid.UUID) {
	if uc.outboxDispatcher == nil || messageID == uuid.Nil {
		return
	}
	if err := uc.outboxDispatcher.Dispatch(ctx, messageID); err != nil {
		// The recovery dispatcher retries rows whose dispatched_at remains NULL.
		log.Printf("failed to dispatch outbox message %s: %v", messageID, err)
		return
	}
}

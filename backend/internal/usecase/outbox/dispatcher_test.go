package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/outboxmessage"
)

type fakeMessageRepository struct {
	messages        []*outboxmessage.Message
	dispatched      []uuid.UUID
	dispatchFailure []uuid.UUID
}

func (f *fakeMessageRepository) Create(context.Context, outboxmessage.CreateOptions) (*outboxmessage.Message, error) {
	return nil, errors.New("Create should not be called")
}
func (f *fakeMessageRepository) Read(context.Context, uuid.UUID) (*outboxmessage.Message, error) {
	return nil, errors.New("Read should not be called")
}
func (f *fakeMessageRepository) ListDispatchable(context.Context, time.Time, int) ([]*outboxmessage.Message, error) {
	return f.messages, nil
}
func (f *fakeMessageRepository) MarkDispatched(_ context.Context, id uuid.UUID, _ time.Time) error {
	f.dispatched = append(f.dispatched, id)
	return nil
}
func (f *fakeMessageRepository) RecordDispatchFailure(_ context.Context, id uuid.UUID, _ string) error {
	f.dispatchFailure = append(f.dispatchFailure, id)
	return nil
}
func (f *fakeMessageRepository) MarkProcessed(context.Context, uuid.UUID, time.Time) error {
	return errors.New("MarkProcessed should not be called")
}

type fakePublisher struct {
	published []uuid.UUID
	failID    uuid.UUID
}

func (f *fakePublisher) Publish(_ context.Context, id uuid.UUID) error {
	f.published = append(f.published, id)
	if id == f.failID {
		return errors.New("cloud tasks unavailable")
	}
	return nil
}

func TestDispatcherRecordsSuccessAndLeavesFailuresForRecovery(t *testing.T) {
	successID := uuid.New()
	failureID := uuid.New()
	repository := &fakeMessageRepository{messages: []*outboxmessage.Message{
		{ID: successID, Payload: json.RawMessage(`{}`)},
		{ID: failureID, Payload: json.RawMessage(`{}`)},
	}}
	publisher := &fakePublisher{failID: failureID}
	dispatcher := NewDispatcher(repository, publisher)

	err := dispatcher.DispatchPending(context.Background(), 100)
	if err == nil {
		t.Fatal("expected aggregate dispatch error")
	}
	if len(repository.dispatched) != 1 || repository.dispatched[0] != successID {
		t.Fatalf("unexpected dispatched messages: %v", repository.dispatched)
	}
	if len(repository.dispatchFailure) != 1 || repository.dispatchFailure[0] != failureID {
		t.Fatalf("unexpected failed messages: %v", repository.dispatchFailure)
	}
}

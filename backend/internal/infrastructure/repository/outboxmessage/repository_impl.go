package outboxmessage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	domain "github.com/koo-arch/adjusta-backend/internal/domain/outboxmessage"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/outboxmessage"
	infraerr "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository/infraerr"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) Create(ctx context.Context, opt domain.CreateOptions) (*domain.Message, error) {
	payload := map[string]any{}
	if len(opt.Payload) > 0 {
		if err := json.Unmarshal(opt.Payload, &payload); err != nil {
			return nil, err
		}
	}

	entity, err := r.client.OutboxMessage.Create().
		SetMessageType(opt.MessageType).
		SetAggregateType(opt.AggregateType).
		SetAggregateID(opt.AggregateID).
		SetPayload(payload).
		SetAvailableAt(opt.AvailableAt).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDomain(entity)
}

func (r *Repository) Read(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	entity, err := r.client.OutboxMessage.Query().
		Where(outboxmessage.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toDomain(entity)
}

func (r *Repository) ListDispatchable(ctx context.Context, now time.Time, limit int) ([]*domain.Message, error) {
	query := r.client.OutboxMessage.Query().
		Where(
			outboxmessage.DispatchedAtIsNil(),
			outboxmessage.ProcessedAtIsNil(),
			outboxmessage.AvailableAtLTE(now),
		).
		Order(ent.Asc(outboxmessage.FieldCreatedAt))
	if limit > 0 {
		query = query.Limit(limit)
	}

	entities, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Message, 0, len(entities))
	for _, entity := range entities {
		message, err := toDomain(entity)
		if err != nil {
			return nil, err
		}
		result = append(result, message)
	}
	return result, nil
}

func (r *Repository) MarkDispatched(ctx context.Context, id uuid.UUID, dispatchedAt time.Time) error {
	_, err := r.client.OutboxMessage.UpdateOneID(id).
		SetDispatchedAt(dispatchedAt).
		ClearLastDispatchError().
		Save(ctx)
	return infraerr.MapNotFound(err)
}

func (r *Repository) RecordDispatchFailure(ctx context.Context, id uuid.UUID, message string) error {
	_, err := r.client.OutboxMessage.UpdateOneID(id).
		AddDispatchAttempts(1).
		SetLastDispatchError(message).
		Save(ctx)
	return infraerr.MapNotFound(err)
}

func (r *Repository) MarkProcessed(ctx context.Context, id uuid.UUID, processedAt time.Time) error {
	_, err := r.client.OutboxMessage.UpdateOneID(id).
		SetProcessedAt(processedAt).
		Save(ctx)
	return infraerr.MapNotFound(err)
}

func toDomain(entity *ent.OutboxMessage) (*domain.Message, error) {
	payload, err := json.Marshal(entity.Payload)
	if err != nil {
		return nil, err
	}
	return &domain.Message{
		ID:                entity.ID,
		MessageType:       entity.MessageType,
		AggregateType:     entity.AggregateType,
		AggregateID:       entity.AggregateID,
		Payload:           payload,
		DispatchAttempts:  entity.DispatchAttempts,
		AvailableAt:       entity.AvailableAt,
		DispatchedAt:      entity.DispatchedAt,
		ProcessedAt:       entity.ProcessedAt,
		LastDispatchError: entity.LastDispatchError,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}, nil
}

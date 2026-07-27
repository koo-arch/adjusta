package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/mixins"
)

type OutboxMessage struct {
	ent.Schema
}

func (OutboxMessage) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("message_type").NotEmpty(),
		field.String("aggregate_type").NotEmpty(),
		field.UUID("aggregate_id", uuid.UUID{}),
		field.JSON("payload", map[string]any{}).Default(map[string]any{}),
		field.Int("dispatch_attempts").Default(0).NonNegative(),
		field.Time("available_at"),
		field.Time("dispatched_at").Optional().Nillable(),
		field.Time("processed_at").Optional().Nillable(),
		field.Text("last_dispatch_error").Optional().Nillable(),
	}
}

func (OutboxMessage) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("available_at").
			Annotations(entsql.IndexWhere("dispatched_at IS NULL AND processed_at IS NULL")),
		index.Fields("processed_at"),
		index.Fields("aggregate_type", "aggregate_id", "created_at"),
	}
}

func (OutboxMessage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.SchemaMixin{},
		mixins.TimeMixin{},
	}
}

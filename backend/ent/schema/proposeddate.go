package schema

import (
	"context"
	"errors"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	gen "github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/hook"
	"github.com/koo-arch/adjusta-backend/ent/mixins"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

// ProposedDate holds the schema definition for the ProposedDate entity.
type ProposedDate struct {
	ent.Schema
}

// Fields of the ProposedDate.
func (ProposedDate) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.UUID("event_id", uuid.UUID{}),
		field.String("google_event_id").Optional().Nillable(),
		field.Time("start_time"),
		field.Time("end_time"),
		field.Int("priority").Default(0),
		field.Enum("status").
			Values(
				string(domainvalue.ProposedDateStatusActive),
				string(domainvalue.ProposedDateStatusConfirmed),
				string(domainvalue.ProposedDateStatusNotSelected),
				string(domainvalue.ProposedDateStatusCancelled),
			).
			Default(string(domainvalue.ProposedDateStatusActive)),
		field.Enum("sync_status").
			Values(
				string(domainvalue.SyncStatusNotSynced),
				string(domainvalue.SyncStatusPending),
				string(domainvalue.SyncStatusSynced),
				string(domainvalue.SyncStatusFailed),
			).
			Default(string(domainvalue.SyncStatusNotSynced)),
		field.Time("last_synced_at").Optional().Nillable(),
		field.Text("last_sync_error").Optional().Nillable(),
	}
}

// Edges of the ProposedDate.
func (ProposedDate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).Ref("proposed_dates").Field("event_id").Unique().Required(),
	}
}

func (ProposedDate) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("event_id"),
		index.Fields("start_time"),
		index.Fields("status"),
		index.Fields("sync_status"),
		index.Fields("event_id", "priority").
			Unique().
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
	}
}

func (ProposedDate) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(proposeddateHook, ent.OpCreate|ent.OpUpdate),
	}
}

func proposeddateHook(next ent.Mutator) ent.Mutator {
	return hook.ProposedDateFunc(func(ctx context.Context, m *gen.ProposedDateMutation) (ent.Value, error) {
		if startTime, ok := m.StartTime(); ok {
			if endTime, ok := m.EndTime(); ok {
				if startTime.After(endTime) {
					return nil, errors.New("start_time must be before end_time")
				}
			}
		}
		return next.Mutate(ctx, m)
	})
}

func (ProposedDate) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

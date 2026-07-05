package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/mixins"
)

const (
	StatusDraft     = value.StatusDraft
	StatusActive    = value.StatusActive
	StatusConfirmed = value.StatusConfirmed
	StatusCancelled = value.StatusCancelled
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

// Fields of the Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("primary_calendar_id", uuid.UUID{}),
		field.String("title").NotEmpty(),
		field.Text("description").Optional(),
		field.String("location").Optional(),
		field.Enum("status").
			Values(
				string(StatusDraft),
				string(StatusActive),
				string(StatusConfirmed),
				string(StatusCancelled),
			).
			Default(string(StatusActive)),
		field.UUID("confirmed_date_id", uuid.UUID{}).Optional(),
		field.String("confirmed_google_event_id").Optional().Nillable(),
		field.Enum("sync_status").
			Values(
				string(value.SyncStatusNotSynced),
				string(value.SyncStatusPending),
				string(value.SyncStatusSynced),
				string(value.SyncStatusFailed),
			).
			Default(string(value.SyncStatusNotSynced)),
		field.Time("last_synced_at").Optional().Nillable(),
		field.Text("last_sync_error").Optional().Nillable(),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("events").Field("user_id").Unique().Required(),
		edge.From("primary_calendar", Calendar.Type).Ref("primary_events").Field("primary_calendar_id").Unique().Required(),
		edge.To("confirmed_date", ProposedDate.Type).Field("confirmed_date_id").Unique(),
		edge.To("proposed_dates", ProposedDate.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Event) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("primary_calendar_id"),
		index.Fields("confirmed_date_id"),
		index.Fields("status"),
		index.Fields("sync_status"),
	}
}

func (Event) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/mixins"
)

// Calendar holds the schema definition for the Calendar entity.
type Calendar struct {
	ent.Schema
}

// Fields of the Calendar.
func (Calendar) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("google_calendar_id").Optional().Nillable().Unique(),
		field.String("summary").Optional().Nillable(),
		field.Text("description").Optional().Nillable(),
		field.String("timezone").Optional().Nillable(),
	}
}

// Edges of the Calendar.
func (Calendar) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_calendars", UserCalendar.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("primary_events", Event.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Calendar) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

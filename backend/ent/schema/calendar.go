package schema

import (
	"entgo.io/ent"
	"github.com/google/uuid"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/dialect/entsql"
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
	}
}

// Edges of the Calendar.
func (Calendar) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("calendars").Unique(),
		edge.To("google_calendar_infos", GoogleCalendarInfo.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("events", Event.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Calendar) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

package schema

import (
	"entgo.io/ent"
	"github.com/google/uuid"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/koo-arch/adjusta-backend/ent/mixins"
)

// GoogleCalendarInfo holds the schema definition for the GoogleCalendarInfo entity.
type GoogleCalendarInfo struct {
	ent.Schema
}

// Fields of the GoogleCalendarInfo.
func (GoogleCalendarInfo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("google_calendar_id").NotEmpty().Unique(),
		field.String("summary").Optional(),
		field.Bool("is_primary").Default(false),
	}
}

// Edges of the GoogleCalendarInfo.
func (GoogleCalendarInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("calendars", Calendar.Type).Ref("google_calendar_infos"),
	}
}

func (GoogleCalendarInfo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/mixins"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

// UserCalendar holds the schema definition for the UserCalendar entity.
type UserCalendar struct {
	ent.Schema
}

// Fields of the UserCalendar.
func (UserCalendar) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("calendar_id", uuid.UUID{}),
		field.Enum("role").
			Values(
				string(value.UserCalendarRolePrimary),
				string(value.UserCalendarRoleAdjustaCandidate),
				string(value.UserCalendarRoleReference),
			),
		field.Bool("is_visible").Default(true),
		field.Bool("sync_proposed_dates").Default(false),
	}
}

// Edges of the UserCalendar.
func (UserCalendar) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("user_calendars").
			Field("user_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("calendar", Calendar.Type).
			Ref("user_calendars").
			Field("calendar_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (UserCalendar) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "calendar_id").Unique(),
		index.Fields("user_id"),
		index.Fields("calendar_id"),
		index.Fields("role"),
		index.Fields("user_id").
			Unique().
			StorageKey("usercalendar_adjusta_candidate_user_id").
			Annotations(entsql.IndexWhere("role = 'adjusta_candidate' AND deleted_at IS NULL")),
		index.Fields("user_id").
			Unique().
			StorageKey("usercalendar_primary_user_id").
			Annotations(entsql.IndexWhere("role = 'primary' AND deleted_at IS NULL")),
	}
}

func (UserCalendar) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

package schema

import (
	"context"
	"errors"
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	gen "github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/hook"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/mixins"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Email正規表現
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("email").NotEmpty().Unique(),
		field.String("name").Optional().Nillable(),
		field.Text("avatar_url").Optional().Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("account", Account.Type).Unique().Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("sessions", Session.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("user_calendars", UserCalendar.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("events", Event.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(userhook, ent.OpCreate|ent.OpUpdate),
	}
}

func userhook(next ent.Mutator) ent.Mutator {
	return hook.UserFunc(func(ctx context.Context, m *gen.UserMutation) (ent.Value, error) {
		if email, ok := m.Email(); ok {
			if !emailRegex.MatchString(email) {
				return nil, errors.New("invalid email address")
			}
		}
		return next.Mutate(ctx, m)
	})
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

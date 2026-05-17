package schema

import (
	"context"
	"regexp"
	"errors"

	gen "github.com/koo-arch/adjusta-backend/ent"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/dialect/entsql"
	"github.com/koo-arch/adjusta-backend/ent/hook"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/mixins"
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
		field.String("refresh_token").Sensitive().Optional(),
		field.Time("refresh_token_expiry").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("oauth_token", OAuthToken.Type).Unique().Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("calendars", Calendar.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
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
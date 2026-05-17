package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/mixins"
)

// OAuthToken holds the schema definition for the OAuthToken entity.
type OAuthToken struct {
	ent.Schema
}

// Fields of the OAuthToken.
func (OAuthToken) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("access_token").Sensitive().Optional(),
		field.String("refresh_token").Sensitive().Optional(),
		field.Time("expiry").Optional(),
	}
}

// Edges of the OAuthToken.
func (OAuthToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("oauth_token").Unique(),
	}
}

func (OAuthToken) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

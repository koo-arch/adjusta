package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/koo-arch/adjusta-backend/ent/mixins"
)

// JWTkey holds the schema definition for the JWTkey entity.
type JWTKey struct {
	ent.Schema
}

// Fields of the JWTkey.
func (JWTKey) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").NotEmpty().Sensitive(),
		field.String("type").NotEmpty().Default("access"), // access or refresh
		field.Time("expires_at").Immutable(),
	}
}

// Edges of the JWTkey.
func (JWTKey) Edges() []ent.Edge {
	return nil
}

func (JWTKey) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key", "expires_at").StorageKey("idx_type_expires"),
	}
}

func (JWTKey) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}
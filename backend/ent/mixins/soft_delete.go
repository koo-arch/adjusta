package mixins

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/koo-arch/adjusta-backend/ent/intercept"
)

// SoftDeleteMixin is a mixin for adding deleted_at field to an entity.
type SoftDeleteMixin struct {
	mixin.Schema
}

// Fields of the SoftDeleteMixin.
func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

type softDeleteKey struct{}

// SoftDelete is a mixin for soft deleting an entity.
func SkipSoftDelete(ctx context.Context) context.Context {
	return context.WithValue(ctx, softDeleteKey{}, true)
}

func (SoftDeleteMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}
			
			q.WhereP(func(s *sql.Selector) {
				s.Where(sql.IsNull(s.C("deleted_at")))
			})
			return nil
		}),
	}
}
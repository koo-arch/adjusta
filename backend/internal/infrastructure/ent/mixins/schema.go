package mixins

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/mixin"
)

const DatabaseSchema = "adjusta"

// SchemaMixinは、entのテーブルをAdjusta専用schemaへ配置する。
type SchemaMixin struct {
	mixin.Schema
}

func (SchemaMixin) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Schema(DatabaseSchema),
	}
}

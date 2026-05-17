package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/dialect/entsql"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/hook"
	gen "github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/utils"
	"github.com/koo-arch/adjusta-backend/ent/mixins"
)

const (
	StatusPending = models.StatusPending
	StatusConfirmed = models.StatusConfirmed
	StatusCancelled = models.StatusCancelled

	randomStringLen = 4
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

// Fields of the Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("summary").Optional(),
		field.String("description").Optional(),
		field.String("location").Optional(),
		field.Enum("status").
			Values(string(StatusPending), string(StatusConfirmed), string(StatusCancelled)).
			Default(string(StatusPending)),
		field.UUID("confirmed_date_id", uuid.UUID{}).Optional(),
		field.String("google_event_id").Optional(),
		field.String("slug").Unique(),

	}
}

func (Event) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(generateSlug, ent.OpCreate|ent.OpUpdate),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("calendar", Calendar.Type).Ref("events").Unique(),
		edge.To("proposed_dates", ProposedDate.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Event) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func generateSlug(next ent.Mutator) ent.Mutator {
	return hook.EventFunc(func(ctx context.Context, m *gen.EventMutation) (ent.Value, error) {
		summary, ok := m.Summary()
		if !ok {
			println("summary is not set")
			return next.Mutate(ctx, m)
		}

		// updateの場合に古いsummaryと比較
		if m.Op().Is(ent.OpUpdate) {
			oldSummary, err := m.OldSummary(ctx)
			if err == nil && summary == oldSummary {
				return next.Mutate(ctx, m)
			}
		}

		baseSlug, err := utils.NormalizeToSlug(ctx, summary)
		if err != nil {
			return nil, err
		}
		
		// 既存のスラッグを取得
		existingSlugs, err := m.Client().Event.
			Query().
			Where(event.SlugContains(baseSlug)).
			Select(event.FieldSlug).
			All(ctx)
		if err != nil {
			return nil, err
		}

		uniqueSlug := baseSlug
		existingSlugMap := make(map[string]struct{})
		for _, e := range existingSlugs {
			existingSlugMap[e.Slug] = struct{}{}
		}

		// 既存のスラッグと競合しないスラッグを生成
		uniqueSlug = utils.EnsureUniqueSlug(ctx, existingSlugMap, baseSlug, randomStringLen)

		m.SetSlug(uniqueSlug)
		return next.Mutate(ctx, m)
	})
}
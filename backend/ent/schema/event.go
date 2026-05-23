package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	gen "github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/hook"
	"github.com/koo-arch/adjusta-backend/ent/mixins"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/utils"
)

const (
	StatusDraft     = models.StatusDraft
	StatusActive    = models.StatusActive
	StatusPending   = models.StatusPending
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
		field.UUID("user_id", uuid.UUID{}).Optional().Nillable(),
		field.UUID("primary_calendar_id", uuid.UUID{}).Optional().Nillable(),
		field.String("summary").Optional(),
		field.String("title").Optional().Nillable(),
		field.String("description").Optional(),
		field.String("location").Optional(),
		field.Enum("status").
			Values(
				string(StatusDraft),
				string(StatusActive),
				string(StatusPending),
				string(StatusConfirmed),
				string(StatusCancelled),
			).
			Default(string(StatusPending)),
		field.UUID("confirmed_date_id", uuid.UUID{}).Optional(),
		field.String("google_event_id").Optional(),
		field.String("confirmed_google_event_id").Optional().Nillable(),
		field.Enum("sync_status").
			Values(
				string(models.SyncStatusNotSynced),
				string(models.SyncStatusPending),
				string(models.SyncStatusSynced),
				string(models.SyncStatusFailed),
			).
			Default(string(models.SyncStatusNotSynced)),
		field.Time("last_synced_at").Optional().Nillable(),
		field.Text("last_sync_error").Optional().Nillable(),
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
		edge.From("user", User.Type).Ref("events").Field("user_id").Unique(),
		edge.From("primary_calendar", Calendar.Type).Ref("primary_events").Field("primary_calendar_id").Unique(),
		edge.To("confirmed_date", ProposedDate.Type).Field("confirmed_date_id").Unique(),
		edge.To("proposed_dates", ProposedDate.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Event) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("primary_calendar_id"),
		index.Fields("confirmed_date_id"),
		index.Fields("status"),
		index.Fields("sync_status"),
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

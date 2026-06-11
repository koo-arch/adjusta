package event

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

func TestBuildProposedDateChangeSet(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	newStart := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	newEnd := newStart.Add(time.Hour)

	changeSet := BuildProposedDateChangeSet(
		[]DraftProposedDate{
			{
				ID:       &id1,
				Start:    newStart,
				End:      newEnd,
				Priority: 2,
			},
			{
				Start:    newStart.Add(2 * time.Hour),
				End:      newEnd.Add(2 * time.Hour),
				Priority: 3,
			},
		},
		[]ExistingProposedDate{
			{
				ID:       id1,
				Start:    newStart.Add(-24 * time.Hour),
				End:      newEnd.Add(-24 * time.Hour),
				Priority: 1,
			},
			{
				ID:       id2,
				Start:    newStart.Add(-48 * time.Hour),
				End:      newEnd.Add(-48 * time.Hour),
				Priority: 2,
			},
		},
	)

	if len(changeSet.Updates) != 1 {
		t.Fatalf("expected 1 update, got %d", len(changeSet.Updates))
	}
	if changeSet.Updates[0].ID != id1 {
		t.Fatalf("expected update for %s, got %s", id1, changeSet.Updates[0].ID)
	}
	if len(changeSet.Creates) != 1 {
		t.Fatalf("expected 1 create, got %d", len(changeSet.Creates))
	}
	if len(changeSet.Deletes) != 1 {
		t.Fatalf("expected 1 delete, got %d", len(changeSet.Deletes))
	}
	if changeSet.Deletes[0] != id2 {
		t.Fatalf("expected delete for %s, got %s", id2, changeSet.Deletes[0])
	}
}

func TestBuildProposedDateChangeSetReusesGapsForMovedItems(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	changeSet := BuildProposedDateChangeSet(
		[]DraftProposedDate{
			{ID: &id1, Start: start, End: end},
			{ID: &id3, Start: start.Add(4 * time.Hour), End: end.Add(4 * time.Hour)},
			{ID: &id2, Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour)},
		},
		[]ExistingProposedDate{
			{ID: id1, Start: start, End: end, Priority: 3072},
			{ID: id2, Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour), Priority: 2048},
			{ID: id3, Start: start.Add(4 * time.Hour), End: end.Add(4 * time.Hour), Priority: 1024},
		},
	)

	if len(changeSet.Creates) != 0 {
		t.Fatalf("expected no create, got %+v", changeSet.Creates)
	}
	if len(changeSet.Deletes) != 0 {
		t.Fatalf("expected no delete, got %+v", changeSet.Deletes)
	}
	if len(changeSet.Updates) != 1 {
		t.Fatalf("expected 1 update, got %+v", changeSet.Updates)
	}
	if changeSet.Updates[0].ID != id3 {
		t.Fatalf("expected moved date %s to update, got %+v", id3, changeSet.Updates)
	}
	if changeSet.Updates[0].Priority != 2560 {
		t.Fatalf("expected moved date priority 2560, got %+v", changeSet.Updates[0])
	}
}

func TestBuildProposedDateChangeSetFallsBackToResequenceWhenGapIsInsufficient(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	changeSet := BuildProposedDateChangeSet(
		[]DraftProposedDate{
			{ID: &id1, Start: start, End: end},
			{ID: &id3, Start: start.Add(4 * time.Hour), End: end.Add(4 * time.Hour)},
			{ID: &id2, Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour)},
		},
		[]ExistingProposedDate{
			{ID: id1, Start: start, End: end, Priority: 3},
			{ID: id2, Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour), Priority: 2},
			{ID: id3, Start: start.Add(4 * time.Hour), End: end.Add(4 * time.Hour), Priority: 1},
		},
	)

	if len(changeSet.Updates) != 3 {
		t.Fatalf("expected full resequence update, got %+v", changeSet.Updates)
	}

	expected := map[uuid.UUID]int{
		id1: 3072,
		id3: 2048,
		id2: 1024,
	}
	for _, update := range changeSet.Updates {
		if expected[update.ID] != update.Priority {
			t.Fatalf("unexpected update priority for %s: %+v", update.ID, update)
		}
	}
}

func TestBuildConfirmationPlanForExistingDate(t *testing.T) {
	id := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	existing := []ExistingProposedDate{
		{ID: id, Start: start, End: end, Priority: 1024},
		{ID: uuid.New(), Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour), Priority: 2048},
	}

	plan, err := BuildConfirmationPlan(DraftProposedDate{
		ID:       &id,
		Start:    start,
		End:      end,
		Priority: 3,
	}, existing)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if plan.Status != domainvalue.StatusConfirmed {
		t.Fatalf("expected confirmed status, got %s", plan.Status)
	}
	if plan.Update == nil || plan.Update.ID != id || plan.Update.Priority != 3072 {
		t.Fatalf("expected existing date update to highest priority, got %+v", plan.Update)
	}
	if plan.Create != nil {
		t.Fatalf("expected no create plan, got %+v", plan.Create)
	}
}

func TestBuildConfirmationPlanForNewDate(t *testing.T) {
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	existing := []ExistingProposedDate{
		{ID: uuid.New(), Start: start.Add(-2 * time.Hour), End: end.Add(-2 * time.Hour), Priority: 2048},
		{ID: uuid.New(), Start: start.Add(-4 * time.Hour), End: end.Add(-4 * time.Hour), Priority: 1024},
	}

	plan, err := BuildConfirmationPlan(DraftProposedDate{
		Start:    start,
		End:      end,
		Priority: 2,
	}, existing)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if plan.Create == nil || plan.Create.Priority != 3072 {
		t.Fatalf("expected create plan with highest priority, got %+v", plan.Create)
	}
	if plan.Update != nil {
		t.Fatalf("expected no update plan, got %+v", plan.Update)
	}
}

func TestNormalizeDraftProposedDatesByOrderAssignsSparseDescendingPriorities(t *testing.T) {
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	normalized := NormalizeDraftProposedDatesByOrder([]DraftProposedDate{
		{Start: start, End: end},
		{Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour)},
		{Start: start.Add(4 * time.Hour), End: end.Add(4 * time.Hour)},
	})

	if len(normalized) != 3 {
		t.Fatalf("expected 3 dates, got %d", len(normalized))
	}
	if normalized[0].Priority != 3072 {
		t.Fatalf("expected first date priority 3072, got %d", normalized[0].Priority)
	}
	if normalized[1].Priority != 2048 {
		t.Fatalf("expected second date priority 2048, got %d", normalized[1].Priority)
	}
	if normalized[2].Priority != 1024 {
		t.Fatalf("expected third date priority 1024, got %d", normalized[2].Priority)
	}
}

func TestBuildConfirmationPlanRejectsInvalidRange(t *testing.T) {
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)

	_, err := BuildConfirmationPlan(DraftProposedDate{
		Start: start,
		End:   start,
	}, nil)
	if err == nil {
		t.Fatal("expected invalid range error")
	}
}

package event

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPlanProposedDateChanges(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	newStart := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	newEnd := newStart.Add(time.Hour)

	changeSet := PlanProposedDateChanges(
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

func TestPlanProposedDateChangesReusesGapsForMovedItems(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	changeSet := PlanProposedDateChanges(
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

func TestPlanProposedDateChangesFallsBackToResequenceWhenGapIsInsufficient(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	changeSet := PlanProposedDateChanges(
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

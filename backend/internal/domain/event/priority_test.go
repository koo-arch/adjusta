package event

import (
	"testing"
	"time"
)

func TestAssignPrioritiesByOrderAssignsSparseDescendingPriorities(t *testing.T) {
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	assigned := AssignPrioritiesByOrder([]DraftProposedDate{
		{Start: start, End: end},
		{Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour)},
		{Start: start.Add(4 * time.Hour), End: end.Add(4 * time.Hour)},
	})

	if len(assigned) != 3 {
		t.Fatalf("expected 3 dates, got %d", len(assigned))
	}
	if assigned[0].Priority != 3072 {
		t.Fatalf("expected first date priority 3072, got %d", assigned[0].Priority)
	}
	if assigned[1].Priority != 2048 {
		t.Fatalf("expected second date priority 2048, got %d", assigned[1].Priority)
	}
	if assigned[2].Priority != 1024 {
		t.Fatalf("expected third date priority 1024, got %d", assigned[2].Priority)
	}
}

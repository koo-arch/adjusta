package events

import (
	"testing"
	"time"
)

func TestAssignSelectedDatePrioritiesAssignsSparseDescendingValues(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)

	assigned := assignSelectedDatePriorities([]SelectedDate{
		{Start: start, End: start.Add(time.Hour), Priority: 1},
		{Start: start.Add(2 * time.Hour), End: start.Add(3 * time.Hour), Priority: 2},
		{Start: start.Add(4 * time.Hour), End: start.Add(5 * time.Hour), Priority: 3},
	})

	if len(assigned) != 3 {
		t.Fatalf("expected 3 selected dates, got %d", len(assigned))
	}
	if assigned[0].Priority != 3072 || assigned[1].Priority != 2048 || assigned[2].Priority != 1024 {
		t.Fatalf("unexpected assigned priorities: %#v", assigned)
	}
}

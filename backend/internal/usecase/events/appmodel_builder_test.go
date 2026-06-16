package events

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestBuildAppProposedDatesSortsByDescendingPriority(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)

	proposedDates := buildAppProposedDates([]*ProposedDateRecord{
		{ID: uuid.New(), StartTime: start, EndTime: start.Add(time.Hour), Priority: 1024},
		{ID: uuid.New(), StartTime: start.Add(2 * time.Hour), EndTime: start.Add(3 * time.Hour), Priority: 3072},
		{ID: uuid.New(), StartTime: start.Add(4 * time.Hour), EndTime: start.Add(5 * time.Hour), Priority: 2048},
	})

	if len(proposedDates) != 3 {
		t.Fatalf("expected 3 proposed dates, got %d", len(proposedDates))
	}
	if proposedDates[0].Priority != 3072 || proposedDates[1].Priority != 2048 || proposedDates[2].Priority != 1024 {
		t.Fatalf("unexpected priority order: %#v", proposedDates)
	}
}

package events

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func TestBuildProposedDateOutputsSortsByDescendingPriority(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)

	proposedDates := buildProposedDateOutputs([]*ProposedDateRecord{
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

func TestBuildUpcomingEventOutputUsesConfirmedDate(t *testing.T) {
	t.Parallel()

	eventID := uuid.New()
	confirmedDateID := uuid.New()
	otherDateID := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)

	output, err := buildUpcomingEventOutput(&EventRecord{
		ID:              eventID,
		Title:           "Upcoming",
		Status:          value.StatusConfirmed,
		ConfirmedDateID: confirmedDateID,
		ProposedDates: []*ProposedDateRecord{
			{ID: otherDateID, StartTime: start.Add(2 * time.Hour), EndTime: start.Add(3 * time.Hour)},
			{ID: confirmedDateID, StartTime: start, EndTime: start.Add(time.Hour)},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output == nil {
		t.Fatal("expected upcoming output")
	}
	if output.ID != eventID {
		t.Fatalf("unexpected event id: %s", output.ID)
	}
	if output.Start != start {
		t.Fatalf("unexpected start: %s", output.Start)
	}
}

func TestBuildNeedsActionDraftOutputMarksOverdueDraft(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	start := now.Add(-time.Hour)

	output, err := buildNeedsActionDraftOutput(&EventRecord{
		ID:     uuid.New(),
		Title:  "Needs action",
		Status: value.StatusActive,
		ProposedDates: []*ProposedDateRecord{
			{ID: uuid.New(), StartTime: start, EndTime: start.Add(time.Hour)},
		},
	}, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output == nil {
		t.Fatal("expected needs action output")
	}
	if !output.NeedsAttention {
		t.Fatal("expected needs attention")
	}
	if output.Start != start {
		t.Fatalf("unexpected start: %s", output.Start)
	}
}

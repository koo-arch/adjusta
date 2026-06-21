package event

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func TestPlanConfirmationChangesForExistingDate(t *testing.T) {
	id := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	existing := []ExistingProposedDate{
		{ID: id, Start: start, End: end, Priority: 1024},
		{ID: uuid.New(), Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour), Priority: 2048},
	}

	changeSet, err := PlanConfirmationChanges(DraftProposedDate{
		ID:       &id,
		Start:    start,
		End:      end,
		Priority: 3,
	}, existing)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if changeSet.Status != value.StatusConfirmed {
		t.Fatalf("expected confirmed status, got %s", changeSet.Status)
	}
	if changeSet.Update == nil || changeSet.Update.ID != id || changeSet.Update.Priority != 3072 {
		t.Fatalf("expected existing date update to highest priority, got %+v", changeSet.Update)
	}
	if changeSet.Create != nil {
		t.Fatalf("expected no create plan, got %+v", changeSet.Create)
	}
}

func TestPlanConfirmationChangesForNewDate(t *testing.T) {
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	existing := []ExistingProposedDate{
		{ID: uuid.New(), Start: start.Add(-2 * time.Hour), End: end.Add(-2 * time.Hour), Priority: 2048},
		{ID: uuid.New(), Start: start.Add(-4 * time.Hour), End: end.Add(-4 * time.Hour), Priority: 1024},
	}

	changeSet, err := PlanConfirmationChanges(DraftProposedDate{
		Start:    start,
		End:      end,
		Priority: 2,
	}, existing)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if changeSet.Create == nil || changeSet.Create.Priority != 3072 {
		t.Fatalf("expected create plan with highest priority, got %+v", changeSet.Create)
	}
	if changeSet.Update != nil {
		t.Fatalf("expected no update plan, got %+v", changeSet.Update)
	}
}

func TestPlanConfirmationChangesRejectsInvalidRange(t *testing.T) {
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)

	_, err := PlanConfirmationChanges(DraftProposedDate{
		Start: start,
		End:   start,
	}, nil)
	if err == nil {
		t.Fatal("expected invalid range error")
	}
}

func TestPlanConfirmationChangesForExistingDateMarksOthersNotSelected(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	changeSet, err := PlanConfirmationChanges(
		DraftProposedDate{
			ID:       &id2,
			Start:    start.Add(2 * time.Hour),
			End:      end.Add(2 * time.Hour),
			Priority: 2,
		},
		[]ExistingProposedDate{
			{ID: id1, Start: start, End: end, Priority: 3072},
			{ID: id2, Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour), Priority: 2048},
			{ID: id3, Start: start.Add(4 * time.Hour), End: end.Add(4 * time.Hour), Priority: 1024},
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if changeSet.Update == nil || changeSet.Update.ID != id2 {
		t.Fatalf("expected update for confirmed date, got %+v", changeSet.Update)
	}
	if len(changeSet.MarkNotSelected) != 2 {
		t.Fatalf("expected 2 not selected dates, got %+v", changeSet.MarkNotSelected)
	}
	if changeSet.MarkNotSelected[0] != id1 || changeSet.MarkNotSelected[1] != id3 {
		t.Fatalf("unexpected not selected ids: %+v", changeSet.MarkNotSelected)
	}
}

func TestPlanConfirmationChangesForNewDateMarksAllExistingNotSelected(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	changeSet, err := PlanConfirmationChanges(
		DraftProposedDate{
			Start:    start.Add(6 * time.Hour),
			End:      end.Add(6 * time.Hour),
			Priority: 1,
		},
		[]ExistingProposedDate{
			{ID: id1, Start: start, End: end, Priority: 2048},
			{ID: id2, Start: start.Add(2 * time.Hour), End: end.Add(2 * time.Hour), Priority: 1024},
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if changeSet.Create == nil {
		t.Fatalf("expected create plan, got %+v", changeSet)
	}
	if len(changeSet.MarkNotSelected) != 2 {
		t.Fatalf("expected all existing dates to become not selected, got %+v", changeSet.MarkNotSelected)
	}
}

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

func TestBuildConfirmationPlanForExistingDate(t *testing.T) {
	id := uuid.New()
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	plan, err := BuildConfirmationPlan(DraftProposedDate{
		ID:       &id,
		Start:    start,
		End:      end,
		Priority: 3,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if plan.Status != domainvalue.StatusConfirmed {
		t.Fatalf("expected confirmed status, got %s", plan.Status)
	}
	if plan.Update == nil || plan.Update.ID != id || plan.Update.Priority != 0 {
		t.Fatalf("expected existing date update to priority 0, got %+v", plan.Update)
	}
	if plan.Create != nil {
		t.Fatalf("expected no create plan, got %+v", plan.Create)
	}
	if plan.PriorityAdjustment != PriorityAdjustmentReorderRemaining {
		t.Fatalf("expected reorder adjustment, got %s", plan.PriorityAdjustment)
	}
}

func TestBuildConfirmationPlanForNewDate(t *testing.T) {
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	plan, err := BuildConfirmationPlan(DraftProposedDate{
		Start:    start,
		End:      end,
		Priority: 2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if plan.Create == nil || plan.Create.Priority != 1 {
		t.Fatalf("expected create plan with priority 1, got %+v", plan.Create)
	}
	if plan.Update != nil {
		t.Fatalf("expected no update plan, got %+v", plan.Update)
	}
	if plan.PriorityAdjustment != PriorityAdjustmentIncrementOthers {
		t.Fatalf("expected increment adjustment, got %s", plan.PriorityAdjustment)
	}
}

func TestBuildConfirmationPlanRejectsInvalidRange(t *testing.T) {
	start := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)

	_, err := BuildConfirmationPlan(DraftProposedDate{
		Start: start,
		End:   start,
	})
	if err == nil {
		t.Fatal("expected invalid range error")
	}
}

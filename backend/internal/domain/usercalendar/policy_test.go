package usercalendar

import (
	"testing"

	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func TestExternalSyncRole(t *testing.T) {
	t.Parallel()

	if role := ExternalSyncRole(true); role != value.UserCalendarRolePrimary {
		t.Fatalf("expected primary role, got %s", role)
	}
	if role := ExternalSyncRole(false); role != value.UserCalendarRoleReference {
		t.Fatalf("expected reference role, got %s", role)
	}
}

func TestIsExternalSyncRole(t *testing.T) {
	t.Parallel()

	if !IsExternalSyncRole(value.UserCalendarRolePrimary) {
		t.Fatal("expected primary to be externally synced")
	}
	if !IsExternalSyncRole(value.UserCalendarRoleReference) {
		t.Fatal("expected reference to be externally synced")
	}
	if IsExternalSyncRole(value.UserCalendarRoleAdjustaCandidate) {
		t.Fatal("expected adjusta candidate to be excluded from external sync")
	}
}

func TestIsAdjustaCandidateCalendarSummary(t *testing.T) {
	t.Parallel()

	if !IsAdjustaCandidateCalendarSummary(AdjustaCandidateCalendarSummary) {
		t.Fatal("expected adjusta candidate summary to be recognized")
	}
	if IsAdjustaCandidateCalendarSummary("Primary") {
		t.Fatal("did not expect unrelated summary to be recognized")
	}
}

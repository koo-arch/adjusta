package events

import (
	"testing"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func TestDraftEventUpdateSetsSyncStatus(t *testing.T) {
	status := value.StatusActive

	pending := draftEventUpdate(&status, true)
	if pending.SyncStatus == nil || *pending.SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected pending sync status: %#v", pending.SyncStatus)
	}

	notSynced := draftEventUpdate(&status, false)
	if notSynced.SyncStatus == nil || *notSynced.SyncStatus != value.SyncStatusNotSynced {
		t.Fatalf("unexpected not synced status: %#v", notSynced.SyncStatus)
	}
	if !notSynced.ClearLastSyncError {
		t.Fatal("expected last sync error to be cleared")
	}
}

func TestConfirmedEventPendingUpdate(t *testing.T) {
	confirmedDateID := uuid.New()

	update := confirmedEventPendingUpdate(confirmedDateID)

	if update.Status == nil || *update.Status != value.StatusConfirmed {
		t.Fatalf("unexpected event status: %#v", update.Status)
	}
	if update.ConfirmedDateID == nil || *update.ConfirmedDateID != confirmedDateID {
		t.Fatalf("unexpected confirmed date id: %#v", update.ConfirmedDateID)
	}
	if update.SyncStatus == nil || *update.SyncStatus != value.SyncStatusPending {
		t.Fatalf("unexpected sync status: %#v", update.SyncStatus)
	}
}

func TestGoogleEventIDOrEmpty(t *testing.T) {
	id := "google-event-id"
	if got := googleEventIDOrEmpty(&id); got != id {
		t.Fatalf("unexpected google event id: %q", got)
	}
	if got := googleEventIDOrEmpty(nil); got != "" {
		t.Fatalf("unexpected fallback google event id: %q", got)
	}
}

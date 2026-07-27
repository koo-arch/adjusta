package googlecalendar

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func TestConfirmedEventSyncedUpdate(t *testing.T) {
	confirmedDateID := uuid.New()
	syncedAt := time.Date(2026, 7, 27, 10, 0, 0, 0, time.UTC)

	update := confirmedEventSyncedUpdate(confirmedDateID, "google-event-id", syncedAt)

	if update.Status == nil || *update.Status != value.StatusConfirmed {
		t.Fatalf("unexpected event status: %#v", update.Status)
	}
	if update.ConfirmedDateID == nil || *update.ConfirmedDateID != confirmedDateID {
		t.Fatalf("unexpected confirmed date id: %#v", update.ConfirmedDateID)
	}
	if update.ConfirmedGoogleEventID == nil || *update.ConfirmedGoogleEventID != "google-event-id" {
		t.Fatalf("unexpected google event id: %#v", update.ConfirmedGoogleEventID)
	}
	if update.SyncStatus == nil || *update.SyncStatus != value.SyncStatusSynced {
		t.Fatalf("unexpected sync status: %#v", update.SyncStatus)
	}
	if update.LastSyncedAt == nil || !update.LastSyncedAt.Equal(syncedAt) {
		t.Fatalf("unexpected synced at: %#v", update.LastSyncedAt)
	}
	if !update.ClearLastSyncError {
		t.Fatal("expected last sync error to be cleared")
	}
}

func TestEventFailedUpdate(t *testing.T) {
	update := eventFailedUpdate(errors.New("google unavailable"))

	if update.SyncStatus == nil || *update.SyncStatus != value.SyncStatusFailed {
		t.Fatalf("unexpected sync status: %#v", update.SyncStatus)
	}
	if update.LastSyncError == nil || *update.LastSyncError != "google unavailable" {
		t.Fatalf("unexpected sync error: %#v", update.LastSyncError)
	}
}

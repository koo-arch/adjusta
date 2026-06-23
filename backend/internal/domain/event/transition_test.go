package event

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

func TestNewPendingEventChange(t *testing.T) {
	status := value.StatusActive

	change := NewPendingEventChange(&status)

	if change.Status == nil || *change.Status != value.StatusActive {
		t.Fatalf("unexpected status: %#v", change.Status)
	}
	if change.Sync.Status != value.SyncStatusPending {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
}

func TestNewNotSyncedEventChange(t *testing.T) {
	status := value.StatusActive

	change := NewNotSyncedEventChange(&status)

	if change.Status == nil || *change.Status != value.StatusActive {
		t.Fatalf("unexpected status: %#v", change.Status)
	}
	if change.Sync.Status != value.SyncStatusNotSynced {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
	if !change.Sync.ClearLastSyncError {
		t.Fatal("expected clear last sync error flag")
	}
}

func TestNewDraftEventChange(t *testing.T) {
	status := value.StatusActive

	pendingChange := NewDraftEventChange(&status, true)
	if pendingChange.Sync.Status != value.SyncStatusPending {
		t.Fatalf("unexpected pending sync status: %s", pendingChange.Sync.Status)
	}

	notSyncedChange := NewDraftEventChange(&status, false)
	if notSyncedChange.Sync.Status != value.SyncStatusNotSynced {
		t.Fatalf("unexpected not synced status: %s", notSyncedChange.Sync.Status)
	}
	if !notSyncedChange.Sync.ClearLastSyncError {
		t.Fatal("expected clear last sync error flag")
	}
}

func TestNewSyncedEventChange(t *testing.T) {
	confirmedDateID := uuid.New()
	syncedAt := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)

	change := NewSyncedEventChange(value.StatusConfirmed, confirmedDateID, "google-event-id", syncedAt)

	if change.Status == nil || *change.Status != value.StatusConfirmed {
		t.Fatalf("unexpected status: %#v", change.Status)
	}
	if change.ConfirmedDateID == nil || *change.ConfirmedDateID != confirmedDateID {
		t.Fatalf("unexpected confirmed date id: %#v", change.ConfirmedDateID)
	}
	if change.ConfirmedGoogleEventID == nil || *change.ConfirmedGoogleEventID != "google-event-id" {
		t.Fatalf("unexpected confirmed google event id: %#v", change.ConfirmedGoogleEventID)
	}
	if change.Sync.Status != value.SyncStatusSynced {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
	if change.Sync.LastSyncedAt == nil || !change.Sync.LastSyncedAt.Equal(syncedAt) {
		t.Fatalf("unexpected synced at: %#v", change.Sync.LastSyncedAt)
	}
	if !change.Sync.ClearLastSyncError {
		t.Fatal("expected clear last sync error flag")
	}
}

func TestNewFailedEventChange(t *testing.T) {
	change := NewFailedEventChange(errors.New("google unavailable"))

	if change.Sync.Status != value.SyncStatusFailed {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
	if change.Sync.LastSyncError == nil || *change.Sync.LastSyncError != "google unavailable" {
		t.Fatalf("unexpected last sync error: %#v", change.Sync.LastSyncError)
	}
}

func TestNewSyncedEventSyncChange(t *testing.T) {
	syncedAt := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)

	change := NewSyncedEventSyncChange(syncedAt)

	if change.Status != nil {
		t.Fatalf("expected status to remain unchanged, got %#v", change.Status)
	}
	if change.Sync.Status != value.SyncStatusSynced {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
	if change.Sync.LastSyncedAt == nil || !change.Sync.LastSyncedAt.Equal(syncedAt) {
		t.Fatalf("unexpected synced at: %#v", change.Sync.LastSyncedAt)
	}
	if !change.Sync.ClearLastSyncError {
		t.Fatal("expected clear last sync error flag")
	}
}

func TestNewPendingProposedDateChange(t *testing.T) {
	start := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	priority := 1024
	status := value.ProposedDateStatusConfirmed

	change := NewPendingProposedDateChange(&start, &end, &priority, &status)

	if change.Start == nil || !change.Start.Equal(start) {
		t.Fatalf("unexpected start: %#v", change.Start)
	}
	if change.End == nil || !change.End.Equal(end) {
		t.Fatalf("unexpected end: %#v", change.End)
	}
	if change.Priority == nil || *change.Priority != priority {
		t.Fatalf("unexpected priority: %#v", change.Priority)
	}
	if change.Status == nil || *change.Status != status {
		t.Fatalf("unexpected status: %#v", change.Status)
	}
	if change.Sync.Status != value.SyncStatusPending {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
}

func TestNewNotSyncedProposedDateChange(t *testing.T) {
	start := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	priority := 1024
	status := value.ProposedDateStatusConfirmed

	change := NewNotSyncedProposedDateChange(&start, &end, &priority, &status)

	if change.Start == nil || !change.Start.Equal(start) {
		t.Fatalf("unexpected start: %#v", change.Start)
	}
	if change.End == nil || !change.End.Equal(end) {
		t.Fatalf("unexpected end: %#v", change.End)
	}
	if change.Priority == nil || *change.Priority != priority {
		t.Fatalf("unexpected priority: %#v", change.Priority)
	}
	if change.Status == nil || *change.Status != status {
		t.Fatalf("unexpected status: %#v", change.Status)
	}
	if change.Sync.Status != value.SyncStatusNotSynced {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
	if !change.Sync.ClearLastSyncError {
		t.Fatal("expected clear last sync error flag")
	}
}

func TestNewDraftProposedDateChange(t *testing.T) {
	status := value.ProposedDateStatusActive

	pendingChange := NewDraftProposedDateChange(nil, nil, nil, &status, true)
	if pendingChange.Sync.Status != value.SyncStatusPending {
		t.Fatalf("unexpected pending sync status: %s", pendingChange.Sync.Status)
	}

	notSyncedChange := NewDraftProposedDateChange(nil, nil, nil, &status, false)
	if notSyncedChange.Sync.Status != value.SyncStatusNotSynced {
		t.Fatalf("unexpected not synced status: %s", notSyncedChange.Sync.Status)
	}
	if !notSyncedChange.Sync.ClearLastSyncError {
		t.Fatal("expected clear last sync error flag")
	}
}

func TestNewConfirmedProposedDateChange(t *testing.T) {
	start := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	priority := 1024

	change := NewConfirmedProposedDateChange(&start, &end, &priority)

	if change.Status == nil || *change.Status != value.ProposedDateStatusConfirmed {
		t.Fatalf("unexpected status: %#v", change.Status)
	}
	if change.Sync.Status != value.SyncStatusPending {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
}

func TestNewNotSelectedProposedDateChange(t *testing.T) {
	change := NewNotSelectedProposedDateChange()

	if change.Status == nil || *change.Status != value.ProposedDateStatusNotSelected {
		t.Fatalf("unexpected status: %#v", change.Status)
	}
	if change.Sync.Status != value.SyncStatusPending {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
}

func TestNewSyncedProposedDateChange(t *testing.T) {
	syncedAt := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)

	change := NewSyncedProposedDateChange("google-event-id", syncedAt)

	if change.GoogleEventID == nil || *change.GoogleEventID != "google-event-id" {
		t.Fatalf("unexpected google event id: %#v", change.GoogleEventID)
	}
	if change.Sync.Status != value.SyncStatusSynced {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
	if change.Sync.LastSyncedAt == nil || !change.Sync.LastSyncedAt.Equal(syncedAt) {
		t.Fatalf("unexpected synced at: %#v", change.Sync.LastSyncedAt)
	}
	if !change.Sync.ClearLastSyncError {
		t.Fatal("expected clear last sync error flag")
	}
}

func TestNewFailedProposedDateChange(t *testing.T) {
	change := NewFailedProposedDateChange(errors.New("calendar unavailable"))

	if change.Sync.Status != value.SyncStatusFailed {
		t.Fatalf("unexpected sync status: %s", change.Sync.Status)
	}
	if change.Sync.LastSyncError == nil || *change.Sync.LastSyncError != "calendar unavailable" {
		t.Fatalf("unexpected last sync error: %#v", change.Sync.LastSyncError)
	}
}

func TestResolveGoogleEventID(t *testing.T) {
	confirmedGoogleEventID := "confirmed-google-event-id"

	if got := ResolveGoogleEventID(&confirmedGoogleEventID); got != "confirmed-google-event-id" {
		t.Fatalf("unexpected resolved google event id: %q", got)
	}
	if got := ResolveGoogleEventID(nil); got != "" {
		t.Fatalf("unexpected fallback google event id: %q", got)
	}
}

func TestResolveReusableGoogleEventID(t *testing.T) {
	confirmDateID := uuid.New()
	confirmedGoogleEventID := "confirmed-google-event-id"

	if got := ResolveReusableGoogleEventID(nil, &confirmedGoogleEventID, "requested-google-event-id"); got != nil {
		t.Fatalf("expected nil without confirmed date id, got %#v", got)
	}

	got := ResolveReusableGoogleEventID(&confirmDateID, &confirmedGoogleEventID, "requested-google-event-id")
	if got == nil || *got != confirmedGoogleEventID {
		t.Fatalf("unexpected resolved confirmed google event id: %#v", got)
	}

	got = ResolveReusableGoogleEventID(&confirmDateID, nil, "requested-google-event-id")
	if got == nil || *got != "requested-google-event-id" {
		t.Fatalf("unexpected requested google event id: %#v", got)
	}

	got = ResolveReusableGoogleEventID(&confirmDateID, nil, "")
	if got != nil {
		t.Fatalf("expected nil without reusable google event id, got %#v", got)
	}
}

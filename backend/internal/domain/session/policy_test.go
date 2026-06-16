package session

import (
	"testing"
	"time"
)

func TestSessionIsExpiredAt(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC)

	expiredSession := Session{ExpiresAt: now.Add(-time.Minute)}
	if !expiredSession.IsExpiredAt(now) {
		t.Fatal("expected session to be expired")
	}

	currentSession := Session{ExpiresAt: now}
	if currentSession.IsExpiredAt(now) {
		t.Fatal("expected session expiring at now to remain valid")
	}

	futureSession := Session{ExpiresAt: now.Add(time.Minute)}
	if futureSession.IsExpiredAt(now) {
		t.Fatal("expected future session to remain valid")
	}
}

func TestExpiresAtFrom(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC)
	lifetime := 24 * time.Hour

	expiresAt := ExpiresAtFrom(now, lifetime)
	if !expiresAt.Equal(now.Add(lifetime)) {
		t.Fatalf("unexpected expiresAt: %s", expiresAt)
	}
}

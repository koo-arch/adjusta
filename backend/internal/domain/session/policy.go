package session

import "time"

func (s Session) IsExpiredAt(now time.Time) bool {
	return s.ExpiresAt.Before(now)
}

func ExpiresAtFrom(now time.Time, lifetime time.Duration) time.Time {
	return now.Add(lifetime)
}

package cache

import (
	"time"

	"github.com/google/uuid"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
	"github.com/patrickmn/go-cache"
)

type CalendarCache struct {
	c *cache.Cache
}

func NewCalendarCache(ttl, cleanupInterval time.Duration) *CalendarCache {
	return &CalendarCache{
		c: cache.New(ttl, cleanupInterval),
	}
}

func (cc *CalendarCache) Get(userID uuid.UUID) ([]*usecaseCalendar.ExternalCalendar, bool) {
	value, found := cc.c.Get(cc.key(userID))
	if !found {
		return nil, false
	}

	calendars, ok := value.([]*usecaseCalendar.ExternalCalendar)
	if !ok {
		cc.c.Delete(cc.key(userID))
		return nil, false
	}
	return calendars, ok
}

func (cc *CalendarCache) Set(userID uuid.UUID, calendars []*usecaseCalendar.ExternalCalendar) {
	cc.c.Set(cc.key(userID), calendars, cache.DefaultExpiration)
}

func (cc *CalendarCache) Invalidate(userID uuid.UUID) {
	cc.c.Delete(cc.key(userID))
}

func (cc *CalendarCache) key(userID uuid.UUID) string {
	return "calendars:" + userID.String()
}

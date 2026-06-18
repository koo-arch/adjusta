package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	CalendarCache *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		CalendarCache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

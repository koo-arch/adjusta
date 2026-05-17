package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Cache struct {
	JWTKeyCache *cache.Cache
	CalendarCache *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		JWTKeyCache: cache.New(5*time.Minute, 10*time.Minute),
		CalendarCache: cache.New(5*time.Minute, 10*time.Minute),
	}
}
package inmemorycache

import (
	"sync"
	"time"

	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
)

type cacheItem struct {
	snapshot  weather.WeatherSnapshot
	expiresAt time.Time
}

type inMemoryCache struct {
	mu       sync.RWMutex
	items    map[string]cacheItem
	ttl      time.Duration
	stopChan chan struct{}
}

func NewMemoryCache(ttl time.Duration) *inMemoryCache {
	c := &inMemoryCache{
		items:    make(map[string]cacheItem),
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}
	go c.cleanupLoop(time.Minute)
	return c
}

func (c *inMemoryCache) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopChan:
			return
		}
	}
}

func (c *inMemoryCache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for key, entry := range c.items {
		if entry.expiresAt.Before(now) {
			delete(c.items, key)
		}
	}
}

func (c *inMemoryCache) Stop() {
	close(c.stopChan)
}

func (c *inMemoryCache) Set(city string, snapshot weather.WeatherSnapshot) {
	c.mu.Lock()
	defer c.mu.Unlock()
	snapshot.Cached = true
	c.items[city] = cacheItem{
		snapshot:  snapshot,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *inMemoryCache) Get(city string) (weather.WeatherSnapshot, bool) {
	c.mu.RLock()
	entry, ok := c.items[city]
	c.mu.RUnlock()

	if !ok {
		return weather.WeatherSnapshot{}, false
	}

	if entry.expiresAt.Before(time.Now()) {
		return weather.WeatherSnapshot{}, false
	}
	return entry.snapshot, true
}

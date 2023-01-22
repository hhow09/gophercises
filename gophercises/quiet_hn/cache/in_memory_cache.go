package cache

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type CacheKey string

const (
	TOP_STORIES CacheKey = "TOP_STORIES"
)

type CacheItem struct {
	data   interface{}
	expire time.Time
}

type InMemoryCache struct {
	mu      sync.RWMutex
	dataMap map[CacheKey]CacheItem
}

func (c *InMemoryCache) Set(key CacheKey, val interface{}) {
	c.mu.Lock()
	c.dataMap[key] = CacheItem{data: val, expire: time.Now().Add(3 * time.Minute)}
	c.mu.Unlock()
}

func debug(s string) {
	debug := os.Getenv("DEBUG")
	if debug == "1" {
		fmt.Println(s)
	}
}

func (c *InMemoryCache) Get(key CacheKey) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.dataMap[key]
	if !ok {
		return nil
	}
	debug(fmt.Sprintf("expire %v", item.expire))
	if item.expire.Before(time.Now()) {
		delete(c.dataMap, key)
		return nil
	}

	return c.dataMap[key].data
}
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{dataMap: map[CacheKey]CacheItem{}}
}

package db

import (
	"log"
	"sync"
)

type MemoryCache struct {
	sync.RWMutex
	values map[string]interface{}
}

var Cache MemoryCache = MemoryCache{}

func CacheInit() {
	Cache = MemoryCache{values: make(map[string]interface{})}
}

func (c *MemoryCache) Set(key string, value interface{}) {
	c.RLock()
	defer c.RUnlock()
	c.values[key] = value
	log.Printf("Set %s to CACHE", key)
}

func (c *MemoryCache) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()

	item, found := c.values[key]
	if !found {
		return nil
	}
	log.Printf("Get %s from CACHE", key)
	return item
}

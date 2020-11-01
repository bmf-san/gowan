package gowan

import (
	"log"
	"sync"
	"time"
)

// item is the data to be cached.
type item struct {
	value   interface{}
	expires int64
}

// Cache is a struct for caching.
type Cache struct {
	items map[string]*item
	mu    sync.Mutex
}

// New creates a Cache.
func New() *Cache {
	c := &Cache{items: make(map[string]*item)}
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				c.mu.Lock()
				for k, v := range c.items {
					if v.Expired(time.Now().UnixNano()) {
						log.Printf("%v has expires at %d", c.items, time.Now().UnixNano())
						delete(c.items, k)
					}
				}
				c.mu.Unlock()
			}
		}
	}()
	return c
}

// Expired determines if it has expires.
func (i *item) Expired(time int64) bool {
	if i.expires == 0 {
		return true
	}
	return time > i.expires
}

// Get gets a value from a cache.
func (c *Cache) Get(key string) interface{} {
	c.mu.Lock()
	var i interface{}
	if v, ok := c.items[key]; ok {
		i = v.value
	}
	c.mu.Unlock()
	return i
}

// Put puts a value to a cache. If a key and value exists, overwrite it.
func (c *Cache) Put(key string, value interface{}, expires int64) {
	c.mu.Lock()
	c.items[key] = &item{
		value:   value,
		expires: expires,
	}
	c.mu.Unlock()
}

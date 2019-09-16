// Package gocash provides a concurrency-safe in-memory cache. It supports
// key expiration by using Deadlines and Timeouts and can also set a default
// timeout.
package gocash

import (
	"sync"
	"time"
)

// NeverExpires is a placeholder deadline used to designate keys that never
// expire.
var NeverExpires = time.Time{}

// Cache stores arbitrary values indexed by string keys, similar to a Go
// map[string]interface{}. Entries support expiration by Deadline or Timeout.
type Cache struct {
	items sync.Map
	opts  CacheOptions
	mutex sync.Mutex
}

type item struct {
	value    interface{}
	deadline time.Time
}

// CacheOptions holds the options used to initialize the cache.
type CacheOptions struct {
	// The default timeout for keys inserted in the cache without an explicit
	// timeout
	DefaultTimeout time.Duration
}

// NewCache initializaes a new Cache using the given options and returns it.
func NewCache(opts CacheOptions) *Cache {
	return &Cache{
		opts: opts,
	}
}

// Set associates value to key in the cache and stores it for further retrieval.
// If the Cache was initialized with a DefaultTimeout, the key will expire
// after the timeout. Otherwise, the key will never expire, and will have to be
// removed with an explicit Delete.
//
// Returns the actual expiration deadline set in the cache.
func (c *Cache) Set(key string, value interface{}) time.Time {
	if c.opts.DefaultTimeout == 0 {
		return c.SetWithDeadline(key, value, NeverExpires)
	}

	return c.SetWithDeadline(key, value, time.Now().Add(c.opts.DefaultTimeout))
}

// SetWithTimeout associates value to key in the cache and stores it for further
// retrieval. The key will expire after the timeout.
//
// Returns the actual expiration deadline set in the cache.
func (c *Cache) SetWithTimeout(
	key string,
	value interface{},
	timeout time.Duration,
) time.Time {
	return c.SetWithDeadline(key, value, time.Now().Add(timeout))
}

// SetWithDeadline associates value to key in the cache and stores it for
// further retrieval. The key will expire after the timeout.
//
// Returns the actual expiration deadline set in the cache. It should never
// differ from the deadline passed as parameter.
func (c *Cache) SetWithDeadline(
	key string,
	value interface{},
	deadline time.Time,
) time.Time {
	if value == nil {
		panic("cannot set nil value in cache")
	}
	c.items.Store(key, item{
		value:    value,
		deadline: deadline,
	})
	return deadline
}

// Has returns true if the Cache contains an entry for the key and it hasn't
// expired, false otherwise. If an entry was found it also returns the expiration
// deadline of the key.
func (c *Cache) Has(key string) (bool, time.Time) {
	it, t := c.Get(key)
	return (it != nil), t
}

// Get returns the value the Cache contains for the key if it exists and it
// hasn't expired, false otherwise. If an entry was found it also returns the
// expiration deadline of the key.
func (c *Cache) Get(key string) (interface{}, time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	it, ok := c.items.Load(key)
	if !ok {
		return nil, NeverExpires
	}
	it2, ok := it.(item)
	if !ok {
		panic("corrupted cache: unexpected item")
	}
	if it2.deadline != NeverExpires && time.Now().After(it2.deadline) {
		// Item is expired, remove from cache
		c.items.Delete(key)
		return nil, NeverExpires
	}
	return it2.value, it2.deadline
}

// Delete removes an entry from the Cache associated to key. If no such entry
// exists, nothing is done.
func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

// Prune purges from the cache all keys that have expired, releasing resources
// for garbage collection if applicable.
func (c *Cache) Prune() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items.Range(func(k interface{}, it interface{}) bool {
		_, ok := k.(string)
		if !ok {
			panic("corrupted cache: unexpected key")
		}
		it2, ok := it.(item)
		if !ok {
			panic("corrupted cache: unexpected item")
		}
		if it2.deadline != NeverExpires && time.Now().After(it2.deadline) {
			// Item is expired, remove from cache
			c.items.Delete(k)
		}
		return true
	})
}

func (c *Cache) count() int {
	count := 0
	c.items.Range(func(k interface{}, it interface{}) bool {
		count++
		return true
	})
	return count
}

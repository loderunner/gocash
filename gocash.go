package gocash

import (
	"sync"
	"time"
)

// NeverExpires is a placeholder deadline used to designate keys that never
// expire
var NeverExpires = time.Time{}

// Cache is
type Cache struct {
	items sync.Map
	opts  CacheOptions
	mutex sync.Mutex
}

type item struct {
	value    interface{}
	deadline time.Time
}

type CacheOptions struct {
	DefaultTimeout time.Duration
}

func NewCache(opts CacheOptions) *Cache {
	return &Cache{
		opts: opts,
	}
}

func (c *Cache) Set(key string, value interface{}) time.Time {
	if c.opts.DefaultTimeout == 0 {
		return c.SetWithDeadline(key, value, NeverExpires)
	}

	return c.SetWithDeadline(key, value, time.Now().Add(c.opts.DefaultTimeout))
}

func (c *Cache) SetWithTimeout(
	key string,
	value interface{},
	timeout time.Duration,
) time.Time {
	return c.SetWithDeadline(key, value, time.Now().Add(timeout))
}

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

func (c *Cache) Has(key string) (bool, time.Time) {
	it, t := c.Get(key)
	return (it != nil), t
}

func (c *Cache) Get(key string) (interface{}, time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	it, ok := c.items.Load(key)
	if !ok {
		return nil, time.Time{}
	}
	it2, ok := it.(item)
	if !ok {
		panic("corrupted cache: unexpected item")
	}
	if it2.deadline != NeverExpires && time.Now().After(it2.deadline) {
		// Item is expired, remove from cache
		c.items.Delete(key)
		return nil, it2.deadline
	}
	return it2.value, it2.deadline
}

func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

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

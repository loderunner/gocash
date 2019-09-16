package gocash

import (
	"math"
	"sync"
	"time"
)

type Cache struct {
	items sync.Map
	opts  CacheOptions
}

type item struct {
	value    interface{}
	deadline time.Time
}

type CacheOptions struct {
	DefaultTimeout time.Duration
}

func NewCache(opts CacheOptions) *Cache {
	if opts.DefaultTimeout == 0 {
		opts.DefaultTimeout = math.MaxInt64
	}
	return &Cache{
		opts: opts,
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.SetWithDeadline(key, value, time.Now().Add(c.opts.DefaultTimeout))
}

func (c *Cache) SetWithTimeout(key string, value interface{}, timeout time.Duration) {
	c.SetWithDeadline(key, value, time.Now().Add(timeout))
}

func (c *Cache) SetWithDeadline(key string, value interface{}, deadline time.Time) {
	if value == nil {
		panic("cannot set nil value in cache")
	}
	c.items.Store(
		key,
		item{
			value:    value,
			deadline: deadline,
		},
	)
}

func (c *Cache) Has(key string) bool {
	it := c.Get(key)
	return it != nil
}

func (c *Cache) Get(key string) interface{} {
	it, ok := c.items.Load(key)
	if !ok {
		return nil
	}
	it2, ok := it.(item)
	if !ok {
		panic("corrupted cache: unexpected item")
	}
	if time.Now().After(it2.deadline) {
		// Item is expired
		return nil
	}
	return it2.value
}

func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

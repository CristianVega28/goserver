package models

import (
	it "github.com/CristianVega28/goserver/core/interface"
)

type (
	Cache struct {
		store map[string]Store
		it.CacheInterface
	}

	Store struct {
		Value any
		Ttl   int64
	}
)

var Cache_ *Cache

func InitializeCache() {
	Cache_ = &Cache{
		store: make(map[string]Store),
	}
}

func (c *Cache) Get(key string) (*Store, bool) {
	store, exists := c.store[key]
	if !exists {
		return nil, false
	}
	return &store, true
}

func (c *Cache) Set(key string, value any, ttl int64) {
	c.store[key] = Store{
		Value: value,
		Ttl:   ttl,
	}
}

func (c *Cache) Delete(key string) {
	delete(c.store, key)
}

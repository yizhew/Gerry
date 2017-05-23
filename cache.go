package gerry

import (
	gc "github.com/patrickmn/go-cache"
	"time"
)

var cache *gc.Cache

var DefaultCache *Cache

type Cache struct {
}

func init() {
	cache = gc.New(5*time.Minute, 30*time.Second)
	DefaultCache = &Cache{}
}

func (c *Cache) Put(key string, value interface{}) {
	cache.Set(key, value, gc.DefaultExpiration)
}

func (c *Cache) Set(key string, value interface{}, t time.Duration) error {
	cache.Set(key, value, t)

	return nil
}

func (c *Cache) SetForever(key string, value interface{}) {
	cache.Set(key, value, gc.NoExpiration)
}

func (c *Cache) Get(key string) interface{} {
	if value, found := cache.Get(key); found {
		return value
	}

	return nil
}

func (c *Cache) Delete(key string) error {
	cache.Delete(key)
	return nil
}

func (c *Cache) IsExist(key string) bool {
	_, found := cache.Get(key)

	return found
}

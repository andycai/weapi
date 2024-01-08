package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

// https://github.com/allegro/bigcache
// https://github.com/coocood/freecache
// https://github.com/VictoriaMetrics/fastcache

// 序列化
// github.com/tinylib/msgp
// github.com/gogo/protobuf/protoc-gen-gogofaster

type Cache struct {
	cache *bigcache.BigCache
}

func New(eviction time.Duration) *Cache {
	c := new(Cache)
	if bc, err := bigcache.New(context.Background(), bigcache.DefaultConfig(eviction)); err != nil {
		panic("bigcache.New failed: " + err.Error())
	} else {
		c.cache = bc
	}
	return c
}

func (c *Cache) Set(key string, value []byte) error {
	return c.cache.Set(key, value)
}

func (c *Cache) Get(key string) ([]byte, error) {
	return c.cache.Get(key)
}

func (c *Cache) Delete(key string) error {
	return c.cache.Delete(key)
}

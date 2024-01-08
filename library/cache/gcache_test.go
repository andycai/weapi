package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheExpire(t *testing.T) {
	cache := New(1 * time.Second)
	cache.Set("name", []byte("andy"))

	name, err := cache.Get("name")
	assert.Nil(t, err)
	assert.Equal(t, string(name), "andy")

	time.Sleep(2100 * time.Millisecond)
	name, _ = cache.Get("name")
	assert.Nil(t, name)

	cache.Set("name", []byte("andy"))
	cache.Delete("name")
	name, err = cache.Get("name")
	assert.NotNil(t, err)
	assert.Nil(t, name)
}

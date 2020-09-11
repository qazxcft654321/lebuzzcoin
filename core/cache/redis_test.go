package cache

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	cache, err := NewCache(Options{})

	assert.IsType(t, cache, &Client{})
	assert.NoError(t, err)
}

func TestNewTestCache(t *testing.T) {
	cache, err := NewTestCache()

	assert.IsType(t, cache, &Client{})
	assert.NoError(t, err)
}

func TestCacheSet(t *testing.T) {
	cache, err := NewTestCache()
	assert.NoError(t, err)

	err = cache.Set("test", "data", 0)
	assert.NoError(t, err)
}

func TestCacheGet(t *testing.T) {
	cache, err := NewTestCache()
	assert.NoError(t, err)

	val, err := cache.Get("test")
	assert.IsType(t, redis.Nil, err)
	assert.IsType(t, "", val)
}

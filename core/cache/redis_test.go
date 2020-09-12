package cache

import (
	"testing"

	"github.com/go-redis/redis/v8"
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

func TestCacheZadd(t *testing.T) {
	cache, err := NewTestCache()
	assert.NoError(t, err)

	count, err := cache.ZAdd("testSet", &ZMember{Score: 1, Member: "testValue"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

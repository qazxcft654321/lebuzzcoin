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

func TestCacheZAdd(t *testing.T) {
	cache, err := NewTestCache()
	assert.NoError(t, err)

	count, err := cache.ZAdd("testSet", &ZMember{Score: 1, Member: "testValue"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestCacheZIncr(t *testing.T) {
	cache, err := NewTestCache()
	assert.NoError(t, err)

	_, err = cache.ZAdd("testSet", &ZMember{Score: 1, Member: "testValue"})
	assert.NoError(t, err)

	max := 5
	score := float64(0)
	for i := 0; i < max; i++ {
		score, err = cache.ZIncr("testSet", &ZMember{Score: 1, Member: "testValue"})
		assert.NoError(t, err)
	}

	assert.Equal(t, float64(max)+1, score)
}

func ZRevRangeWithScores(t *testing.T) {
	cache, err := NewTestCache()
	assert.NoError(t, err)

	_, err = cache.ZAdd("testSet", &ZMember{Score: 1, Member: "testValue1"})
	assert.NoError(t, err)
	_, err = cache.ZAdd("testSet", &ZMember{Score: 1, Member: "testValue2"})
	assert.NoError(t, err)
	_, err = cache.ZAdd("testSet", &ZMember{Score: 1, Member: "testValue3"})
	assert.NoError(t, err)
	_, err = cache.ZAdd("testSet", &ZMember{Score: 1, Member: "testValue4"})
	assert.NoError(t, err)

	start, stop := int64(0), int64(3)
	members, err := cache.ZRevRangeWithScores("testSet", start, stop)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(members))
}

package cache

import (
	"context"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
)

const Empty = "redis: nil"

type Store interface {
	Set(key string, value interface{}, exp time.Duration) error
	Get(key string) (string, error)
	ZAdd(key string, member *ZMember) (int64, error)
	ZIncr(key string, member *ZMember) (float64, error)
}

type Client struct {
	client *redis.Client
}

type Options struct {
	Address  string
	Password string
	DB       int
}

type ZMember struct {
	Score  float64
	Member interface{}
}

func NewCache(options Options) (*Client, error) {
	client := redis.NewClient(&redis.Options{Addr: options.Address})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()

	return &Client{client: client}, err
}

func NewTestCache() (*Client, error) {
	mr, err := miniredis.Run()
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	return &Client{client: client}, err
}

func (c *Client) Set(key string, value interface{}, exp time.Duration) error {
	ctx := context.Background()
	return c.client.Set(ctx, key, value, exp).Err()
}

func (c *Client) Get(key string) (string, error) {
	ctx := context.Background()
	get, err := c.client.Get(ctx, key).Result()
	return get, err
}

func (c *Client) ZAdd(key string, member *ZMember) (int64, error) {
	ctx := context.Background()
	count, err := c.client.ZAdd(ctx, key, &redis.Z{Score: member.Score, Member: member.Member}).Result()
	return count, err
}

func (c *Client) ZIncr(key string, member *ZMember) (float64, error) {
	ctx := context.Background()
	count, err := c.client.ZIncr(ctx, key, &redis.Z{Score: member.Score, Member: member.Member}).Result()
	return count, err
}

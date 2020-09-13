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
	ZRevRangeWithScores(key string, start, stop int64) ([]ZMember, error)
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
	return c.client.Get(ctx, key).Result()
}

func (c *Client) ZAdd(key string, member *ZMember) (int64, error) {
	ctx := context.Background()
	return c.client.ZAdd(ctx, key, &redis.Z{Score: member.Score, Member: member.Member}).Result()
}

func (c *Client) ZIncr(key string, member *ZMember) (float64, error) {
	ctx := context.Background()
	return c.client.ZIncr(ctx, key, &redis.Z{Score: member.Score, Member: member.Member}).Result()
}

func (c *Client) ZRevRangeWithScores(key string, start, stop int64) ([]ZMember, error) {
	ctx := context.Background()
	members, err := c.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	res := make([]ZMember, 0, len(members))
	for _, v := range members {
		if v.Member != nil {
			res = append(res, ZMember{Score: v.Score, Member: v.Member})
		}
	}

	return res, err
}

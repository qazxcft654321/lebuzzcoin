package cache

import (
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

type Store interface {
	Set(key string, value interface{}, exp time.Duration) error
	Get(key string) (string, error)
}

type Client struct {
	client *redis.Client
}

type Options struct {
	Address  string
	Password string
	DB       int
}

func NewCache(options Options) (*Client, error) {
	client := redis.NewClient(&redis.Options{Addr: options.Address})
	_, err := client.Ping().Result()

	return &Client{client: client}, err
}

func NewTestCache() (*Client, error) {
	mr, err := miniredis.Run()
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	return &Client{client: client}, err
}

func (c *Client) Set(key string, value interface{}, exp time.Duration) error {
	return c.client.Set(key, value, exp).Err()
}
func (c *Client) Get(key string) (string, error) {
	get, err := c.client.Get(key).Result()
	return get, err
}

package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Redis interface {
	Client() *redis.Client
	Close() error
}

type Client struct {
	Cli *redis.Client
}

func NewRedis(conf *Config) (Redis, error) {
	err := conf.Validate()
	if err != nil {
		return &Client{}, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password, // no password set
		DB:       conf.DB,       // use default DB
	})

	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Init Redis Success \n")
		return &Client{client}, err
	}
}

func (c *Client) Client() *redis.Client {
	return c.Cli
}

func (c *Client) Close() error {
	return c.Cli.Close()
}

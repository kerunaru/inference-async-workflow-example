package infrastructure

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    *context.Context
}

func NewRedisClient(addr string, password string, db int) *RedisClient {
	log.Println(fmt.Sprintf("Connecting to %s", addr))

	ctx := context.Background()

	c := RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
		ctx: &ctx,
	}

	return &c
}

func (r *RedisClient) Get() (string, error) {
	record, err := r.client.LPop(*r.ctx, "default").Result()
	if err != nil {
		return "", err
	}

	return record, nil
}

func (r *RedisClient) Set(value string) error {
	_, err := r.client.LPush(*r.ctx, "default", value).Result()
	if err != nil {
		return err
	}

	return nil
}

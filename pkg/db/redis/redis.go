package redis

import (
	"context"
	"fmt"

	"github.com/SultanKs4/wassistant/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.Db,
	})
	cmdStat := rdb.Ping(context.Background())
	fmt.Println("redis: ", cmdStat.Val())
	return rdb, nil
}

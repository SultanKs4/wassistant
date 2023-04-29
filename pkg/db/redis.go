package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedis(addr, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	cmdStat := rdb.Ping(context.Background())
	fmt.Println("redis: ", cmdStat.Val())
	return rdb, nil
}

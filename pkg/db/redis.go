package db

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
)

type Redis struct {
	Rdb *redis.Client
}

func NewRedis(rdb *redis.Client) *Redis {
	return &Redis{Rdb: rdb}
}

func CreateDbRedis() (*redis.Client, error) {
	dbInt, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbInt,
	})
	cmdStat := rdb.Ping(context.Background())
	fmt.Println("redis: ", cmdStat.Val())
	return rdb, nil
}

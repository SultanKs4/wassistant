package entity

import (
	"context"
)

type Cache struct {
	message
}

type cache interface {
	GetCache(ctx context.Context, key string) (res *Cache, err error)
	SetCache(ctx context.Context, key string, value *Cache, timeout int) error
}

type CacheRepository interface {
	cache
}

type CacheService interface {
	cache
}

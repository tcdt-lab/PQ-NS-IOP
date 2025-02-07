package data_access

import (
	"gateway/data"
	"github.com/redis/go-redis/v9"
)

type CacheHandlerDA struct {
	client *redis.Client
}

func NewCacheHandlerDA() *CacheHandlerDA {
	client, err := getRedisConnection()
	if err != nil {
		return nil
	}
	return &CacheHandlerDA{
		client: client,
	}
}

func (c *CacheHandlerDA) GetUserAdminId(key string) (int, error) {
	cacheHandler := data.CacheHandler{}

	res, err := cacheHandler.GetUserAdminId(c.client, key)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (c *CacheHandlerDA) SetUserAdminId(key string, value int) error {
	cacheHandler := data.CacheHandler{}

	err := cacheHandler.SetUserAdminId(c.client, key, value)
	if err != nil {
		return err
	}
	return nil
}

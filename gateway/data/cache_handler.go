package data

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
)

type CacheHandler struct {
}

func (c *CacheHandler) GetUserAdminId(client *redis.Client, key string) (int, error) {
	ctx := context.Background()
	val, err := client.Get(ctx, key).Int()
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (c *CacheHandler) SetUserAdminId(client *redis.Client, key string, value int) error {
	ctx := context.Background()
	_, err := client.Set(ctx, key, value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandler) GetBootstrapVerifierId(client *redis.Client, key string) (int, error) {
	ctx := context.Background()
	val, err := client.Get(ctx, key).Int()
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (c *CacheHandler) SetBootstrapVerifierId(client *redis.Client, key string, value int) error {
	ctx := context.Background()
	_, err := client.Set(ctx, key, value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// returns request stored data title and stored data
func (c *CacheHandler) GetRequestInformation(client *redis.Client, requestId string) (string, string, error) {
	ctx := context.Background()
	data, err := client.HGet(ctx, requestId, "data").Result()
	title, err := client.HGet(ctx, requestId, "title").Result()
	if err != nil {
		return title, data, err
	}
	return title, data, nil
}

// sets request stored data title and stored data
func (c *CacheHandler) SetRequestInformation(client *redis.Client, requestId string, title string, data string) error {
	ctx := context.Background()
	_, err := client.HSet(ctx, requestId, "title", title, "data", data).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandler) SetRequestNumber(client *redis.Client, value int64) error {
	ctx := context.Background()
	_, err := client.Set(ctx, "request_number", value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandler) GetRequestNumber(client *redis.Client) (int64, error) {
	ctx := context.Background()
	val, err := client.Get(ctx, "request_number").Int64()
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (c *CacheHandler) GenerateRequestNumber(client *redis.Client) (int64, error) {
	ctx := context.Background()
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	newReqId, err := client.Incr(ctx, "request_number").Result()
	if err != nil {
		return -1, err
	}
	return newReqId, nil
}

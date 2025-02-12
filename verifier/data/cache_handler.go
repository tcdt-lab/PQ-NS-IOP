package data

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
)

type CacheHandler struct {
}

func (c *CacheHandler) GetUserAdminId(client *redis.Client, key string) (int64, error) {
	ctx := context.Background()
	val, err := client.Get(ctx, key).Int64()
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (c *CacheHandler) SetUserAdminId(client *redis.Client, key string, value int64) error {
	ctx := context.Background()
	_, err := client.Set(ctx, key, value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandler) GetBootstrapVerifierId(client *redis.Client, key string) (int64, error) {
	ctx := context.Background()
	val, err := client.Get(ctx, key).Int64()
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (c *CacheHandler) SetBootstrapVerifierId(client *redis.Client, key string, value int64) error {
	ctx := context.Background()
	_, err := client.Set(ctx, key, value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// returns request stored data title and stored data
func (c *CacheHandler) GetRequestInformation(client *redis.Client, requestId string, key string) (string, error) {
	ctx := context.Background()
	data, err := client.HGet(ctx, requestId, key).Result()

	if err != nil {
		return "", err
	}
	return data, nil
}

// sets request stored data title and stored data
func (c *CacheHandler) SetRequestInformation(client *redis.Client, requestId string, key string, data string) error {
	ctx := context.Background()
	_, err := client.HSet(ctx, requestId, key, data).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandler) RemoveRequestInformation(client *redis.Client, requestId string) error {
	ctx := context.Background()
	_, err := client.Del(ctx, requestId).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandler) RemoveDataFromRequestInformation(client *redis.Client, requestId string, key string) error {
	ctx := context.Background()
	_, err := client.HDel(ctx, requestId, key).Result()
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

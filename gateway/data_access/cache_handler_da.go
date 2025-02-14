package data_access

import (
	"gateway/data"
	"github.com/redis/go-redis/v9"
	"strconv"
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

func (c *CacheHandlerDA) GetUserAdminId() (int64, error) {
	cacheHandler := data.CacheHandler{}

	res, err := cacheHandler.GetUserAdminId(c.client, "user_admin_id")
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (c *CacheHandlerDA) SetUserAdminId(value int64) error {
	cacheHandler := data.CacheHandler{}

	err := cacheHandler.SetUserAdminId(c.client, "user_admin_id", value)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandlerDA) SetBootstrapVerifierId(value int64) error {
	cacheHandler := data.CacheHandler{}

	err := cacheHandler.SetBootstrapVerifierId(c.client, "bootstrap_verifier_id", value)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandlerDA) GetBootstrapVerifierId() (int64, error) {
	cacheHandler := data.CacheHandler{}

	res, err := cacheHandler.GetBootstrapVerifierId(c.client, "bootstrap_verifier_id")
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (c *CacheHandlerDA) GetRequestInformation(requestId int64, key string) (string, error) {
	cacheHandler := data.CacheHandler{}

	data, err := cacheHandler.GetRequestInformation(c.client, strconv.FormatInt(requestId, 10), key)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (c *CacheHandlerDA) SetRequestInformation(requestId int64, title string, reqData string) error {
	cacheHandler := data.CacheHandler{}

	err := cacheHandler.SetRequestInformation(c.client, strconv.FormatInt(requestId, 10), title, reqData)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandlerDA) RemoveRequestInformation(requestId int64) error {
	cacheHandler := data.CacheHandler{}

	err := cacheHandler.RemoveRequestInformation(c.client, strconv.FormatInt(requestId, 10))
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheHandlerDA) RemoveTitleDataFromRequestInformation(requestId int64, key string) error {
	cacheHandler := data.CacheHandler{}

	err := cacheHandler.RemoveDataFromRequestInformation(c.client, strconv.FormatInt(requestId, 10), key)
	if err != nil {
		return err
	}
	return nil

}

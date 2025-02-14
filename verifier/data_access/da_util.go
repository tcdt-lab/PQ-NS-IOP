package data_access

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"verifier/config"
)

func getConfig() (config.Config, error) {
	cfg, err := config.ReadYaml()
	return *cfg, err
}

func getDbConnection() (*sql.DB, error) {
	c, err := getConfig()
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)

	if err != nil {
		return nil, err
	}
	return db, nil

}

func getRedisConnection() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       1,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

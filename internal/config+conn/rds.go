package config_conn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func RedisConn() *redis.Client {
	cfg := GetConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opt, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		log.Fatalf("failed to parse redis url: %v", err)
	}

	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return client
}

func InvalidateCache(userID uuid.UUID, client *redis.Client) error {
	pattern := fmt.Sprintf("Api:%v*", userID)
	var keys []string
	iter := client.Scan(context.Background(), 0, pattern, 0).Iterator()
	for iter.Next(context.Background()) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to iterate redis keys: %v", err)
	}
	if len(keys) == 0 {
		return nil
	}
	if err := client.Del(context.Background(), keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete redis keys: %v", err)
	}
	return nil
}

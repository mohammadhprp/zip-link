package configs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func ConnectRedisDB() *RedisClient {
	redisURL := os.Getenv("REDIS_URL")
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis ping failed: %v", err)
	}

	log.Println("Connected to Redis successfully")
	return &RedisClient{Client: client}
}

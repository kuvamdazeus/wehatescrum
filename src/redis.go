package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client = nil

func getRedisClient() (*redis.Client, error) {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return redisClient, err
	}

	if redisClient == nil {
		redisClient = redis.NewClient(opts)
	}

	return redisClient, nil
}
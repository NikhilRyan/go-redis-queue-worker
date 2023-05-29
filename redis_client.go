package main

import "github.com/go-redis/redis/v8"

// NewRedisClient creates a new Redis client instance
func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Redis password
		DB:       0,                // Redis database index
	})
	return client
}

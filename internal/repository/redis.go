// redis.go
package repository

import (
	"github.com/go-redis/redis/v8"
)

// NewClient creates a new Redis client and returns it
func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

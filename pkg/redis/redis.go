package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client    *redis.Client
	Ctx       = context.Background()
	IsEnabled bool
	TTL       time.Duration
)

// InitRedis initializes Redis connection from environment variables
func InitRedis() error {
	// Check if Redis is enabled
	enabled := os.Getenv("REDIS_ENABLED")
	IsEnabled = enabled == "true"

	if !IsEnabled {
		log.Println("[REDIS] Redis is disabled in configuration")
		return nil
	}

	// Get Redis configuration from environment
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	unixSocket := os.Getenv("REDIS_UNIX_SOCKET")
	password := os.Getenv("REDIS_PASSWORD")

	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		if parsedDB, err := strconv.Atoi(dbStr); err == nil {
			db = parsedDB
		}
	}

	// Get TTL configuration
	ttlStr := os.Getenv("REDIS_TTL")
	TTL = 3600 * time.Second // Default 1 hour
	if ttlStr != "" {
		if parsedTTL, err := strconv.Atoi(ttlStr); err == nil {
			TTL = time.Duration(parsedTTL) * time.Second
		}
	}

	// Create Redis client with Unix socket or TCP connection
	var options *redis.Options
	if unixSocket != "" {
		// Use Unix socket connection
		options = &redis.Options{
			Network:  "unix",
			Addr:     unixSocket,
			Password: password,
			DB:       db,
		}
		log.Printf("[REDIS] Connecting to Redis via Unix socket: %s", unixSocket)
	} else {
		// Use TCP connection
		options = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       db,
		}
		log.Printf("[REDIS] Connecting to Redis via TCP: %s:%s", host, port)
	}

	Client = redis.NewClient(options)

	// Test connection
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Printf("[REDIS] Failed to connect to Redis: %v", err)
		IsEnabled = false
		return err
	}

	if unixSocket != "" {
		log.Printf("[REDIS] Successfully connected to Redis via Unix socket: %s (DB: %d)", unixSocket, db)
	} else {
		log.Printf("[REDIS] Successfully connected to Redis via TCP: %s:%s (DB: %d)", host, port, db)
	}
	log.Printf("[REDIS] Cache TTL set to %v", TTL)
	return nil
}

// Set stores a key-value pair in Redis with TTL
func Set(key string, value interface{}) error {
	if !IsEnabled || Client == nil {
		return nil
	}

	err := Client.Set(Ctx, key, value, TTL).Err()
	if err != nil {
		log.Printf("[REDIS] Error setting key %s: %v", key, err)
	}
	return err
}

// Get retrieves a value from Redis by key
func Get(key string) (string, error) {
	if !IsEnabled || Client == nil {
		return "", fmt.Errorf("redis not enabled")
	}

	val, err := Client.Get(Ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	}
	if err != nil {
		log.Printf("[REDIS] Error getting key %s: %v", key, err)
		return "", err
	}
	return val, nil
}

// Delete removes a key from Redis
func Delete(key string) error {
	if !IsEnabled || Client == nil {
		return nil
	}

	err := Client.Del(Ctx, key).Err()
	if err != nil {
		log.Printf("[REDIS] Error deleting key %s: %v", key, err)
	}
	return err
}

// DeletePattern removes all keys matching a pattern
func DeletePattern(pattern string) error {
	if !IsEnabled || Client == nil {
		return nil
	}

	keys, err := Client.Keys(Ctx, pattern).Result()
	if err != nil {
		log.Printf("[REDIS] Error finding keys with pattern %s: %v", pattern, err)
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	err = Client.Del(Ctx, keys...).Err()
	if err != nil {
		log.Printf("[REDIS] Error deleting keys with pattern %s: %v", pattern, err)
	} else {
		log.Printf("[REDIS] Deleted %d keys with pattern %s", len(keys), pattern)
	}
	return err
}

// FlushAll clears all Redis data (use with caution)
func FlushAll() error {
	if !IsEnabled || Client == nil {
		return nil
	}

	err := Client.FlushAll(Ctx).Err()
	if err != nil {
		log.Printf("[REDIS] Error flushing all data: %v", err)
	} else {
		log.Println("[REDIS] All Redis data flushed")
	}
	return err
}

// Close closes the Redis connection
func Close() error {
	if Client != nil {
		log.Println("[REDIS] Closing Redis connection")
		return Client.Close()
	}
	return nil
}

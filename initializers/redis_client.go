package initializers

import (
    "context"
    "github.com/redis/go-redis/v9"
    "log"
    "os"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),    // Use REDIS_ADDR from environment variables
        Password: os.Getenv("REDIS_PASSWORD"), // No password set by default
        DB:       0,                          // Default DB
    })

    // Test the connection
    _, err := RedisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Println("Connected to Redis")
}

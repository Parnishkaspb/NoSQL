package redishelper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	ctx      = context.Background()
	initOnce sync.Once
	clients  = make(map[int]*redis.Client)
)

func getClient(dbNumber int) (*redis.Client, error) {
	initOnce.Do(func() {
		_ = godotenv.Load(".env") // безопасно для продакшена
	})

	if client, ok := clients[dbNumber]; ok {
		return client, nil
	}

	sentinels := os.Getenv("REDIS_SENTINELS")
	masterName := os.Getenv("REDIS_MASTER_NAME")

	if sentinels == "" || masterName == "" {
		return nil, fmt.Errorf("❌ переменные REDIS_SENTINELS или REDIS_MASTER_NAME не заданы")
	}

	addrs := strings.Split(sentinels, ",")

	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addrs,
		DB:            dbNumber,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if _, err := client.Ping(pingCtx).Result(); err != nil {
		return nil, fmt.Errorf("Redis Sentinel не отвечает: %v", err)
	}

	clients[dbNumber] = client
	log.Printf("✅ Redis подключен: DB %d", dbNumber)
	return client, nil
}

func Read[T any](dbNumber int, key string) ([]T, error) {
	rb, err := getClient(dbNumber)
	if err != nil {
		return nil, err
	}

	val, err := rb.Get(ctx, key).Result()
	if err == redis.Nil {
		return make([]T, 0), nil
	} else if err != nil {
		return nil, err
	}

	var items []T
	if err := json.Unmarshal([]byte(val), &items); err != nil {
		return nil, err
	}

	return items, nil
}

func CreateUpdate[T any](dbNumber int, key string, value T) error {
	rb, err := getClient(dbNumber)
	if err != nil {
		return err
	}

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rb.Set(ctx, key, valueJSON, 0).Err()
}

func Delete[T any](dbNumber int, key string) error {
	rb, err := getClient(dbNumber)
	if err != nil {
		return err
	}

	return rb.Del(ctx, key).Err()
}

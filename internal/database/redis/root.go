package redishelper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	client   *redis.Client
	initOnce sync.Once
	initErr  error
	ctx      = context.Background()
)

func initRedis(typeDB int) {
	initErr = godotenv.Load(".env")
	if initErr != nil {
		log.Printf("не удалось загрузить .env: %v", initErr)
		return
	}

	host := os.Getenv("HOST")
	port := os.Getenv("REDIS_PORT")
	if host == "" || port == "" {
		initErr = fmt.Errorf("HOST или REDIS_PORT не заданы в .env")
		log.Printf("HOST или REDIS_PORT не заданы: %s %s", host, port)
		return
	}

	client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       typeDB,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if _, err := client.Ping(pingCtx).Result(); err != nil {
		initErr = fmt.Errorf("ошибка подключения к Redis: %v", err)
		log.Printf("Redis ошибка: %v", err)
		return
	}

	log.Println("✅ Успешное подключение к Redis")
}

func getClient(typeDB int) (*redis.Client, error) {
	initOnce.Do(func() {
		initRedis(typeDB)
	})
	return client, initErr
}

func Read[T any](typeDB int, key string) ([]T, error) {
	rb, err := getClient(typeDB)
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

func CreateUpdate[T any](typeDB int, key string, value T) error {
	rb, err := getClient(typeDB)
	if err != nil {
		return err
	}

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rb.Set(ctx, key, valueJSON, 0).Err()
}

func Delete[T any](typeDB int, key string) error {
	rb, err := getClient(typeDB)
	if err != nil {
		return err
	}

	return rb.Del(ctx, key).Err()
}

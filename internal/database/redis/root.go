package redishelper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

func connectRedis(typeDB int) (*redis.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("не удалось загрузить .env: %v", err)
		return nil, fmt.Errorf("не удалось загрузить .env: %v", err)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("REDIS_PORT")
	if host == "" || port == "" {
		return nil, fmt.Errorf("HOST или REDIS_PORT не заданы в .env")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       typeDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("ошибка подключения к Redis: %v\n", err)
		return nil, fmt.Errorf("ошибка подключения к Redis: %v", err)
	}

	log.Println("✅ Успешное подключение к Redis")
	return rdb, nil
}

func Read[T any](typeDB int, key string) ([]T, error) {
	rb, err := connectRedis(typeDB)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	val, err := rb.Get(ctx, key).Result()
	if err == redis.Nil {
		return make([]T, 0), err
	} else if err != nil {
		return nil, err
	}

	var items []T
	err = json.Unmarshal([]byte(val), &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func CreateUpdate[T any](typeDB int, key string, value T) error {
	ctx := context.Background()
	rb, err := connectRedis(typeDB)
	if err != nil {
		return err
	}

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rb.Set(ctx, key, valueJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func Delete(typeDB int, key string) error {
	ctx := context.Background()
	rb, err := connectRedis(typeDB)
	if err != nil {
		return err
	}
	err = rb.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

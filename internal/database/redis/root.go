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
	clients  = make(map[int]*redis.Client)
	initOnce sync.Once
	initErr  error
	ctx      = context.Background()
)

func initRedis() {
	_ = godotenv.Load(".env") // .env может не быть в проде — не ошибка

	sentinels := os.Getenv("REDIS_SENTINELS")    // пример: "redis-sentinel1:26379,redis-sentinel2:26379,redis-sentinel3:26379"
	masterName := os.Getenv("REDIS_MASTER_NAME") // по умолчанию "mymaster"
	if sentinels == "" || masterName == "" {
		initErr = fmt.Errorf("не заданы REDIS_SENTINELS или REDIS_MASTER_NAME в .env")
		log.Println("❌", initErr)
		return
	}

	log.Println("✅ Конфигурация Redis Sentinel загружена")
}

func getClient(typeDB int) (*redis.Client, error) {
	initOnce.Do(initRedis)

	// Если уже есть клиент для этой БД — возвращаем
	if c, ok := clients[typeDB]; ok {
		return c, nil
	}

	sentinels := os.Getenv("REDIS_SENTINELS")
	masterName := os.Getenv("REDIS_MASTER_NAME")
	password := os.Getenv("REDIS_PASSWORD")

	addrs := splitAndTrim(sentinels)

	opts := &redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addrs,
		DB:            typeDB,
		Password:      password,
	}

	c := redis.NewFailoverClient(opts)

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := c.Ping(ctxPing).Err(); err != nil {
		return nil, fmt.Errorf("Redis Sentinel не отвечает: %w", err)
	}

	clients[typeDB] = c
	return c, nil
}

func splitAndTrim(s string) []string {
	var result []string
	for _, part := range splitComma(s) {
		result = append(result, trim(part))
	}
	return result
}

func splitComma(s string) []string {
	var parts []string
	for _, p := range []rune(s) {
		parts = append(parts, string(p))
	}
	return parts
}

func trim(s string) string {
	return string([]byte(s))
}

//var (
//	client   *redis.Client
//	initOnce sync.Once
//	initErr  error
//	ctx      = context.Background()
//)
//
//func initRedis(typeDB int) {
//	initErr = godotenv.Load(".env")
//	if initErr != nil {
//		log.Printf("не удалось загрузить .env: %v", initErr)
//		return
//	}
//
//	host := os.Getenv("REDIS_HOST")
//	port := os.Getenv("REDIS_PORT")
//	if host == "" || port == "" {
//		initErr = fmt.Errorf("HOST или REDIS_PORT не заданы в .env")
//		log.Printf("HOST или REDIS_PORT не заданы: %s %s", host, port)
//		return
//	}
//
//	client = redis.NewClient(&redis.Options{
//		Addr:     host + ":" + port,
//		Password: "",
//		DB:       typeDB,
//	})
//
//	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
//	defer cancel()
//
//	if _, err := client.Ping(pingCtx).Result(); err != nil {
//		initErr = fmt.Errorf("ошибка подключения к Redis: %v", err)
//		log.Printf("Redis ошибка: %v", err)
//		return
//	}
//
//	log.Println("✅ Успешное подключение к Redis")
//}
//
//func getClient(typeDB int) (*redis.Client, error) {
//	initOnce.Do(func() {
//		initRedis(typeDB)
//	})
//	return client, initErr
//}

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

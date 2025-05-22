package mongohelper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client    *mongo.Client
	database  *mongo.Database
	initOnce  sync.Once
	initErr   error
	ctx       = context.Background()
	defaultDB string
)

func initMongo() {
	initErr = godotenv.Load(".env")
	if initErr != nil {
		log.Printf("Не удалось загрузить .env: %v", initErr)
		return
	}

	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	defaultDB = os.Getenv("MONGO_DB")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")

	if username == "" || password == "" || defaultDB == "" || host == "" || port == "" {
		log.Printf("Не заданы обязательные переменные окружения: %s %s %s %s", username, password, host, port)
		initErr = fmt.Errorf("Не заданы обязательные переменные окружения")
		return
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	clientOpts := options.Client().ApplyURI(uri)

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, initErr = mongo.Connect(connectCtx, clientOpts)
	if initErr != nil {
		log.Printf("Ошибка подключения к MongoDB: %v", initErr)
		return
	}

	if err := client.Ping(connectCtx, nil); err != nil {
		log.Printf("MongoDB не отвечает: %v", err)
		initErr = fmt.Errorf("MongoDB не отвечает: %v", err)
		return
	}

	database = client.Database(defaultDB)
	log.Println("✅ Успешное подключение к MongoDB")
}

func getDB() (*mongo.Database, error) {
	initOnce.Do(initMongo)
	return database, initErr
}

func ReadAll[T any](collectionName string) ([]T, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}

	cursor, err := db.Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var item T
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, cursor.Err()
}

func Read[T any](collectionName string, filter any) ([]T, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}

	cursor, err := db.Collection(collectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var item T
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, cursor.Err()
}

func Create[T any](collectionName string, value T) (string, error) {
	db, err := getDB()
	if err != nil {
		return "", err
	}

	res, err := db.Collection(collectionName).InsertOne(ctx, value)
	if err != nil {
		return "", err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "inserted (non-objectID)", nil
}

func Delete[T any](collectionName string, filter bson.M) (int64, error) {
	db, err := getDB()
	if err != nil {
		return 0, err
	}

	res, err := db.Collection(collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

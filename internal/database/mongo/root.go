package mongohelper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	_ = godotenv.Load(".env")

	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	defaultDB = os.Getenv("MONGO_DB")
	replicaSet := os.Getenv("MONGO_REPLICA_SET")
	hosts := os.Getenv("MONGO_HOSTS")

	if username == "" || password == "" || defaultDB == "" || hosts == "" || replicaSet == "" {
		initErr = fmt.Errorf("не заданы обязательные переменные окружения: MONGO_USERNAME, MONGO_PASSWORD, MONGO_DB, MONGO_HOSTS, MONGO_REPLICA_SET")
		log.Println("❌", initErr)
		return
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s/?replicaSet=%s", username, password, hosts, replicaSet)

	clientOpts := options.Client().
		ApplyURI(uri).
		SetReadPreference(readpref.SecondaryPreferred()).
		SetRetryWrites(true)

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, initErr = mongo.Connect(connectCtx, clientOpts)
	if initErr != nil {
		log.Printf("❌ Ошибка подключения к MongoDB: %v", initErr)
		return
	}

	if err := client.Ping(connectCtx, nil); err != nil {
		initErr = fmt.Errorf("MongoDB не отвечает: %v", err)
		log.Println("❌", initErr)
		return
	}

	database = client.Database(defaultDB)
	log.Println("✅ Подключение к MongoDB Replica Set успешно")
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

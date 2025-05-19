package mongohelper

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func connectMongo() (*mongo.Client, *mongo.Database, context.Context, context.CancelFunc, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	database := os.Getenv("MONGO_DB")
	host := os.Getenv("HOST")
	port := os.Getenv("MONGO_PORT")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		fmt.Println("Не удалось подключиться к MongoDB:", err)
		return nil, nil, nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		cancel()
		fmt.Println("MongoDB не отвечает:", err)
		return nil, nil, nil, nil, err
	}

	fmt.Println("✅ Успешное подключение к MongoDB")

	return client, client.Database(database), ctx, cancel, nil
}

func ReadAll[T any](collectionName string) ([]T, error) {
	client, db, ctx, cancel, err := connectMongo()
	if err != nil {
		cancel()
		return nil, err
	}
	defer cancel()
	defer client.Disconnect(ctx)

	collection := db.Collection(collectionName)
	cursor, err := collection.Find(ctx, bson.M{})
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

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func Create[T any](collectionName string, value T) (string, error) {
	client, db, ctx, cancel, err := connectMongo()
	if err != nil {
		return "Проблема с подключением к БД", err
	}
	defer cancel()
	defer client.Disconnect(ctx)

	collection := db.Collection(collectionName)

	res, err := collection.InsertOne(ctx, value)
	if err != nil {
		return "Проблема с вставкой", err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return "Запись успешно создана: " + oid.Hex(), nil
	}

	return "Запись создана (не ObjectID)", nil
}

func Read[T any](collectionName string, filter any) ([]T, error) {
	client, db, ctx, cancel, err := connectMongo()
	if err != nil {
		cancel()
		return nil, err
	}
	defer cancel()
	defer client.Disconnect(ctx)

	collection := db.Collection(collectionName)
	cursor, err := collection.Find(ctx, filter)
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

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func DeleteT(collectionName string, filter bson.M) (int64, error) {
	client, db, ctx, cancel, err := connectMongo()
	if err != nil {
		return 0, err
	}
	defer cancel()
	defer client.Disconnect(ctx)

	res, err := db.Collection(collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

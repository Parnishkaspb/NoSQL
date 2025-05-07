package mongohelper

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func ConnectMongo() (*mongo.Client, *mongo.Database, context.Context, context.CancelFunc, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	database := os.Getenv("MONGO_DB")
	host := os.Getenv("MONGO_HOST")
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

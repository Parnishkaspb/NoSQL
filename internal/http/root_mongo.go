package http

import (
	mongohelper "NoSQL/internal/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadAll[T any](collectionName string) ([]T, error) {
	client, db, ctx, cancel, err := mongohelper.ConnectMongo()
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
	client, db, ctx, cancel, err := mongohelper.ConnectMongo()
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
	client, db, ctx, cancel, err := mongohelper.ConnectMongo()
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

//func main() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Ошибка загрузки .env файла")
//	}
//
//	username := os.Getenv("MONGO_USERNAME")
//	password := os.Getenv("MONGO_PASSWORD")
//	host := os.Getenv("MONGO_HOST")
//	port := os.Getenv("MONGO_PORT")
//	dbname := os.Getenv("MONGO_DB")
//
//	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	clientOptions := options.Client().ApplyURI(uri)
//	client, err := mongo.Connect(ctx, clientOptions)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Disconnect(ctx)
//
//	err = client.Ping(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println("Успешное подключение к MongoDB!")
//
//	collection := client.Database(dbname).Collection("users")
//	fmt.Println("Работаем с коллекцией:", collection.Name())
//}

package http

import (
	mongohelper "NoSQL/internal/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCountry(entity string) string {
	client, database, ctx, cancel, err := mongohelper.ConnectMongo()

	collection := database.Collection("countries")
	if err != nil {
		cancel()
		return "Произошла проблема с подключение к БД: " + err.Error()
	}

	country := Country{
		Name: entity,
	}

	res, err := collection.InsertOne(ctx, country)

	if err != nil {
		return "Вставка не произошла успешно: " + err.Error()
	}

	defer cancel()
	defer client.Disconnect(ctx)

	return "Запись успешно создана: " + res.InsertedID.(primitive.ObjectID).Hex()
}

func ReadCountry() ([]Country, error) {
	client, db, ctx, cancel, err := mongohelper.ConnectMongo()
	if err != nil {
		cancel()
		return nil, err
	}
	defer cancel()
	defer client.Disconnect(ctx)

	cursor, err := db.Collection("countries").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var countries []Country
	for cursor.Next(ctx) {
		var country Country
		if err := cursor.Decode(&country); err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	return countries, err
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

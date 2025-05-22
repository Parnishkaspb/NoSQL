package main

import (
	mongohelper "NoSQL/internal/database/mongo"
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Country struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `faker:"name" bson:"name"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	_ = context.Background()
	const count = 10

	for i := 0; i < count; i++ {
		user := User{
			Age: 10 + rand.Intn(40),
		}
		err := faker.FakeData(&user)
		if err != nil {
			log.Fatalf("Ошибка генерации данных: %v", err)
		}

		_, err = mongohelper.Create[User]("users", user)
		if err != nil {
			log.Printf("Ошибка сохранения пользователя %s: %v", user.Login, err)
			continue
		}

		log.Printf("✅ Добавлен пользователь: %s %s (%s)", user.Name, user.Surname, user.Login)
	}

	log.Println("🌱 Сидирование завершено!")
}

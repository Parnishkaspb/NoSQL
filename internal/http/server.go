package http

import (
	mongohelper "NoSQL/internal/database/mongo"
	redishelper "NoSQL/internal/database/redis"
	"NoSQL/internal/pkg"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func createHandler(s *Serve) {
	s.mux.HandleFunc("POST /create/{entity}", func(w http.ResponseWriter, r *http.Request) {
		entity := r.PathValue("entity")

		var tB tableBody

		if err := json.NewDecoder(r.Body).Decode(&tB); err != nil {
			if err == io.EOF {
				pkg.WriteApiResponse(w, nil, "Пустое тело запроса", http.StatusBadRequest)
				return
			}
			pkg.WriteApiResponse(w, nil, "Ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		var answer string
		var err error
		switch entity {
		case "country":
			if tB.Name == "" {
				pkg.WriteApiResponse(w, nil, "Поле name не может быть пустым", http.StatusBadRequest)
				return
			}

			answer, err = mongohelper.Create[Country]("countries", Country{Name: tB.Name})

			if err != nil {
				pkg.WriteApiResponse(w, nil, answer+err.Error(), http.StatusInternalServerError)
				return
			}
		case "platform":
			if tB.Name == "" {
				pkg.WriteApiResponse(w, nil, "Поле name не может быть пустым", http.StatusBadRequest)
				return
			}

			answer, err = mongohelper.Create[Platform]("platforms", Platform{Name: tB.Name})

			if err != nil {
				pkg.WriteApiResponse(w, nil, answer+err.Error(), http.StatusInternalServerError)
				return
			}
		case "game":
			if tB.Title == "" {
				pkg.WriteApiResponse(w, nil, "Поле title не может быть пустым", http.StatusBadRequest)
				return
			}
			if tB.Description == "" {
				pkg.WriteApiResponse(w, nil, "Поле description не может быть пустым", http.StatusBadRequest)
				return
			}
			if tB.Price <= 0 {
				pkg.WriteApiResponse(w, nil, "Поле price не может быть меньше(равно) 0", http.StatusBadRequest)
				return
			}
			if tB.Age <= 0 {
				pkg.WriteApiResponse(w, nil, "Поле age не может быть меньше(равно) 0", http.StatusBadRequest)
				return
			}

			game := Game{
				Title:       tB.Title,
				Description: tB.Description,
				Price:       tB.Price,
				Age:         tB.Age,
				ApprovedAge: tB.ApprovedAge != nil && *tB.ApprovedAge,
			}

			answer, err = mongohelper.Create[Game]("games", game)

			if err != nil {
				pkg.WriteApiResponse(w, nil, answer+err.Error(), http.StatusInternalServerError)
				return
			}
		case "user":
			if tB.Name == "" {
				pkg.WriteApiResponse(w, nil, "Поле title не может быть пустым", http.StatusBadRequest)
				return
			}
			if tB.Surname == "" {
				pkg.WriteApiResponse(w, nil, "Поле description не может быть пустым", http.StatusBadRequest)
				return
			}
			if tB.Login == "" {
				pkg.WriteApiResponse(w, nil, "Поле login не может быть пустым", http.StatusBadRequest)
				return
			}

			users, err := mongohelper.Read[User]("users", bson.M{"login": tB.Login})
			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(users) > 0 {
				pkg.WriteApiResponse(w, nil, "Пользователь с таким логином уже существует!", http.StatusConflict)
				return
			}

			if tB.Password == "" {
				pkg.WriteApiResponse(w, nil, "Поле password не может быть пустым", http.StatusBadRequest)
				return
			}

			if tB.Age <= 0 {
				pkg.WriteApiResponse(w, nil, "Поле age не может быть меньше(равно) 0", http.StatusBadRequest)
				return
			}

			user := User{
				Name:     tB.Name,
				Surname:  tB.Surname,
				Login:    tB.Login,
				Age:      tB.Age,
				Password: tB.Password,
			}

			answer, err = mongohelper.Create[User]("users", user)

			if err != nil {
				pkg.WriteApiResponse(w, nil, answer+err.Error(), http.StatusInternalServerError)
				return
			}
		case "gamereview":
			if tB.GameID.IsZero() {
				pkg.WriteApiResponse(w, nil, "Поле game_id не может быть пустым", http.StatusBadRequest)
				return
			}
			if tB.UserID.IsZero() {
				pkg.WriteApiResponse(w, nil, "Поле user_id не может быть пустым", http.StatusBadRequest)
				return
			}
			if tB.Review == "" {
				pkg.WriteApiResponse(w, nil, "Поле text не может быть пустым", http.StatusBadRequest)
				return
			}
			if tB.Rating < 1 || tB.Rating > 5 {
				pkg.WriteApiResponse(w, nil, "Рейтинг должен быть от 1 до 5", http.StatusBadRequest)
				return
			}

			users, err := mongohelper.Read[User]("users", bson.M{"_id": tB.UserID})
			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(users) == 0 {
				pkg.WriteApiResponse(w, nil, "Такого пользователя не существует!", http.StatusConflict)
				return
			}

			games, err := mongohelper.Read[Game]("games", bson.M{"_id": tB.GameID})
			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(games) == 0 {
				pkg.WriteApiResponse(w, nil, "Такой игры не существует!", http.StatusConflict)
				return
			}

			gamereview := GameReview{
				GameID:    tB.GameID,
				UserID:    tB.UserID,
				Rating:    tB.Rating,
				Review:    tB.Review,
				CreatedAt: time.Now(),
			}

			answer, err = mongohelper.Create[GameReview]("gamesreviews", gamereview)

			if err != nil {
				pkg.WriteApiResponse(w, nil, answer+err.Error(), http.StatusInternalServerError)
				return
			}
		case "gameactualprice":
			if tB.GameID.IsZero() {
				pkg.WriteApiResponse(w, nil, "Поле game_id не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}

			if tB.Price <= 0 {
				pkg.WriteApiResponse(w, nil, "Поле price не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}

			if tB.Discount <= 0 {
				pkg.WriteApiResponse(w, nil, "Поле discount не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}
			typeDB, err := strconv.Atoi(os.Getenv("REDIS_GamesActualPrice"))
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемы с типом БД", http.StatusBadRequest)
			}
			gameactualprice := GameActualPrice{
				Price:    tB.Price,
				Discount: tB.Discount,
			}
			err = redishelper.CreateUpdate[GameActualPrice](typeDB, tB.GameID.Hex(), gameactualprice)

			if err != nil {
				pkg.WriteApiResponse(w, nil, answer+err.Error(), http.StatusInternalServerError)
				return
			}
		case "userscart":
			if tB.UserID.IsZero() {
				pkg.WriteApiResponse(w, nil, "Поле user_id не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}

			if tB.GameID.IsZero() {
				pkg.WriteApiResponse(w, nil, "Поле game_id не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}

			if tB.Quantity <= 0 {
				pkg.WriteApiResponse(w, nil, "Поле discount не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}

			if tB.Price <= 0 {
				pkg.WriteApiResponse(w, nil, "Поле price не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}

			typeDB, err := strconv.Atoi(os.Getenv("REDIS_UsersCart"))
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемы с типом БД", http.StatusBadRequest)
				return
			}
			userscart := UsersCart{
				GameID:   tB.GameID,
				Quantity: tB.Quantity,
				Price:    tB.Price,
			}

			err = redishelper.CreateUpdate[UsersCart](typeDB, tB.GameID.Hex(), userscart)

			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемa : "+err.Error(), http.StatusBadRequest)
				return
			}

			pkg.WriteApiResponse(w, userscart, "", http.StatusOK)
			return

		}

		pkg.WriteApiResponse(w, nil, "Прочитано "+entity, http.StatusOK)
		return
	})
}

func readHandler(s *Serve) {
	s.mux.HandleFunc("GET /read/{entity}/{id}", func(w http.ResponseWriter, r *http.Request) {
		entity := r.PathValue("entity")
		id := r.PathValue("id")
		switch entity {
		case "country":
			answer, err := mongohelper.ReadAll[Country]("countries")

			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, answer, "", http.StatusOK)
			return
		case "platform":
			answer, err := mongohelper.ReadAll[Platform]("platforms")
			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, answer, "", http.StatusOK)
			return
		case "game":
			answer, err := mongohelper.ReadAll[Game]("games")
			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, answer, "", http.StatusOK)
			return
		case "user":
			answer, err := mongohelper.ReadAll[User]("users")
			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, answer, "", http.StatusOK)
			return
		case "gamereview":
			answer, err := mongohelper.ReadAll[GameReview]("gamesreviews")
			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, answer, "", http.StatusOK)
			return
		case "gameactualprice":
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Некорректный ID", http.StatusBadRequest)
				return
			}
			typeDB, err := strconv.Atoi(os.Getenv("REDIS_GamesActualPrice"))
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемы с типом БД", http.StatusBadRequest)
				return
			}
			gameactualprice, err := redishelper.Read[GameActualPrice](typeDB, id.Hex())

			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемa : "+err.Error(), http.StatusBadRequest)
				return
			}

			pkg.WriteApiResponse(w, gameactualprice, "", http.StatusOK)
			return
		case "userscart":
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Некорректный ID", http.StatusBadRequest)
				return
			}
			typeDB, err := strconv.Atoi(os.Getenv("REDIS_UsersCart"))
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемы с типом БД", http.StatusBadRequest)
				return
			}
			userscart, err := redishelper.Read[UsersCart](typeDB, id.Hex())

			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемa : "+err.Error(), http.StatusBadRequest)
				return
			}

			pkg.WriteApiResponse(w, userscart, "", http.StatusOK)
			return
		}
		pkg.WriteApiResponse(w, nil, "Прочитано: "+entity, http.StatusOK)
		return
	})
}
func updateHandler(s *Serve) {
	s.mux.HandleFunc("PUT /update/{entity}/{id}", func(w http.ResponseWriter, r *http.Request) {
		entity := r.PathValue("entity")
		id := r.PathValue("id")
		var tB tableBody

		if err := json.NewDecoder(r.Body).Decode(&tB); err != nil {
			if err == io.EOF {
				pkg.WriteApiResponse(w, nil, "Пустое тело запроса", http.StatusBadRequest)
				return
			}
			pkg.WriteApiResponse(w, nil, "Ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch entity {
		case "country":
			answer, err := mongohelper.ReadAll[Country]("countries")

			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, answer, "", http.StatusOK)
			return
		case "platform":
			answer, err := mongohelper.ReadAll[Platform]("platforms")
			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, answer, "", http.StatusOK)
			return
		case "userscart":
			if tB.GameID.IsZero() {
				pkg.WriteApiResponse(w, nil, "Поле game_id не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}

			if tB.Quantity <= 0 {
				pkg.WriteApiResponse(w, nil, "Поле discount не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}

			if tB.Price <= 0 {
				pkg.WriteApiResponse(w, nil, "Поле price не может быть равно или меньше 0", http.StatusBadRequest)
				return
			}

			typeDB, err := strconv.Atoi(os.Getenv("REDIS_UsersCart"))
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемы с типом БД", http.StatusBadRequest)
				return
			}
			userscart := UsersCart{
				GameID:   tB.GameID,
				Quantity: tB.Quantity,
				Price:    tB.Price,
			}

			err = redishelper.CreateUpdate[UsersCart](typeDB, id, userscart)

			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемa : "+err.Error(), http.StatusBadRequest)
				return
			}

			pkg.WriteApiResponse(w, userscart, "", http.StatusOK)
			return
		}
		pkg.WriteApiResponse(w, nil, "Прочитано: "+entity, http.StatusOK)
		return
	})
}
func deleteHandler(s *Serve) {
	s.mux.HandleFunc("DELETE /delete/{entity}/{id}", func(w http.ResponseWriter, r *http.Request) {
		entity := r.PathValue("entity")
		id := r.PathValue("id")
		switch entity {
		case "country":
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Некорректный ID", http.StatusBadRequest)
				return
			}
			_, err = mongohelper.DeleteT("countries", bson.M{"_id": id})

			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, "Удаление произошло успешно", "", http.StatusOK)
			return
		case "platform":
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Некорректный ID", http.StatusBadRequest)
				return
			}
			_, err = mongohelper.DeleteT("platforms", bson.M{"_id": id})

			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, "Удаление произошло успешно", "", http.StatusOK)
			return
		case "game":
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Некорректный ID", http.StatusBadRequest)
				return
			}
			_, err = mongohelper.DeleteT("games", bson.M{"_id": id})

			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, "Удаление произошло успешно", "", http.StatusOK)
			return
		case "user":
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Некорректный ID", http.StatusBadRequest)
				return
			}
			_, err = mongohelper.DeleteT("users", bson.M{"_id": id})

			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, "Удаление произошло успешно", "", http.StatusOK)
			return
		case "gamereview":
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Некорректный ID", http.StatusBadRequest)
				return
			}
			_, err = mongohelper.DeleteT("gamesreviews", bson.M{"_id": id})

			if err != nil {
				pkg.WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			pkg.WriteApiResponse(w, "Удаление произошло успешно", "", http.StatusOK)
			return
		case "gameactualprice":
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Некорректный ID", http.StatusBadRequest)
				return
			}
			typeDB, err := strconv.Atoi(os.Getenv("REDIS_GamesActualPrice"))
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемы с типом БД", http.StatusBadRequest)
			}
			err = redishelper.Delete(typeDB, id.Hex())
		case "userscart":
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Некорректный ID", http.StatusBadRequest)
				return
			}
			typeDB, err := strconv.Atoi(os.Getenv("REDIS_UsersCart"))
			if err != nil {
				pkg.WriteApiResponse(w, nil, "Проблемы с типом БД", http.StatusBadRequest)
			}
			err = redishelper.Delete(typeDB, id.Hex())

		}
		pkg.WriteApiResponse(w, nil, "Прочитано: "+entity, http.StatusOK)
		return
	})
}

func StartServer() {

	s := &Serve{mux: http.NewServeMux()}

	createHandler(s)
	readHandler(s)
	updateHandler(s)
	deleteHandler(s)

	fmt.Printf("сервер запущен! http://localhost:8080/\n")
	err := http.ListenAndServe(":8080", s.mux)
	if err != nil {
		log.Fatalf("Сервер упал: %v", err)
	}
}

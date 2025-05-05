package http

import (
	"encoding/json"
	"io"
	"net/http"
)

type Serve struct {
	mux *http.ServeMux
}

type tableBody struct {
	name string `json:"name"`
}

func createHandler(s *Serve) {
	s.mux.HandleFunc("POST /create/{entity}", func(w http.ResponseWriter, r *http.Request) {
		entity := r.PathValue("entity")

		var tB tableBody

		if err := json.NewDecoder(r.Body).Decode(&tB); err != nil {
			if err == io.EOF {
				WriteApiResponse(w, nil, "Пустое тело запроса", http.StatusBadRequest)
				return
			}
			WriteApiResponse(w, nil, "Ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch entity {
		case "country":
			if tB.name == "" {
				WriteApiResponse(w, nil, "Поле name не может быть пустым", http.StatusBadRequest)
			}
			Create(entity, "hello")
		}

		WriteApiResponse(w, nil, "Создано: "+entity, http.StatusOK)
	})

	s.mux.HandleFunc("GET /read/{entity}", func(w http.ResponseWriter, r *http.Request) {
		entity := r.PathValue("entity")
		WriteApiResponse(w, nil, "Прочитано: "+entity, http.StatusOK)
	})
}

func readHandler(s *Serve)   {}
func updateHandler(s *Serve) {}
func deleteHandler(s *Serve) {}

func StartServer() {
	s := &Serve{mux: http.NewServeMux()}

	createHandler(s)
	readHandler(s)
	updateHandler(s)
	deleteHandler(s)

	print("сервер запущен! http://localhost:8080/")
	err := http.ListenAndServe(":8080", s.mux)
	if err != nil {
		return
	}
}

func WriteApiResponse(w http.ResponseWriter, result any, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := ApiResponse{
		Message: message,
		Data:    result,
	}
	json.NewEncoder(w).Encode(response)
}

package http

import (
	"encoding/json"
	"io"
	"net/http"
)

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

		var answer string
		switch entity {
		case "country":
			if tB.Name == "" {
				WriteApiResponse(w, nil, "Поле name не может быть пустым", http.StatusBadRequest)
				return
			}
			answer = CreateCountry(tB.Name)
		}

		WriteApiResponse(w, nil, answer, http.StatusOK)

	})
}

func readHandler(s *Serve) {
	s.mux.HandleFunc("GET /read/{entity}", func(w http.ResponseWriter, r *http.Request) {
		entity := r.PathValue("entity")
		switch entity {
		case "country":
			answer, err := ReadCountry()
			if err != nil {
				WriteApiResponse(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}

			WriteApiResponse(w, answer, "", http.StatusOK)
			return
		}
		WriteApiResponse(w, nil, "Прочитано: "+entity, http.StatusOK)
		return
	})
}
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

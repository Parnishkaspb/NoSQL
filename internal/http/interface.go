package http

import "net/http"

type APIArm interface {
	Create()
	Read()
	Update()
	Delete()
}

type Server interface {
	createHandler()
	readHandler()
	updateHandler()
	deleteHandler()

	WriteApiResponse(w http.ResponseWriter, result any, message string, code int)
	StartServer()
}

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

package http

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type APIArm interface {
	CreateCountry()
	ReadCountry()
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

type Serve struct {
	mux *http.ServeMux
}

type tableBody struct {
	Name string `json:"name"`
}

type Country struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

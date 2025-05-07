package http

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type APIArm interface {
	Create(string, any) (string, error)
	ReadAll(string, any) ([]any, error)
	Update(string, any) (string, error)
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

type Serve struct {
	mux *http.ServeMux
}

type tableBody struct {
	Name        string             `json:"name"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Price       float64            `json:"price"`
	Age         int                `json:"age"`
	ApprovedAge *bool              `json:"approved_age"`
	Surname     string             `json:"surname"`
	Login       string             `json:"login"`
	Password    string             `json:"password"`
	GameID      primitive.ObjectID `bson:"game_id"`
	UserID      primitive.ObjectID `bson:"user_id"`
	Rating      int                `bson:"rating"`
	Review      string             `bson:"text"`
}

type Country struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

type Platform struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

type Game struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	Age         int                `bson:"age"`
	ApprovedAge bool               `bson:"approved_age"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Surname  string             `bson:"surname"`
	Login    string             `bson:"login"`
	Password string             `bson:"password"`
	Age      int                `bson:"age"`
}

type GameReview struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	GameID    primitive.ObjectID `bson:"game_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Rating    int                `bson:"rating"`
	Review    string             `bson:"text"`
	CreatedAt time.Time          `bson:"created_at"`
}

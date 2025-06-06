package http

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	Name        string  `json:"name"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Age         int     `json:"age"`
	ApprovedAge *bool   `json:"approved_age"`
	Surname     string  `json:"surname"`
	Login       string  `json:"login"`
	Password    string  `json:"password"`
	GameID      string  `json:"game_id"`
	UserID      string  `json:"user_id"`
	Rating      int     `json:"rating"`
	Review      string  `json:"text"`
	Discount    int     `json:"discount"`
	Quantity    int     `json:"quantity"`
	Code        string  `json:"code"`
	FriendID    string  `json:"friend_id"`
	PlatformID  string  `json:"platform_id"`
	CountryID   string  `json:"country_id"`
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

type GameActualPrice struct {
	Price    float64 `bson:"price"`
	Discount int     `bson:"discount"`
}

type UsersCart struct {
	GameID   primitive.ObjectID `bson:"game_id"`
	Quantity int                `bson:"quantity"`
	Price    float64            `bson:"price"`
}

type UsersResetPassword struct {
	Code string `json:"code"`
}

type Router struct {
	mux *http.ServeMux
}

type httpServerStruct struct {
	httpServer *http.Server
}

type UserFriend struct {
	UserID   primitive.ObjectID `bson:"user_id"`
	FriendID primitive.ObjectID `bson:"friend_id"`
}

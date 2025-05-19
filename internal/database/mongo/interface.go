package mongohelper

import "go.mongodb.org/mongo-driver/bson"

type Mongo interface {
	ConnectMongo()
	ReadAll(collectionName string) ([]any, error)
	Create(collectionName string, value any) (string, error)
	Read(collectionName string, filter any) ([]any, error)
	DeleteT(collectionName string, filter bson.M) (int64, error)
}

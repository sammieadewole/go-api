package repo

import (
	"context"
	"go-api/db"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Mongo's Generic repo struct
type MongoRepo[T interface{}] struct {
	collectionName string
}

// Creates a new mongodb repo
func NewMongoRepo[T interface{}](collectionName string) *MongoRepo[T] {
	return &MongoRepo[T]{collectionName: collectionName}
}

// Creates a new entity
func (r *MongoRepo[T]) Create(entity *T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.GetMongoCollection(r.collectionName).InsertOne(ctx, entity)
	return err
}

// Gets all entities
func (r *MongoRepo[T]) Get() ([]*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var entities []*T
	filter := bson.M{"deleted_at": nil}
	cursor, err := db.GetMongoCollection(r.collectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}

// Gets one entity by ID
func (r *MongoRepo[T]) GetByID(id string) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var entity T
	filter := bson.M{"_id": id}
	if err := db.GetMongoCollection(r.collectionName).FindOne(ctx, filter).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

// Gets one entity by email
func (r *MongoRepo[T]) GetByEmail(email string) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var entity T
	filter := bson.M{"email": email, "deleted_at": nil}

	err := db.GetMongoCollection(r.collectionName).FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &entity, nil
}

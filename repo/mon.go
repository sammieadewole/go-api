package repo

import (
	"context"
	"go-api/db"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Mongo's Generic repo struct
type MongoRepo[T any] struct {
	collectionName string
}

// Creates a new mongodb repo
func NewMongoRepo[T any](collectionName string) *MongoRepo[T] {
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

// Gets one entity by id
func (r *MongoRepo[T]) GetOne(id string) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var entity T
	filter := bson.M{"id": id, "deleted_at": nil}

	err := db.GetMongoCollection(r.collectionName).FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
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

// Updates an entity
func (r *MongoRepo[T]) Update(id string, entity *T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	setEntityID(entity, id)
	filter := bson.M{"id": id, "deleted_at": nil}
	update := bson.M{"$set": entity}

	result, err := db.GetMongoCollection(r.collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Soft deletes an entity
func (r *MongoRepo[T]) SoftDelete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id, "deleted_at": nil}
	now := time.Now()
	update := bson.M{"$set": bson.M{"deleted_at": &now}}

	result, err := db.GetMongoCollection(r.collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Hard deletes an entity
func (r *MongoRepo[T]) HardDelete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	result, err := db.GetMongoCollection(r.collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Admin functions

// Gets all entities whether soft deleted or not
func (r *MongoRepo[T]) AdminGet() ([]*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var entities []*T
	cursor, err := db.GetMongoCollection(r.collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}

// Gets one entity by id whether soft deleted or not
func (r *MongoRepo[T]) AdminGetOne(id string) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var entity T
	filter := bson.M{"id": id}

	err := db.GetMongoCollection(r.collectionName).FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Gets one entity by email whether soft deleted or not
func (r *MongoRepo[T]) AdminGetByEmail(email string) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var entity T
	filter := bson.M{"email": email}

	err := db.GetMongoCollection(r.collectionName).FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Updates an entity whether soft deleted or not
func (r *MongoRepo[T]) AdminUpdate(id string, entity *T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	setEntityID(entity, id)
	filter := bson.M{"id": id}
	update := bson.M{"$set": entity}

	result, err := db.GetMongoCollection(r.collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

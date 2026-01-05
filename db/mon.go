package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client

// Connect to Mongo DB
func ConnectMongo() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Println("You must set your 'MONGODB_URI' environment variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}
	// Uses the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// Defines the options for the MongoDB client
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Creates a new client and connects to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Sends a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	log.Println("Successfully connected to MongoDB!")

	MongoClient = client
}

// Get a collection
//
// - name: Colection name
func GetMongoCollection(name string) *mongo.Collection {
	dbName := os.Getenv("MONGO_DB_NAME")
	return MongoClient.Database(dbName).Collection(name)
}

// Migrates all collections models to database
//
// - collectionNames: Collection names
func MigrateMongo(collectionNames ...string) error {
	if MongoClient == nil {
		return fmt.Errorf("MONGODB is not connected")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		return fmt.Errorf("MONGO_DB_NAME is not set")
	}

	db := MongoClient.Database(dbName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, name := range collectionNames {
		collections, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: name}})
		if err != nil {
			return fmt.Errorf("Failed to find collections: %w", err)
		}

		if len(collections) == 0 {
			if err := db.CreateCollection(ctx, name); err != nil {
				fmt.Printf("Collection %s created in MONGODB\n", name)
			}
		}
	}
	return nil
}

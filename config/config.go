package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = CreateMongoClient()

func CreateMongoClient() *mongo.Client {
	mongoURI := os.Getenv("MONGO_SERVER")

	if mongoURI == "" {
		log.Fatal("[ConnectDB] - Error: MONGO_SERVER is not set in .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("[ConnectDB] - Error: Connection error!")
	}

	// ตรวจสอบการเชื่อมต่อด้วย Ping
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("[ConnectDB] - Error: Ping error!")
	}

	log.Println("[ConnectDB] - DEBUG: Connected to MongoDB!")

	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var database = os.Getenv("MONGO_INITDB_DATABASE")

	return client.Database(database).Collection(collectionName)
}

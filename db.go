package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func db() *mongo.Client {
	godotenv.Load()
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGOLINK"))
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}
